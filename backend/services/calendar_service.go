package services

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// MonthlyCalendar holds the full monthly calendar view for a user.
type MonthlyCalendar struct {
	Year            int           `json:"year"`
	Month           int           `json:"month"`
	TotalAmount     int           `json:"totalAmount"`
	TotalCount      int           `json:"totalCount"`
	RemainingAmount int           `json:"remainingAmount"`
	RemainingCount  int           `json:"remainingCount"`
	Days            []CalendarDay `json:"days"`
}

// CalendarDay represents a single day with scheduled billing subscriptions.
type CalendarDay struct {
	Date          string                 `json:"date"`
	TotalAmount   int                    `json:"totalAmount"`
	Subscriptions []CalendarSubscription `json:"subscriptions"`
}

// CalendarSubscription represents a subscription entry within a calendar day.
type CalendarSubscription struct {
	SubscriptionID string `json:"subscriptionId"`
	ServiceName    string `json:"serviceName"`
	Amount         int    `json:"amount"`
	MonthlyAmount  int    `json:"monthlyAmount"`
	PersonalAmount int    `json:"personalAmount"`
	BillingCycle   string `json:"billingCycle"`
	CategoryName   string `json:"categoryName"`
	CategoryColor  string `json:"categoryColor"`
	AutoRenew      bool   `json:"autoRenew"`
}

// CalendarService handles calendar-related business logic.
type CalendarService struct {
	subRepo   repositories.SubscriptionRepository
	shareRepo repositories.SubscriptionShareRepository
}

// NewCalendarService creates a new CalendarService.
func NewCalendarService(subRepo repositories.SubscriptionRepository, shareRepo repositories.SubscriptionShareRepository) *CalendarService {
	return &CalendarService{subRepo: subRepo, shareRepo: shareRepo}
}

// GetMonthlyCalendar returns the monthly calendar with billing schedule for a user.
func (s *CalendarService) GetMonthlyCalendar(userID string, year, month int) (*MonthlyCalendar, error) {
	// Fetch active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("캘린더 활성 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("캘린더 데이터를 조회할 수 없습니다")
	}

	// Build share map for personal amount calculation.
	shareMap := buildShareMap(s.shareRepo, userID)

	// Target month boundaries.
	targetStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	targetEnd := targetStart.AddDate(0, 1, -1) // last day of month
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Group subscriptions by billing day in the target month.
	dayMap := make(map[int][]CalendarSubscription) // day-of-month -> subscriptions

	totalAmount := 0
	totalCount := 0
	remainingAmount := 0
	remainingCount := 0

	for _, sub := range activeSubs {
		billingDay, inMonth := s.billingDayInMonth(sub, year, month)
		if !inMonth {
			continue
		}

		// Clamp billing day to last day of month.
		lastDay := targetEnd.Day()
		if billingDay > lastDay {
			billingDay = lastDay
		}

		monthlyAmt := sub.MonthlyAmount()
		personalAmt := monthlyAmt
		if share, ok := shareMap[sub.ID.String()]; ok {
			personalAmt = share.PersonalAmount(monthlyAmt)
		}

		catName := "미분류"
		catColor := "#9E9E9E"
		if sub.CategoryID != nil && sub.Category != nil {
			catName = sub.Category.Name
			if sub.Category.Color != nil {
				catColor = *sub.Category.Color
			}
		}

		entry := CalendarSubscription{
			SubscriptionID: sub.ID.String(),
			ServiceName:    sub.ServiceName,
			Amount:         sub.Amount,
			MonthlyAmount:  monthlyAmt,
			PersonalAmount: personalAmt,
			BillingCycle:   string(sub.BillingCycle),
			CategoryName:   catName,
			CategoryColor:  catColor,
			AutoRenew:      sub.AutoRenew,
		}

		dayMap[billingDay] = append(dayMap[billingDay], entry)

		totalAmount += personalAmt
		totalCount++

		billingDate := time.Date(year, time.Month(month), billingDay, 0, 0, 0, 0, time.UTC)
		if !billingDate.Before(today) {
			remainingAmount += personalAmt
			remainingCount++
		}
	}

	// Build sorted CalendarDay slice.
	days := make([]CalendarDay, 0, len(dayMap))
	for day, subs := range dayMap {
		dayTotal := 0
		for _, s := range subs {
			dayTotal += s.PersonalAmount
		}
		days = append(days, CalendarDay{
			Date:          fmt.Sprintf("%04d-%02d-%02d", year, month, day),
			TotalAmount:   dayTotal,
			Subscriptions: subs,
		})
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date < days[j].Date
	})

	return &MonthlyCalendar{
		Year:            year,
		Month:           month,
		TotalAmount:     totalAmount,
		TotalCount:      totalCount,
		RemainingAmount: remainingAmount,
		RemainingCount:  remainingCount,
		Days:            days,
	}, nil
}

// DayDetail holds the billing details for a specific day.
type DayDetail struct {
	Date          string                 `json:"date"`
	TotalAmount   int                    `json:"totalAmount"`
	Subscriptions []CalendarSubscription `json:"subscriptions"`
}

