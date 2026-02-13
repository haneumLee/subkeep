package services

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
)

// ---------------------------------------------------------------------------
// Mock repositories for report tests
// ---------------------------------------------------------------------------

type mockSubRepoForReport struct {
	subs map[string]*models.Subscription
}

func newMockSubRepoForReport() *mockSubRepoForReport {
	return &mockSubRepoForReport{subs: make(map[string]*models.Subscription)}
}

func (m *mockSubRepoForReport) FindByID(id string) (*models.Subscription, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, nil
	}
	return sub, nil
}
func (m *mockSubRepoForReport) FindByUserID(userID string, filter repositories.SubscriptionFilter) ([]*models.Subscription, int64, error) {
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
func (m *mockSubRepoForReport) Create(sub *models.Subscription) error              { return nil }
func (m *mockSubRepoForReport) Update(sub *models.Subscription) error              { return nil }
func (m *mockSubRepoForReport) Delete(id string) error                             { return nil }
func (m *mockSubRepoForReport) Restore(id string) error                            { return nil }
func (m *mockSubRepoForReport) CountByUserID(userID string) (int64, error)         { return 0, nil }
func (m *mockSubRepoForReport) FindDuplicateName(userID, name string) (bool, error) { return false, nil }
func (m *mockSubRepoForReport) FindSimilarInCategory(userID string, categoryID string, excludeSubID string) ([]*models.Subscription, error) {
	return nil, nil
}

type mockShareRepoForReport struct {
	shares map[string]*models.SubscriptionShare
}

func newMockShareRepoForReport() *mockShareRepoForReport {
	return &mockShareRepoForReport{shares: make(map[string]*models.SubscriptionShare)}
}

func (m *mockShareRepoForReport) FindByID(id string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForReport) FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForReport) FindByUserID(userID string) ([]*models.SubscriptionShare, error) {
	var result []*models.SubscriptionShare
	for _, s := range m.shares {
		result = append(result, s)
	}
	return result, nil
}
func (m *mockShareRepoForReport) Create(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForReport) Update(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForReport) Delete(id string) error                        { return nil }
func (m *mockShareRepoForReport) DeleteBySubscriptionID(subscriptionID string) error { return nil }

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func rptMakeCategory(name, color string) *models.Category {
	return &models.Category{
		ID:    uuid.New(),
		Name:  name,
		Color: &color,
	}
}

func rptIntPtr(i int) *int { return &i }

func seedReportSub(
	repo *mockSubRepoForReport,
	userID uuid.UUID,
	name string,
	amount int,
	cycle models.BillingCycle,
	status models.SubscriptionStatus,
	satisfaction *int,
	category *models.Category,
) *models.Subscription {
	sub := &models.Subscription{
		ID:                uuid.New(),
		UserID:            userID,
		ServiceName:       name,
		Amount:            amount,
		BillingCycle:      cycle,
		Currency:          "KRW",
		NextBillingDate:   time.Now().Add(30 * 24 * time.Hour),
		AutoRenew:         true,
		Status:            status,
		SatisfactionScore: satisfaction,
		StartDate:         time.Now().AddDate(-1, 0, 0),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if category != nil {
		sub.CategoryID = &category.ID
		sub.Category = category
	}
	repo.subs[sub.ID.String()] = sub
	return sub
}

// ===========================================================================
// GetOverview
// ===========================================================================

func TestGetOverview_NoSubscriptions(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.CategoryBreakdown) != 0 {
		t.Errorf("expected 0 categories, got %d", len(overview.CategoryBreakdown))
	}
	if overview.AverageCost.Monthly != 0 {
		t.Errorf("expected monthly 0, got %d", overview.AverageCost.Monthly)
	}
	if overview.Summary.TotalSubscriptions != 0 {
		t.Errorf("expected 0 total, got %d", overview.Summary.TotalSubscriptions)
	}
	if len(overview.MonthlyTrend) != 12 {
		t.Errorf("expected 12 months trend, got %d", len(overview.MonthlyTrend))
	}
}

func TestGetOverview_SingleActiveSubscription(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.AverageCost.Monthly != 17000 {
		t.Errorf("expected monthly 17000, got %d", overview.AverageCost.Monthly)
	}
	if overview.Summary.ActiveCount != 1 {
		t.Errorf("expected 1 active, got %d", overview.Summary.ActiveCount)
	}
}

func TestGetOverview_MultipleActiveSubscriptions(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	seedReportSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Monthly: 17000 + 10900 = 27900
	if overview.AverageCost.Monthly != 27900 {
		t.Errorf("expected monthly 27900, got %d", overview.AverageCost.Monthly)
	}
	// Annual: 27900 * 12 = 334800
	if overview.AverageCost.Annual != 334800 {
		t.Errorf("expected annual 334800, got %d", overview.AverageCost.Annual)
	}
}

