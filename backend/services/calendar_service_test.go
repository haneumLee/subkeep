package services

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
)

// ---------------------------------------------------------------------------
// Mock repositories for calendar tests
// ---------------------------------------------------------------------------

type mockSubRepoForCalendar struct {
	subs map[string]*models.Subscription
}

func newMockSubRepoForCalendar() *mockSubRepoForCalendar {
	return &mockSubRepoForCalendar{subs: make(map[string]*models.Subscription)}
}

func (m *mockSubRepoForCalendar) FindByID(id string) (*models.Subscription, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, nil
	}
	return sub, nil
}
func (m *mockSubRepoForCalendar) FindByUserID(userID string, filter repositories.SubscriptionFilter) ([]*models.Subscription, int64, error) {
	var result []*models.Subscription
	for _, sub := range m.subs {
		if sub.UserID.String() != userID {
			continue
		}
		if filter.Status != "" && string(sub.Status) != filter.Status {
			continue
		}
		result = append(result, sub)
	}
	return result, int64(len(result)), nil
}
func (m *mockSubRepoForCalendar) Create(sub *models.Subscription) error              { return nil }
func (m *mockSubRepoForCalendar) Update(sub *models.Subscription) error              { return nil }
func (m *mockSubRepoForCalendar) Delete(id string) error                             { return nil }
func (m *mockSubRepoForCalendar) Restore(id string) error                            { return nil }
func (m *mockSubRepoForCalendar) CountByUserID(userID string) (int64, error)         { return 0, nil }
func (m *mockSubRepoForCalendar) FindDuplicateName(userID, name string) (bool, error) { return false, nil }

type mockShareRepoForCalendar struct {
	shares map[string]*models.SubscriptionShare
}

func newMockShareRepoForCalendar() *mockShareRepoForCalendar {
	return &mockShareRepoForCalendar{shares: make(map[string]*models.SubscriptionShare)}
}

func (m *mockShareRepoForCalendar) FindByID(id string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForCalendar) FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForCalendar) FindByUserID(userID string) ([]*models.SubscriptionShare, error) {
	var result []*models.SubscriptionShare
	for _, s := range m.shares {
		result = append(result, s)
	}
	return result, nil
}
func (m *mockShareRepoForCalendar) Create(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForCalendar) Update(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForCalendar) Delete(id string) error                        { return nil }
func (m *mockShareRepoForCalendar) DeleteBySubscriptionID(subscriptionID string) error { return nil }

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func calMakeCategory(name, color string) *models.Category {
	return &models.Category{
		ID:    uuid.New(),
		Name:  name,
		Color: &color,
	}
}

func seedCalendarSub(
	repo *mockSubRepoForCalendar,
	userID uuid.UUID,
	name string,
	amount int,
	cycle models.BillingCycle,
	nextBilling time.Time,
	category *models.Category,
) *models.Subscription {
	sub := &models.Subscription{
		ID:              uuid.New(),
		UserID:          userID,
		ServiceName:     name,
		Amount:          amount,
		BillingCycle:    cycle,
		Currency:        "KRW",
		NextBillingDate: nextBilling,
		AutoRenew:       true,
		Status:          models.SubscriptionStatusActive,
		StartDate:       time.Now().AddDate(-1, 0, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if category != nil {
		sub.CategoryID = &category.ID
		sub.Category = category
	}
	repo.subs[sub.ID.String()] = sub
	return sub
}

// ===========================================================================
// GetMonthlyCalendar
// ===========================================================================

func TestGetMonthlyCalendar_NoSubscriptions(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.Year != 2026 || cal.Month != 3 {
		t.Errorf("expected 2026/3, got %d/%d", cal.Year, cal.Month)
	}
	if cal.TotalAmount != 0 {
		t.Errorf("expected totalAmount 0, got %d", cal.TotalAmount)
	}
	if cal.TotalCount != 0 {
		t.Errorf("expected totalCount 0, got %d", cal.TotalCount)
	}
	if len(cal.Days) != 0 {
		t.Errorf("expected 0 days, got %d", len(cal.Days))
	}
}

func TestGetMonthlyCalendar_SingleMonthly(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Monthly sub billing on the 15th.
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 1 {
		t.Errorf("expected 1 subscription, got %d", cal.TotalCount)
	}
	if cal.TotalAmount != 17000 {
		t.Errorf("expected totalAmount 17000, got %d", cal.TotalAmount)
	}
	if len(cal.Days) != 1 {
		t.Fatalf("expected 1 day entry, got %d", len(cal.Days))
	}
	if cal.Days[0].Date != "2026-03-15" {
		t.Errorf("expected date 2026-03-15, got %s", cal.Days[0].Date)
	}
}

func TestGetMonthlyCalendar_MultipleSubscriptions(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), nil)
	seedCalendarSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly,
		time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC), nil)
	// Yearly sub in March.
	seedCalendarSub(repo, userID, "iCloud", 132000, models.BillingCycleYearly,
		time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 3 {
		t.Errorf("expected 3 subscriptions, got %d", cal.TotalCount)
	}
	if len(cal.Days) != 3 {
		t.Errorf("expected 3 day entries, got %d", len(cal.Days))
	}
}

