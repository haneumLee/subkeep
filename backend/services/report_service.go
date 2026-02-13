package services

import (
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// ReportOverview holds all report data for the overview endpoint.
type ReportOverview struct {
	CategoryBreakdown []CategoryBreakdown `json:"categoryBreakdown"`
	MonthlyTrend      []MonthlyTrend      `json:"monthlyTrend"`
	AverageCost       AverageCost         `json:"averageCost"`
	Summary           ReportSummary       `json:"summary"`
}

// MonthlyTrend represents the cost trend for a specific month.
type MonthlyTrend struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Amount int `json:"amount"`
	Count  int `json:"count"`
}

// AverageCost holds weekly, monthly, and annual average costs.
type AverageCost struct {
	Monthly int `json:"monthly"`
	Annual  int `json:"annual"`
	Weekly  int `json:"weekly"`
}

// ReportSummary holds high-level subscription statistics.
type ReportSummary struct {
	TotalSubscriptions  int     `json:"totalSubscriptions"`
	ActiveCount         int     `json:"activeCount"`
	PausedCount         int     `json:"pausedCount"`
	MostExpensive       *string `json:"mostExpensive,omitempty"`
	MostExpensiveAmount int     `json:"mostExpensiveAmount"`
	AverageSatisfaction float64 `json:"averageSatisfaction"`
}

// ReportService handles report-related business logic.
type ReportService struct {
	subRepo   repositories.SubscriptionRepository
	shareRepo repositories.SubscriptionShareRepository
}

// NewReportService creates a new ReportService.
func NewReportService(subRepo repositories.SubscriptionRepository, shareRepo repositories.SubscriptionShareRepository) *ReportService {
	return &ReportService{subRepo: subRepo, shareRepo: shareRepo}
}

// GetOverview returns the full report overview for a user.
func (s *ReportService) GetOverview(userID string) (*ReportOverview, error) {
	// Fetch active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("리포트 활성 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("리포트 데이터를 조회할 수 없습니다")
	}

	// Fetch paused subscriptions.
	pausedSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "paused",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("리포트 일시정지 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("리포트 데이터를 조회할 수 없습니다")
	}

	// Build share map for personal amount calculation.
	shareMap := buildShareMap(s.shareRepo, userID)

	// Combine active + paused for trend calculation.
	allSubs := append(activeSubs, pausedSubs...)

	// Calculate personal monthly amounts for each subscription.
	subCosts := make([]subWithCostEntry, 0, len(allSubs))
	for _, sub := range allSubs {
		monthly := sub.MonthlyAmount()
		personal := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personal = share.PersonalAmount(monthly)
		}
		subCosts = append(subCosts, subWithCostEntry{sub: sub, personalAmount: personal})
	}

	// --- Category Breakdown (active only) ---
	categoryBreakdown := s.buildCategoryBreakdown(activeSubs, shareMap)

	// --- Monthly Trend (last 12 months) ---
	monthlyTrend := s.buildMonthlyTrend(subCosts)

	// --- Average Cost (active subscriptions only) ---
	averageCost := s.buildAverageCost(activeSubs, shareMap)

	// --- Report Summary ---
	summary := s.buildSummary(activeSubs, pausedSubs, shareMap)

	return &ReportOverview{
		CategoryBreakdown: categoryBreakdown,
		MonthlyTrend:      monthlyTrend,
		AverageCost:       averageCost,
		Summary:           summary,
	}, nil
}

// buildCategoryBreakdown calculates per-category spending breakdown for active subscriptions.
func (s *ReportService) buildCategoryBreakdown(activeSubs []*models.Subscription, shareMap map[string]*models.SubscriptionShare) []CategoryBreakdown {
	type catGroup struct {
		categoryID   string
		categoryName string
		color        string
		amount       int
		count        int
	}
	categoryMap := make(map[string]*catGroup)
	monthlyTotal := 0

	for _, sub := range activeSubs {
		monthly := sub.MonthlyAmount()
		personal := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personal = share.PersonalAmount(monthly)
		}
		monthlyTotal += personal

		catID := "uncategorized"
		catName := "미분류"
		catColor := "#9E9E9E"

		if sub.CategoryID != nil && sub.Category != nil {
			catID = sub.CategoryID.String()
			catName = sub.Category.Name
			if sub.Category.Color != nil {
				catColor = *sub.Category.Color
			}
		}

		if g, ok := categoryMap[catID]; ok {
			g.amount += personal
			g.count++
		} else {
			categoryMap[catID] = &catGroup{
				categoryID:   catID,
				categoryName: catName,
				color:        catColor,
				amount:       personal,
				count:        1,
			}
		}
	}

	breakdown := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, g := range categoryMap {
		pct := 0.0
		if monthlyTotal > 0 {
			pct = math.Round(float64(g.amount)/float64(monthlyTotal)*1000) / 10
		}
		breakdown = append(breakdown, CategoryBreakdown{
			CategoryID:    g.categoryID,
			CategoryName:  g.categoryName,
			Color:         g.color,
			MonthlyAmount: g.amount,
			Percentage:    pct,
			Count:         g.count,
		})
	}

	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].MonthlyAmount > breakdown[j].MonthlyAmount
	})

	return breakdown
}