func TestGetOverview_CategoryBreakdown_SingleCategory(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	entertainment := rptMakeCategory("Entertainment", "#FF5722")
	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
	seedReportSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.CategoryBreakdown) != 1 {
		t.Fatalf("expected 1 category, got %d", len(overview.CategoryBreakdown))
	}
	if overview.CategoryBreakdown[0].CategoryName != "Entertainment" {
		t.Errorf("expected Entertainment, got %s", overview.CategoryBreakdown[0].CategoryName)
	}
	if overview.CategoryBreakdown[0].MonthlyAmount != 31900 {
		t.Errorf("expected 31900, got %d", overview.CategoryBreakdown[0].MonthlyAmount)
	}
	if overview.CategoryBreakdown[0].Count != 2 {
		t.Errorf("expected count 2, got %d", overview.CategoryBreakdown[0].Count)
	}
	if overview.CategoryBreakdown[0].Percentage != 100.0 {
		t.Errorf("expected 100%%, got %.1f%%", overview.CategoryBreakdown[0].Percentage)
	}
}

func TestGetOverview_CategoryBreakdown_MultipleCategories(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	entertainment := rptMakeCategory("Entertainment", "#FF5722")
	music := rptMakeCategory("Music", "#2196F3")

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
	seedReportSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
	seedReportSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, music)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.CategoryBreakdown) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(overview.CategoryBreakdown))
	}
	// Sorted by amount: Entertainment (31900) first, Music (10900) second.
	if overview.CategoryBreakdown[0].CategoryName != "Entertainment" {
		t.Errorf("expected Entertainment first, got %s", overview.CategoryBreakdown[0].CategoryName)
	}
	// Entertainment: 31900/42800 ≈ 74.5%
	pct := overview.CategoryBreakdown[0].Percentage
	if pct < 74.0 || pct > 75.0 {
		t.Errorf("expected Entertainment ~74.5%%, got %.1f%%", pct)
	}
}

func TestGetOverview_CategoryBreakdown_Uncategorized(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.CategoryBreakdown) != 1 {
		t.Fatalf("expected 1 category, got %d", len(overview.CategoryBreakdown))
	}
	if overview.CategoryBreakdown[0].CategoryID != "uncategorized" {
		t.Errorf("expected 'uncategorized', got '%s'", overview.CategoryBreakdown[0].CategoryID)
	}
	if overview.CategoryBreakdown[0].Color != "#9E9E9E" {
		t.Errorf("expected '#9E9E9E', got '%s'", overview.CategoryBreakdown[0].Color)
	}
}

func TestGetOverview_CategoryBreakdown_SortByAmount(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	catA := rptMakeCategory("Cheap", "#111")
	catB := rptMakeCategory("Mid", "#222")
	catC := rptMakeCategory("Expensive", "#333")

	seedReportSub(repo, userID, "A", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catA)
	seedReportSub(repo, userID, "B", 15000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catB)
	seedReportSub(repo, userID, "C", 30000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catC)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.CategoryBreakdown) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(overview.CategoryBreakdown))
	}
	if overview.CategoryBreakdown[0].CategoryName != "Expensive" {
		t.Errorf("expected Expensive first, got %s", overview.CategoryBreakdown[0].CategoryName)
	}
	if overview.CategoryBreakdown[2].CategoryName != "Cheap" {
		t.Errorf("expected Cheap last, got %s", overview.CategoryBreakdown[2].CategoryName)
	}
}

func TestGetOverview_MonthlyTrend_TwelveMonths(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overview.MonthlyTrend) != 12 {
		t.Errorf("expected 12 months trend, got %d", len(overview.MonthlyTrend))
	}
	// All months should have amount since sub started 1 year ago.
	for i, trend := range overview.MonthlyTrend {
		if trend.Amount == 0 {
			t.Errorf("month %d: expected non-zero amount", i)
		}
	}
}