// GetDayDetail returns the billing subscriptions for a specific date.
func (s *CalendarService) GetDayDetail(userID string, year, month, day int) (*DayDetail, error) {
	// Fetch active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("일별 결제 상세 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("일별 결제 데이터를 조회할 수 없습니다")
	}

	shareMap := buildShareMap(s.shareRepo, userID)

	// Last day of the target month.
	targetEnd := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	lastDay := targetEnd.Day()

	result := &DayDetail{
		Date:          fmt.Sprintf("%04d-%02d-%02d", year, month, day),
		TotalAmount:   0,
		Subscriptions: make([]CalendarSubscription, 0),
	}

	for _, sub := range activeSubs {
		billingDay, inMonth := s.billingDayInMonth(sub, year, month)
		if !inMonth {
			continue
		}

		// Clamp to last day of month.
		if billingDay > lastDay {
			billingDay = lastDay
		}

		if billingDay != day {
			continue
		}

		monthlyAmt := sub.MonthlyAmount()
		personalAmt := monthlyAmt
		if share, ok := shareMap[sub.ID.String()]; ok {
			personalAmt = share.PersonalAmount(monthlyAmt)
		}

		catName := "미분류"
		catColor := "#9E9E9E"
		if sub.CategoryID != nil && sub.Category != nil {
			catName = sub.Category.Name
			if sub.Category.Color != nil {
				catColor = *sub.Category.Color
			}
		}

		entry := CalendarSubscription{
			SubscriptionID: sub.ID.String(),
			ServiceName:    sub.ServiceName,
			Amount:         sub.Amount,
			MonthlyAmount:  monthlyAmt,
			PersonalAmount: personalAmt,
			BillingCycle:   string(sub.BillingCycle),
			CategoryName:   catName,
			CategoryColor:  catColor,
			AutoRenew:      sub.AutoRenew,
		}

		result.Subscriptions = append(result.Subscriptions, entry)
		result.TotalAmount += personalAmt
	}

	return result, nil
}

// UpcomingPayment represents a single upcoming payment entry.
type UpcomingPayment struct {
	Date           string `json:"date"`
	DaysUntil      int    `json:"daysUntil"`
	SubscriptionID string `json:"subscriptionId"`
	ServiceName    string `json:"serviceName"`
	Amount         int    `json:"amount"`
	PersonalAmount int    `json:"personalAmount"`
	CategoryName   string `json:"categoryName"`
	CategoryColor  string `json:"categoryColor"`
}

// GetUpcomingPayments returns payments due within the next N days.
func (s *CalendarService) GetUpcomingPayments(userID string, days int) ([]UpcomingPayment, error) {
	if days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("예정 결제 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("예정 결제 데이터를 조회할 수 없습니다")
	}

	shareMap := buildShareMap(s.shareRepo, userID)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	deadline := today.AddDate(0, 0, days)

	payments := make([]UpcomingPayment, 0)

	for _, sub := range activeSubs {
		nbd := sub.NextBillingDate
		// Only include if NextBillingDate is within range [today, deadline].
		if nbd.Before(today) || nbd.After(deadline) {
			continue
		}

		daysUntil := int(nbd.Sub(today).Hours() / 24)

		monthlyAmt := sub.MonthlyAmount()
		personalAmt := monthlyAmt
		if share, ok := shareMap[sub.ID.String()]; ok {
			personalAmt = share.PersonalAmount(monthlyAmt)
		}

		catName := "미분류"
		catColor := "#9E9E9E"
		if sub.CategoryID != nil && sub.Category != nil {
			catName = sub.Category.Name
			if sub.Category.Color != nil {
				catColor = *sub.Category.Color
			}
		}

		payments = append(payments, UpcomingPayment{
			Date:           nbd.Format("2006-01-02"),
			DaysUntil:      daysUntil,
			SubscriptionID: sub.ID.String(),
			ServiceName:    sub.ServiceName,
			Amount:         sub.Amount,
			PersonalAmount: personalAmt,
			CategoryName:   catName,
			CategoryColor:  catColor,
		})
	}

	// Sort by date ascending.
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].Date < payments[j].Date
	})

	return payments, nil
}

// billingDayInMonth determines whether a subscription has a billing event in
// the given year/month and returns the day-of-month for that event.
func (s *CalendarService) billingDayInMonth(sub *models.Subscription, year, month int) (int, bool) {
	nbd := sub.NextBillingDate

	switch sub.BillingCycle {
	case models.BillingCycleMonthly:
		// Monthly subscriptions bill on the same day every month.
		return nbd.Day(), true

	case models.BillingCycleYearly:
		// Yearly subscriptions bill only if the month matches.
		if int(nbd.Month()) == month {
			return nbd.Day(), true
		}
		return 0, false

	case models.BillingCycleWeekly:
		// Weekly subscriptions: check if NextBillingDate itself falls in this month.
		if nbd.Year() == year && int(nbd.Month()) == month {
			return nbd.Day(), true
		}
		return 0, false

	default:
		// Fallback: check if NextBillingDate is in the target month.
		if nbd.Year() == year && int(nbd.Month()) == month {
			return nbd.Day(), true
		}
		return 0, false
	}
}