func TestGetMonthlyCalendar_SameDayGrouping(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Two subscriptions billing on the same day.
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)
	seedCalendarSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be grouped into a single day entry.
	if len(cal.Days) != 1 {
		t.Fatalf("expected 1 day entry, got %d", len(cal.Days))
	}
	if len(cal.Days[0].Subscriptions) != 2 {
		t.Errorf("expected 2 subscriptions on same day, got %d", len(cal.Days[0].Subscriptions))
	}
	if cal.Days[0].TotalAmount != 27900 {
		t.Errorf("expected day total 27900, got %d", cal.Days[0].TotalAmount)
	}
}

func TestGetMonthlyCalendar_DaySorting(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	seedCalendarSub(repo, userID, "Late", 5000, models.BillingCycleMonthly,
		time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), nil)
	seedCalendarSub(repo, userID, "Early", 3000, models.BillingCycleMonthly,
		time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cal.Days) != 2 {
		t.Fatalf("expected 2 days, got %d", len(cal.Days))
	}
	if cal.Days[0].Date != "2026-03-05" {
		t.Errorf("expected first day 2026-03-05, got %s", cal.Days[0].Date)
	}
	if cal.Days[1].Date != "2026-03-25" {
		t.Errorf("expected second day 2026-03-25, got %s", cal.Days[1].Date)
	}
}

func TestGetMonthlyCalendar_MonthlyBillingCycle(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Monthly shows in every month.
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), nil)

	// Should appear in January.
	cal1, _ := svc.GetMonthlyCalendar(userID.String(), 2026, 1)
	if cal1.TotalCount != 1 {
		t.Errorf("January: expected 1, got %d", cal1.TotalCount)
	}

	// Should appear in June too.
	cal6, _ := svc.GetMonthlyCalendar(userID.String(), 2026, 6)
	if cal6.TotalCount != 1 {
		t.Errorf("June: expected 1, got %d", cal6.TotalCount)
	}
}