func TestGetOverview_MonthlyTrend_NewSubscription(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	// Sub started 3 months ago.
	sub := seedReportSub(repo, userID, "New", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	sub.StartDate = time.Now().AddDate(0, -3, 0)
	repo.subs[sub.ID.String()] = sub

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Earlier months should have 0, later months should have the amount.
	zeroCount := 0
	nonZeroCount := 0
	for _, trend := range overview.MonthlyTrend {
		if trend.Amount == 0 {
			zeroCount++
		} else {
			nonZeroCount++
		}
	}
	if nonZeroCount < 3 {
		t.Errorf("expected at least 3 non-zero months, got %d", nonZeroCount)
	}
}

func TestGetOverview_AverageCost_Monthly(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	seedReportSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.AverageCost.Monthly != 27900 {
		t.Errorf("expected monthly 27900, got %d", overview.AverageCost.Monthly)
	}
}

func TestGetOverview_AverageCost_WithMixedCycles(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	// monthly: 10000
	seedReportSub(repo, userID, "Monthly", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	// yearly: 120000 / 12 = 10000
	seedReportSub(repo, userID, "Yearly", 120000, models.BillingCycleYearly, models.SubscriptionStatusActive, nil, nil)
	// weekly: Round(2500 * 52 / 12) = 10833
	seedReportSub(repo, userID, "Weekly", 2500, models.BillingCycleWeekly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 10000 + 10000 + 10833 = 30833
	if overview.AverageCost.Monthly != 30833 {
		t.Errorf("expected monthly 30833, got %d", overview.AverageCost.Monthly)
	}
	// Annual: 30833 * 12 = 369996
	if overview.AverageCost.Annual != 369996 {
		t.Errorf("expected annual 369996, got %d", overview.AverageCost.Annual)
	}
	// Weekly: Round(369996 / 52) = Round(7115.31) = 7115
	expectedWeekly := int(math.Round(float64(369996) / 52.0))
	if overview.AverageCost.Weekly != expectedWeekly {
		t.Errorf("expected weekly %d, got %d", expectedWeekly, overview.AverageCost.Weekly)
	}
}

func TestGetOverview_Summary_ActiveAndPaused(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	seedReportSub(repo, userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	seedReportSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusPaused, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Summary.ActiveCount != 2 {
		t.Errorf("expected 2 active, got %d", overview.Summary.ActiveCount)
	}
	if overview.Summary.PausedCount != 1 {
		t.Errorf("expected 1 paused, got %d", overview.Summary.PausedCount)
	}
	if overview.Summary.TotalSubscriptions != 3 {
		t.Errorf("expected 3 total, got %d", overview.Summary.TotalSubscriptions)
	}
}

func TestGetOverview_Summary_MostExpensive(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
	seedReportSub(repo, userID, "Expensive", 50000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Summary.MostExpensive == nil {
		t.Fatal("expected most expensive, got nil")
	}
	if *overview.Summary.MostExpensive != "Expensive" {
		t.Errorf("expected 'Expensive', got '%s'", *overview.Summary.MostExpensive)
	}
	if overview.Summary.MostExpensiveAmount != 50000 {
		t.Errorf("expected 50000, got %d", overview.Summary.MostExpensiveAmount)
	}
}

func TestGetOverview_Summary_AverageSatisfaction(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, rptIntPtr(5), nil)
	seedReportSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusActive, rptIntPtr(3), nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// (5 + 3) / 2 = 4.0
	if overview.Summary.AverageSatisfaction != 4.0 {
		t.Errorf("expected satisfaction 4.0, got %.1f", overview.Summary.AverageSatisfaction)
	}
}

func TestGetOverview_WithShareEqual(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	sub := seedReportSub(repo, userID, "Netflix", 20000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	// Equal split 4 members → 20000 / 4 = 5000
	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeEqual,
		TotalMembersSnapshot: 4,
	}

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.AverageCost.Monthly != 5000 {
		t.Errorf("expected monthly 5000, got %d", overview.AverageCost.Monthly)
	}
}

func TestGetOverview_WithShareCustomAmount(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	sub := seedReportSub(repo, userID, "Netflix", 20000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	shareAmt := 7000
	shareRepo.shares[sub.ID.String()] = &models.SubscriptionShare{
		ID:                   uuid.New(),
		SubscriptionID:       sub.ID,
		ShareGroupID:         uuid.New(),
		SplitType:            models.SplitTypeCustomAmount,
		MyShareAmount:        &shareAmt,
		TotalMembersSnapshot: 3,
	}

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.AverageCost.Monthly != 7000 {
		t.Errorf("expected monthly 7000, got %d", overview.AverageCost.Monthly)
	}
}

func TestGetOverview_PausedNotInCategoryBreakdown(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	cat := rptMakeCategory("Entertainment", "#FF5722")
	seedReportSub(repo, userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, cat)
	seedReportSub(repo, userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusPaused, nil, cat)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Category breakdown only includes active subs.
	if len(overview.CategoryBreakdown) != 1 {
		t.Fatalf("expected 1 category entry, got %d", len(overview.CategoryBreakdown))
	}
	if overview.CategoryBreakdown[0].MonthlyAmount != 17000 {
		t.Errorf("expected 17000, got %d", overview.CategoryBreakdown[0].MonthlyAmount)
	}
	if overview.CategoryBreakdown[0].Count != 1 {
		t.Errorf("expected count 1, got %d", overview.CategoryBreakdown[0].Count)
	}
}

func TestGetOverview_WeeklyAverageCost(t *testing.T) {
	repo := newMockSubRepoForReport()
	shareRepo := newMockShareRepoForReport()
	svc := NewReportService(repo, shareRepo)
	userID := uuid.New()

	seedReportSub(repo, userID, "Netflix", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

	overview, err := svc.GetOverview(userID.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Weekly: Round(120000 / 52) = Round(2307.69) = 2308
	expectedWeekly := int(math.Round(float64(120000) / 52.0))
	if overview.AverageCost.Weekly != expectedWeekly {
		t.Errorf("expected weekly %d, got %d", expectedWeekly, overview.AverageCost.Weekly)
	}
}
