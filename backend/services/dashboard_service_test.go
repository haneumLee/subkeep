package services

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/subkeep/backend/models"
)

// suppress unused import warning
var _ = time.Now

// ---------------------------------------------------------------------------
// Mock SubscriptionShareRepository for dashboard tests
// ---------------------------------------------------------------------------

type mockShareRepoForDashboard struct {
	shares map[string]*models.SubscriptionShare
}

func newMockShareRepo() *mockShareRepoForDashboard {
	return &mockShareRepoForDashboard{
		shares: make(map[string]*models.SubscriptionShare),
	}
}

func (m *mockShareRepoForDashboard) FindByID(id string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForDashboard) FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error) {
	return nil, nil
}
func (m *mockShareRepoForDashboard) FindByUserID(userID string) ([]*models.SubscriptionShare, error) {
	var result []*models.SubscriptionShare
	for _, s := range m.shares {
		result = append(result, s)
	}
	return result, nil
}
func (m *mockShareRepoForDashboard) Create(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForDashboard) Update(share *models.SubscriptionShare) error  { return nil }
func (m *mockShareRepoForDashboard) Delete(id string) error                        { return nil }
func (m *mockShareRepoForDashboard) DeleteBySubscriptionID(subscriptionID string) error {
	return nil
}

// ---------------------------------------------------------------------------
// Helper: seed subscription with category and satisfaction score
// ---------------------------------------------------------------------------

func (m *mockSubscriptionRepo) seedSubscriptionWithDetails(
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
		StartDate:         time.Now().Truncate(24 * time.Hour),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if category != nil {
		sub.CategoryID = &category.ID
		sub.Category = category
	}
	m.subs[sub.ID.String()] = sub
	return sub
}

// ---------------------------------------------------------------------------
// Test categories
// ---------------------------------------------------------------------------

func makeCategory(name, color string) *models.Category {
	c := &models.Category{
		ID:    uuid.New(),
		Name:  name,
		Color: &color,
	}
	return c
}

// ===========================================================================
// GetSummary
// ===========================================================================

func TestGetSummary(t *testing.T) {
	userID := uuid.New()

	t.Run("returns correct monthly and annual totals for active subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.MonthlyTotal, 27900)
		assertEqual(t, summary.AnnualTotal, 27900*12)
	})

	t.Run("correctly counts active and paused subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusPaused, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Notion", 5000, models.BillingCycleMonthly, models.SubscriptionStatusPaused, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.ActiveCount, 2)
		assertEqual(t, summary.PausedCount, 2)
	})

	t.Run("correctly groups subscriptions by category with percentage", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		entertainment := makeCategory("Entertainment", "#FF5722")
		music := makeCategory("Music", "#2196F3")

		// Entertainment: 17000 + 14900 = 31900
		// Music: 10900
		// Total: 42800
		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
		repo.seedSubscriptionWithDetails(userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, music)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, len(summary.CategoryBreakdown), 2)

		// Sorted by amount descending: Entertainment first.
		assertEqual(t, summary.CategoryBreakdown[0].CategoryName, "Entertainment")
		assertEqual(t, summary.CategoryBreakdown[0].MonthlyAmount, 31900)
		assertEqual(t, summary.CategoryBreakdown[0].Count, 2)

		assertEqual(t, summary.CategoryBreakdown[1].CategoryName, "Music")
		assertEqual(t, summary.CategoryBreakdown[1].MonthlyAmount, 10900)
		assertEqual(t, summary.CategoryBreakdown[1].Count, 1)

		// Percentages: 31900/42800 ≈ 74.5%, 10900/42800 ≈ 25.5%
		if summary.CategoryBreakdown[0].Percentage < 74.0 || summary.CategoryBreakdown[0].Percentage > 75.0 {
			t.Errorf("expected Entertainment percentage ~74.5%%, got %.1f%%", summary.CategoryBreakdown[0].Percentage)
		}
		if summary.CategoryBreakdown[1].Percentage < 25.0 || summary.CategoryBreakdown[1].Percentage > 26.0 {
			t.Errorf("expected Music percentage ~25.5%%, got %.1f%%", summary.CategoryBreakdown[1].Percentage)
		}
	})

	t.Run("handles zero subscriptions (empty state)", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.MonthlyTotal, 0)
		assertEqual(t, summary.AnnualTotal, 0)
		assertEqual(t, summary.ActiveCount, 0)
		assertEqual(t, summary.PausedCount, 0)
		assertEqual(t, len(summary.CategoryBreakdown), 0)
	})

	t.Run("handles subscriptions without categories (uncategorized)", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, len(summary.CategoryBreakdown), 1)
		assertEqual(t, summary.CategoryBreakdown[0].CategoryID, "uncategorized")
		assertEqual(t, summary.CategoryBreakdown[0].MonthlyAmount, 27900)
		assertEqual(t, summary.CategoryBreakdown[0].Count, 2)
		assertEqual(t, summary.CategoryBreakdown[0].Color, "#9E9E9E")
	})

	t.Run("handles mixed billing cycles (weekly, monthly, yearly)", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// monthly: 10000 → 10000
		// yearly: 120000 → 120000/12 = 10000
		// weekly: 2500 → Round(2500*52/12) = Round(10833.33) = 10833
		repo.seedSubscriptionWithDetails(userID, "Monthly", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Yearly", 120000, models.BillingCycleYearly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Weekly", 2500, models.BillingCycleWeekly, models.SubscriptionStatusActive, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		// 10000 + 10000 + 10833 = 30833
		assertEqual(t, summary.MonthlyTotal, 30833)
		assertEqual(t, summary.AnnualTotal, 30833*12)
	})

	t.Run("only includes active subscriptions in total calculation", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusPaused, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusCancelled, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		// Only Netflix (active) counts toward the total.
		assertEqual(t, summary.MonthlyTotal, 17000)
		assertEqual(t, summary.ActiveCount, 1)
	})

	t.Run("sorts category breakdown by amount descending", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		catA := makeCategory("Cheap", "#111111")
		catB := makeCategory("Mid", "#222222")
		catC := makeCategory("Expensive", "#333333")

		repo.seedSubscriptionWithDetails(userID, "A", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catA)
		repo.seedSubscriptionWithDetails(userID, "B", 15000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catB)
		repo.seedSubscriptionWithDetails(userID, "C", 30000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, catC)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, len(summary.CategoryBreakdown), 3)
		assertEqual(t, summary.CategoryBreakdown[0].CategoryName, "Expensive")
		assertEqual(t, summary.CategoryBreakdown[1].CategoryName, "Mid")
		assertEqual(t, summary.CategoryBreakdown[2].CategoryName, "Cheap")
	})
}