func TestGetMonthlyCalendar_YearlyBillingCycle_MatchingMonth(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	seedCalendarSub(repo, userID, "iCloud", 132000, models.BillingCycleYearly,
		time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 1 {
		t.Errorf("expected 1 subscription in matching month, got %d", cal.TotalCount)
	}
	// Yearly 132000 / 12 = 11000.
	if cal.TotalAmount != 11000 {
		t.Errorf("expected monthly amount 11000, got %d", cal.TotalAmount)
	}
}

func TestGetMonthlyCalendar_YearlyBillingCycle_NonMatchingMonth(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	seedCalendarSub(repo, userID, "iCloud", 132000, models.BillingCycleYearly,
		time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 0 {
		t.Errorf("expected 0 subscriptions in non-matching month, got %d", cal.TotalCount)
	}
}

func TestGetMonthlyCalendar_WeeklyBillingCycle(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Weekly sub with NextBillingDate in March 2026.
	seedCalendarSub(repo, userID, "Gym", 5000, models.BillingCycleWeekly,
		time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 1 {
		t.Errorf("expected 1 subscription, got %d", cal.TotalCount)
	}

	// Should NOT appear in April.
	cal4, _ := svc.GetMonthlyCalendar(userID.String(), 2026, 4)
	if cal4.TotalCount != 0 {
		t.Errorf("expected 0 subscriptions in April, got %d", cal4.TotalCount)
	}
}

func TestGetMonthlyCalendar_MonthlyAmount_YearlySub(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Yearly: 120000 / 12 = 10000
	seedCalendarSub(repo, userID, "Yearly", 120000, models.BillingCycleYearly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.Days[0].Subscriptions[0].MonthlyAmount != 10000 {
		t.Errorf("expected monthly amount 10000, got %d", cal.Days[0].Subscriptions[0].MonthlyAmount)
	}
}

func TestGetMonthlyCalendar_MonthlyAmount_WeeklySub(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Weekly: 2500 * 52 / 12 = 10833.33 → 10833
	seedCalendarSub(repo, userID, "Weekly", 2500, models.BillingCycleWeekly,
		time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.Days[0].Subscriptions[0].MonthlyAmount != 10833 {
		t.Errorf("expected monthly amount 10833, got %d", cal.Days[0].Subscriptions[0].MonthlyAmount)
	}
}

func TestGetMonthlyCalendar_WithShareEqual(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	sub := seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	// Equal split with 4 members → 17000 / 4 = 4250
	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeEqual,
		TotalMembersSnapshot: 4,
	}

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalAmount != 4250 {
		t.Errorf("expected personal amount 4250, got %d", cal.TotalAmount)
	}
	if cal.Days[0].Subscriptions[0].PersonalAmount != 4250 {
		t.Errorf("expected subscription personal amount 4250, got %d", cal.Days[0].Subscriptions[0].PersonalAmount)
	}
}

func TestGetMonthlyCalendar_WithShareCustomAmount(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	sub := seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	shareAmt := 5000
	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeCustomAmount,
		MyShareAmount:        &shareAmt,
		TotalMembersSnapshot: 3,
	}

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalAmount != 5000 {
		t.Errorf("expected personal amount 5000, got %d", cal.TotalAmount)
	}
}

func TestGetMonthlyCalendar_BillingDayOverMonthEnd(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Monthly sub billing on the 31st → clamped to 28 in February.
	seedCalendarSub(repo, userID, "Service", 10000, models.BillingCycleMonthly,
		time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC), nil)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cal.Days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(cal.Days))
	}
	if cal.Days[0].Date != "2026-02-28" {
		t.Errorf("expected date 2026-02-28, got %s", cal.Days[0].Date)
	}
}

func TestGetMonthlyCalendar_CategoryInfo(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Subscription without category → "미분류"
	seedCalendarSub(repo, userID, "NoCat", 5000, models.BillingCycleMonthly,
		time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), nil)

	// Subscription with category
	cat := calMakeCategory("Entertainment", "#FF5722")
	seedCalendarSub(repo, userID, "WithCat", 8000, models.BillingCycleMonthly,
		time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC), cat)

	cal, err := svc.GetMonthlyCalendar(userID.String(), 2026, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cal.TotalCount != 2 {
		t.Fatalf("expected 2, got %d", cal.TotalCount)
	}

	found := map[string]CalendarSubscription{}
	for _, day := range cal.Days {
		for _, s := range day.Subscriptions {
			found[s.ServiceName] = s
		}
	}

	if found["NoCat"].CategoryName != "미분류" {
		t.Errorf("expected '미분류', got '%s'", found["NoCat"].CategoryName)
	}
	if found["NoCat"].CategoryColor != "#9E9E9E" {
		t.Errorf("expected '#9E9E9E', got '%s'", found["NoCat"].CategoryColor)
	}
	if found["WithCat"].CategoryName != "Entertainment" {
		t.Errorf("expected 'Entertainment', got '%s'", found["WithCat"].CategoryName)
	}
	if found["WithCat"].CategoryColor != "#FF5722" {
		t.Errorf("expected '#FF5722', got '%s'", found["WithCat"].CategoryColor)
	}
}

// ===========================================================================
// GetDayDetail
// ===========================================================================

func TestGetDayDetail_NoPayments(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	detail, err := svc.GetDayDetail(userID.String(), 2026, 3, 15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Date != "2026-03-15" {
		t.Errorf("expected date 2026-03-15, got %s", detail.Date)
	}
	if detail.TotalAmount != 0 {
		t.Errorf("expected totalAmount 0, got %d", detail.TotalAmount)
	}
	if len(detail.Subscriptions) != 0 {
		t.Errorf("expected 0 subscriptions, got %d", len(detail.Subscriptions))
	}
}

func TestGetDayDetail_WithPayments(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)
	seedCalendarSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)
	// Different day — should NOT appear.
	seedCalendarSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly,
		time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC), nil)

	detail, err := svc.GetDayDetail(userID.String(), 2026, 3, 15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(detail.Subscriptions) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(detail.Subscriptions))
	}
	if detail.TotalAmount != 27900 {
		t.Errorf("expected totalAmount 27900, got %d", detail.TotalAmount)
	}
}

