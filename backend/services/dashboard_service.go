package services

import (
	"log/slog"
	"math"
	"sort"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// DashboardSummary holds the overall spending summary for a user.
type DashboardSummary struct {
	MonthlyTotal      int                 `json:"monthlyTotal"`
	AnnualTotal       int                 `json:"annualTotal"`
	ActiveCount       int                 `json:"activeCount"`
	PausedCount       int                 `json:"pausedCount"`
	CategoryBreakdown []CategoryBreakdown `json:"categoryBreakdown"`
}

// CategoryBreakdown represents spending breakdown per category.
type CategoryBreakdown struct {
	CategoryID    string  `json:"categoryId"`
	CategoryName  string  `json:"categoryName"`
	Color         string  `json:"color"`
	MonthlyAmount int     `json:"monthlyAmount"`
	Percentage    float64 `json:"percentage"`
	Count         int     `json:"count"`
}

// CancelRecommendation represents a subscription recommended for cancellation.
type CancelRecommendation struct {
	SubscriptionID    string `json:"subscriptionId"`
	ServiceName       string `json:"serviceName"`
	MonthlyAmount     int    `json:"monthlyAmount"`
	AnnualSaving      int    `json:"annualSaving"`
	SatisfactionScore *int   `json:"satisfactionScore"`
	Reason            string `json:"reason"`
}

// DashboardService handles dashboard-related business logic.
type DashboardService struct {
	subRepo repositories.SubscriptionRepository
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(subRepo repositories.SubscriptionRepository) *DashboardService {
	return &DashboardService{subRepo: subRepo}
}

// GetSummary returns the overall spending summary for a user.
func (s *DashboardService) GetSummary(userID string) (*DashboardSummary, error) {
	// Fetch active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("대시보드 활성 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("대시보드 데이터를 조회할 수 없습니다")
	}

	// Fetch paused subscriptions count.
	pausedSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "paused",
		Page:    1,
		PerPage: 1,
	})
	if err != nil {
		slog.Error("대시보드 일시정지 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("대시보드 데이터를 조회할 수 없습니다")
	}

	// Calculate monthly total and category grouping.
	monthlyTotal := 0
	type catGroup struct {
		categoryID   string
		categoryName string
		color        string
		amount       int
		count        int
	}
	categoryMap := make(map[string]*catGroup)

	for _, sub := range activeSubs {
		monthly := sub.MonthlyAmount()
		monthlyTotal += monthly

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
			g.amount += monthly
			g.count++
		} else {
			categoryMap[catID] = &catGroup{
				categoryID:   catID,
				categoryName: catName,
				color:        catColor,
				amount:       monthly,
				count:        1,
			}
		}
	}

	// Build category breakdown.
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

	// Sort by amount descending.
	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].MonthlyAmount > breakdown[j].MonthlyAmount
	})

	// Count paused.
	pausedCount := len(pausedSubs)

	return &DashboardSummary{
		MonthlyTotal:      monthlyTotal,
		AnnualTotal:       monthlyTotal * 12,
		ActiveCount:       len(activeSubs),
		PausedCount:       pausedCount,
		CategoryBreakdown: breakdown,
	}, nil
}

// GetRecommendations returns cancel recommendations based on satisfaction and cost.
func (s *DashboardService) GetRecommendations(userID string) ([]*CancelRecommendation, error) {
	// Fetch all active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("해지 추천 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("해지 추천 데이터를 조회할 수 없습니다")
	}

	if len(activeSubs) == 0 {
		return []*CancelRecommendation{}, nil
	}

	// Calculate monthly amounts and find top 20% cost threshold.
	type subWithCost struct {
		sub     *models.Subscription
		monthly int
	}
	items := make([]subWithCost, len(activeSubs))
	for i, sub := range activeSubs {
		items[i] = subWithCost{sub: sub, monthly: sub.MonthlyAmount()}
	}

	// Sort by cost descending to find top 20% threshold.
	sort.Slice(items, func(i, j int) bool {
		return items[i].monthly > items[j].monthly
	})

	top20Index := int(math.Ceil(float64(len(items)) * 0.2))
	costThreshold := 0
	if top20Index > 0 && top20Index <= len(items) {
		costThreshold = items[top20Index-1].monthly
	}

	// Build recommendations.
	recommendations := make([]*CancelRecommendation, 0)
	for _, item := range items {
		reason := ""

		// Criteria 1: satisfaction_score 1-2 (any amount).
		if item.sub.SatisfactionScore != nil && *item.sub.SatisfactionScore <= 2 {
			reason = "만족도 낮음"
		}

		// Criteria 2: top 20% cost AND satisfaction_score <= 3.
		if reason == "" && item.monthly >= costThreshold && item.sub.SatisfactionScore != nil && *item.sub.SatisfactionScore <= 3 {
			reason = "높은 비용 대비 낮은 만족도"
		}

		if reason == "" {
			continue
		}

		recommendations = append(recommendations, &CancelRecommendation{
			SubscriptionID:    item.sub.ID.String(),
			ServiceName:       item.sub.ServiceName,
			MonthlyAmount:     item.monthly,
			AnnualSaving:      item.monthly * 12,
			SatisfactionScore: item.sub.SatisfactionScore,
			Reason:            reason,
		})
	}

	// Sort by satisfaction ASC then monthlyAmount DESC.
	sort.Slice(recommendations, func(i, j int) bool {
		si := 0
		sj := 0
		if recommendations[i].SatisfactionScore != nil {
			si = *recommendations[i].SatisfactionScore
		}
		if recommendations[j].SatisfactionScore != nil {
			sj = *recommendations[j].SatisfactionScore
		}
		if si != sj {
			return si < sj
		}
		return recommendations[i].MonthlyAmount > recommendations[j].MonthlyAmount
	})

	return recommendations, nil
}