// ===========================================================================
// GetRecommendations
// ===========================================================================

func TestGetRecommendations(t *testing.T) {
	userID := uuid.New()

	t.Run("returns empty list when no subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)
		assertNotNil(t, recs)
		assertEqual(t, len(recs), 0)
	})

	t.Run("recommends subscriptions with satisfaction score 1-2", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "BadService", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(1), nil)
		repo.seedSubscriptionWithDetails(userID, "MehService", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(2), nil)
		repo.seedSubscriptionWithDetails(userID, "GoodService", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(4), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)
		assertEqual(t, len(recs), 2)

		// Both should have reason "만족도 낮음".
		for _, rec := range recs {
			if rec.ServiceName != "BadService" && rec.ServiceName != "MehService" {
				t.Errorf("unexpected recommendation: %s", rec.ServiceName)
			}
		}
	})

	t.Run("recommends high-cost low-satisfaction subscriptions (top 20% cost + satisfaction <= 3)", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// 10 subscriptions: top 20% = top 2 by cost.
		// Satisfaction 3 with high cost → should be recommended.
		for i := 0; i < 8; i++ {
			repo.seedSubscriptionWithDetails(userID, "Cheap"+string(rune('A'+i)), 1000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(5), nil)
		}
		expensiveSub := repo.seedSubscriptionWithDetails(userID, "ExpensiveMeh", 50000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(3), nil)
		repo.seedSubscriptionWithDetails(userID, "ExpensiveGood", 40000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(4), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)

		// ExpensiveMeh has satisfaction 3 and is in top 20% cost → recommended.
		// ExpensiveGood has satisfaction 4 → not recommended.
		found := false
		for _, rec := range recs {
			if rec.SubscriptionID == expensiveSub.ID.String() {
				found = true
				assertEqual(t, rec.Reason, "높은 비용 대비 낮은 만족도")
			}
		}
		if !found {
			t.Error("expected ExpensiveMeh to be recommended")
		}
	})

	t.Run("does NOT recommend satisfaction 4-5 subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "Great", 50000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(4), nil)
		repo.seedSubscriptionWithDetails(userID, "Excellent", 90000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(5), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)
		assertEqual(t, len(recs), 0)
	})

	t.Run("sorts recommendations by satisfaction ASC then amount DESC", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// Create enough subscriptions so cost threshold logic works.
		// We want all of these to be recommended, so use satisfaction <= 2.
		repo.seedSubscriptionWithDetails(userID, "Sat2_Low", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(2), nil)
		repo.seedSubscriptionWithDetails(userID, "Sat1_High", 30000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(1), nil)
		repo.seedSubscriptionWithDetails(userID, "Sat1_Low", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(1), nil)
		repo.seedSubscriptionWithDetails(userID, "Sat2_High", 20000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(2), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)
		assertEqual(t, len(recs), 4)

		// Satisfaction ASC: satisfaction 1 first, then satisfaction 2.
		// Within same satisfaction: amount DESC.
		assertEqual(t, recs[0].ServiceName, "Sat1_High")
		assertEqual(t, recs[1].ServiceName, "Sat1_Low")
		assertEqual(t, recs[2].ServiceName, "Sat2_High")
		assertEqual(t, recs[3].ServiceName, "Sat2_Low")
	})

	t.Run("handles subscriptions with nil satisfaction score", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "NoScore", 50000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "BadScore", 5000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(1), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)

		// NoScore (nil) should NOT be recommended; only BadScore.
		assertEqual(t, len(recs), 1)
		assertEqual(t, recs[0].ServiceName, "BadScore")
	})

	t.Run("returns annual saving correctly (monthly * 12)", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		repo.seedSubscriptionWithDetails(userID, "BadSub", 15000, models.BillingCycleMonthly, models.SubscriptionStatusActive, intPtr(1), nil)

		recs, err := svc.GetRecommendations(userID.String())
		assertNil(t, err)
		assertEqual(t, len(recs), 1)
		assertEqual(t, recs[0].MonthlyAmount, 15000)
		assertEqual(t, recs[0].AnnualSaving, 15000*12)
	})
}