// buildMonthlyTrend calculates the cost trend for the last 12 months.
// A subscription is considered active in a month if:
//   - StartDate <= last day of that month
//   - Status is active or paused
type subWithCostEntry struct {
	sub            *models.Subscription
	personalAmount int
}

func (s *ReportService) buildMonthlyTrend(subCosts []subWithCostEntry) []MonthlyTrend {
	now := time.Now()
	trends := make([]MonthlyTrend, 12)

	for i := 0; i < 12; i++ {
		// Calculate target month: from 11 months ago to current month.
		offset := 11 - i
		targetDate := now.AddDate(0, -offset, 0)
		year := targetDate.Year()
		month := int(targetDate.Month())

		// Last day of the target month.
		lastDay := time.Date(year, targetDate.Month()+1, 0, 23, 59, 59, 0, time.UTC)

		totalAmount := 0
		count := 0

		for _, sc := range subCosts {
			// Subscription was active in this month if it started on or before the last day.
			if sc.sub.StartDate.Before(lastDay) || sc.sub.StartDate.Equal(lastDay) {
				totalAmount += sc.personalAmount
				count++
			}
		}

		trends[i] = MonthlyTrend{
			Year:   year,
			Month:  month,
			Amount: totalAmount,
			Count:  count,
		}
	}

	return trends
}

// buildAverageCost calculates current average costs based on active subscriptions.
func (s *ReportService) buildAverageCost(activeSubs []*models.Subscription, shareMap map[string]*models.SubscriptionShare) AverageCost {
	monthlyTotal := 0
	for _, sub := range activeSubs {
		monthly := sub.MonthlyAmount()
		personal := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personal = share.PersonalAmount(monthly)
		}
		monthlyTotal += personal
	}

	annual := monthlyTotal * 12
	weekly := 0
	if monthlyTotal > 0 {
		weekly = int(math.Round(float64(annual) / 52.0))
	}

	return AverageCost{
		Monthly: monthlyTotal,
		Annual:  annual,
		Weekly:  weekly,
	}
}

// buildSummary builds the report summary statistics.
func (s *ReportService) buildSummary(activeSubs, pausedSubs []*models.Subscription, shareMap map[string]*models.SubscriptionShare) ReportSummary {
	totalCount := len(activeSubs) + len(pausedSubs)

	var mostExpensiveName *string
	mostExpensiveAmount := 0
	satisfactionSum := 0.0
	satisfactionCount := 0

	allSubs := append(activeSubs, pausedSubs...)
	for _, sub := range allSubs {
		// Track satisfaction scores.
		if sub.SatisfactionScore != nil {
			satisfactionSum += float64(*sub.SatisfactionScore)
			satisfactionCount++
		}

		// Find most expensive by personal amount.
		monthly := sub.MonthlyAmount()
		personal := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personal = share.PersonalAmount(monthly)
		}

		if personal > mostExpensiveAmount {
			mostExpensiveAmount = personal
			name := sub.ServiceName
			mostExpensiveName = &name
		}
	}

	avgSatisfaction := 0.0
	if satisfactionCount > 0 {
		avgSatisfaction = math.Round(satisfactionSum/float64(satisfactionCount)*10) / 10
	}

	return ReportSummary{
		TotalSubscriptions:  totalCount,
		ActiveCount:         len(activeSubs),
		PausedCount:         len(pausedSubs),
		MostExpensive:       mostExpensiveName,
		MostExpensiveAmount: mostExpensiveAmount,
		AverageSatisfaction: avgSatisfaction,
	}
}