func TestGetDayDetail_ClampedBillingDay(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	// Billing on 31st → clamped to 28 in Feb.
	seedCalendarSub(repo, userID, "Service", 10000, models.BillingCycleMonthly,
		time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC), nil)

	detail, err := svc.GetDayDetail(userID.String(), 2026, 2, 28)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(detail.Subscriptions) != 1 {
		t.Errorf("expected 1 subscription clamped to 28th, got %d", len(detail.Subscriptions))
	}
}

func TestGetDayDetail_WithShare(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	sub := seedCalendarSub(repo, userID, "Netflix", 20000, models.BillingCycleMonthly,
		time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), nil)

	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeEqual,
		TotalMembersSnapshot: 2,
	}

	detail, err := svc.GetDayDetail(userID.String(), 2026, 3, 15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 20000 / 2 = 10000
	if detail.TotalAmount != 10000 {
		t.Errorf("expected personal amount 10000, got %d", detail.TotalAmount)
	}
}

// ===========================================================================
// GetUpcomingPayments
// ===========================================================================

func TestGetUpcomingPayments_NoPayments(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments, got %d", len(payments))
	}
}

func TestGetUpcomingPayments_WithinRange(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Within 30 days.
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 10), nil)
	// Beyond 30 days.
	seedCalendarSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly,
		today.AddDate(0, 0, 45), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 1 {
		t.Errorf("expected 1 payment, got %d", len(payments))
	}
	if len(payments) > 0 && payments[0].ServiceName != "Netflix" {
		t.Errorf("expected Netflix, got %s", payments[0].ServiceName)
	}
}

func TestGetUpcomingPayments_SortedByDate(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)

	seedCalendarSub(repo, userID, "Later", 5000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 20), nil)
	seedCalendarSub(repo, userID, "Sooner", 3000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 5), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 2 {
		t.Fatalf("expected 2 payments, got %d", len(payments))
	}
	if payments[0].ServiceName != "Sooner" {
		t.Errorf("expected Sooner first, got %s", payments[0].ServiceName)
	}
	if payments[1].ServiceName != "Later" {
		t.Errorf("expected Later second, got %s", payments[1].ServiceName)
	}
}

func TestGetUpcomingPayments_DaysUntilCalculation(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 7), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	if payments[0].DaysUntil != 7 {
		t.Errorf("expected daysUntil 7, got %d", payments[0].DaysUntil)
	}
}

func TestGetUpcomingPayments_DefaultDays(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	// 35 days out — should NOT be included with default 30 days.
	seedCalendarSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 35), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 0) // 0 → defaults to 30
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments with default days, got %d", len(payments))
	}
}

func TestGetUpcomingPayments_MaxDaysClamped(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	// 85 days out — within 90 max.
	seedCalendarSub(repo, userID, "Annual", 100000, models.BillingCycleYearly,
		today.AddDate(0, 0, 85), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 200) // clamped to 90
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 1 {
		t.Errorf("expected 1 payment within clamped range, got %d", len(payments))
	}
}

func TestGetUpcomingPayments_WithShareAmount(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	sub := seedCalendarSub(repo, userID, "Netflix", 20000, models.BillingCycleMonthly,
		today.AddDate(0, 0, 10), nil)

	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeEqual,
		TotalMembersSnapshot: 4,
	}

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	// 20000 / 4 = 5000
	if payments[0].PersonalAmount != 5000 {
		t.Errorf("expected personalAmount 5000, got %d", payments[0].PersonalAmount)
	}
}

func TestGetUpcomingPayments_PastBillingDateExcluded(t *testing.T) {
	repo := newMockSubRepoForCalendar()
	shareRepo := newMockShareRepoForCalendar()
	svc := NewCalendarService(repo, shareRepo)
	userID := uuid.New()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	// Past date — should be excluded.
	seedCalendarSub(repo, userID, "Old", 5000, models.BillingCycleMonthly,
		today.AddDate(0, 0, -5), nil)

	payments, err := svc.GetUpcomingPayments(userID.String(), 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments for past billing, got %d", len(payments))
	}
}