// ===========================================================================
// GetSummary with SubscriptionShare
// ===========================================================================

func TestGetSummaryWithShares(t *testing.T) {
	userID := uuid.New()

	t.Run("equal split reduces monthly total by member count", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// Netflix 17000/month, shared equally among 4 members → 17000/4 = 4250
		sub := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		// Spotify 10900/month, no share → 10900
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		shareRepo.shares[uuid.New().String()] = &models.SubscriptionShare{
			ID:                   uuid.New(),
			SubscriptionID:       sub.ID,
			SplitType:            models.SplitTypeEqual,
			TotalMembersSnapshot: 4,
		}

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		// 4250 + 10900 = 15150
		assertEqual(t, summary.MonthlyTotal, 15150)
		assertEqual(t, summary.AnnualTotal, 15150*12)
	})

	t.Run("custom_amount uses myShareAmount for monthly total", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// Netflix 17000/month, custom_amount = 5000
		sub := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		shareRepo.shares[uuid.New().String()] = &models.SubscriptionShare{
			ID:                   uuid.New(),
			SubscriptionID:       sub.ID,
			SplitType:            models.SplitTypeCustomAmount,
			MyShareAmount:        intPtr(5000),
			TotalMembersSnapshot: 3,
		}

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.MonthlyTotal, 5000)
	})

	t.Run("custom_ratio uses monthlyAmount * ratio for monthly total", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// Netflix 20000/month, custom_ratio = 0.3 → 20000 * 0.3 = 6000
		sub := repo.seedSubscriptionWithDetails(userID, "Netflix", 20000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		ratio := 0.3
		shareRepo.shares[uuid.New().String()] = &models.SubscriptionShare{
			ID:                   uuid.New(),
			SubscriptionID:       sub.ID,
			SplitType:            models.SplitTypeCustomRatio,
			MyShareRatio:         &ratio,
			TotalMembersSnapshot: 3,
		}

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.MonthlyTotal, 6000)
	})

	t.Run("subscriptions without share use full amount", func(t *testing.T) {
		repo := newMockRepo()
		shareRepo := newMockShareRepo()
		svc := NewDashboardService(repo, shareRepo)

		// No shares configured — full amount should be used
		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		summary, err := svc.GetSummary(userID.String())
		assertNil(t, err)
		assertNotNil(t, summary)
		assertEqual(t, summary.MonthlyTotal, 27900)
	})
}
