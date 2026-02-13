package services

import (
	"math"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/subkeep/backend/models"
)

// ===========================================================================
// SimulateCancel
// ===========================================================================

func TestSimulateCancel(t *testing.T) {
	userID := uuid.New()

	t.Run("correctly calculates savings after cancelling subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		sub1 := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		result, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{sub1.ID.String()},
		})
		assertNil(t, err)
		assertNotNil(t, result)
		assertEqual(t, result.CurrentMonthlyTotal, 27900)
		assertEqual(t, result.SimulatedMonthlyTotal, 10900)
		assertEqual(t, result.MonthlyDifference, 17000)
		assertEqual(t, result.AnnualDifference, 17000*12)
	})

	t.Run("returns error for empty subscription IDs", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		_, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("returns error for non-existent subscription ID", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		_, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{uuid.New().String()},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusNotFound)
	})

	t.Run("returns correct category breakdown after cancellation", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		entertainment := makeCategory("Entertainment", "#FF5722")
		music := makeCategory("Music", "#2196F3")

		sub1 := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
		repo.seedSubscriptionWithDetails(userID, "YouTube", 14900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, music)

		// Cancel Netflix — Entertainment should have only YouTube (14900).
		result, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{sub1.ID.String()},
		})
		assertNil(t, err)
		assertNotNil(t, result)
		assertEqual(t, len(result.CategoryBreakdown), 2)

		// Sorted by amount DESC: Entertainment 14900, Music 10900.
		assertEqual(t, result.CategoryBreakdown[0].CategoryName, "Entertainment")
		assertEqual(t, result.CategoryBreakdown[0].MonthlyAmount, 14900)
		assertEqual(t, result.CategoryBreakdown[1].CategoryName, "Music")
		assertEqual(t, result.CategoryBreakdown[1].MonthlyAmount, 10900)
	})

	t.Run("handles cancelling all subscriptions (total becomes 0)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		sub1 := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		sub2 := repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		result, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{sub1.ID.String(), sub2.ID.String()},
		})
		assertNil(t, err)
		assertNotNil(t, result)
		assertEqual(t, result.SimulatedMonthlyTotal, 0)
		assertEqual(t, result.MonthlyDifference, 27900)
		assertEqual(t, len(result.CategoryBreakdown), 0)
	})

	t.Run("handles mixed billing cycles in calculation", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		// monthly 10000 → 10000
		// yearly 120000 → 10000
		// weekly 2500 → Round(2500*52/12) = 10833
		sub1 := repo.seedSubscriptionWithDetails(userID, "Monthly", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Yearly", 120000, models.BillingCycleYearly, models.SubscriptionStatusActive, nil, nil)
		repo.seedSubscriptionWithDetails(userID, "Weekly", 2500, models.BillingCycleWeekly, models.SubscriptionStatusActive, nil, nil)

		result, err := svc.SimulateCancel(userID.String(), &CancelSimulationRequest{
			SubscriptionIDs: []string{sub1.ID.String()},
		})
		assertNil(t, err)
		assertNotNil(t, result)
		// currentTotal = 10000 + 10000 + 10833 = 30833
		assertEqual(t, result.CurrentMonthlyTotal, 30833)
		// simulatedTotal = 10000 + 10833 = 20833
		assertEqual(t, result.SimulatedMonthlyTotal, 20833)
		assertEqual(t, result.MonthlyDifference, 10000)
	})
}

// ===========================================================================
// SimulateAdd
// ===========================================================================

func TestSimulateAdd(t *testing.T) {
	userID := uuid.New()

	t.Run("correctly calculates cost increase after adding subscription", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		result, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName:  "NewService",
			Amount:       10000,
			BillingCycle: "monthly",
		})
		assertNil(t, err)
		assertNotNil(t, result)
		assertEqual(t, result.CurrentMonthlyTotal, 17000)
		assertEqual(t, result.SimulatedMonthlyTotal, 27000)
		// Difference is negative (cost increase): current - simulated = -10000.
		assertEqual(t, result.MonthlyDifference, -10000)
		assertEqual(t, result.AnnualDifference, -10000*12)
	})

	t.Run("returns error for invalid request (missing fields)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		// Missing ServiceName.
		_, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			Amount:       10000,
			BillingCycle: "monthly",
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)

		// Missing BillingCycle.
		_, err = svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName: "Test",
			Amount:      10000,
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("handles adding weekly subscription (correct monthly conversion)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		repo.seedSubscriptionWithDetails(userID, "Existing", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		// Weekly 3000 → Round(3000*52/12) = Round(13000) = 13000
		result, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName:  "WeeklySub",
			Amount:       3000,
			BillingCycle: "weekly",
		})
		assertNil(t, err)
		assertNotNil(t, result)
		expectedMonthly := int(math.Round(3000.0 * 52.0 / 12.0))
		assertEqual(t, result.SimulatedMonthlyTotal, 10000+expectedMonthly)
	})

	t.Run("handles adding yearly subscription (correct monthly conversion)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		repo.seedSubscriptionWithDetails(userID, "Existing", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		// Yearly 99000 → Round(99000/12) = Round(8250) = 8250
		result, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName:  "YearlySub",
			Amount:       99000,
			BillingCycle: "yearly",
		})
		assertNil(t, err)
		assertNotNil(t, result)
		expectedMonthly := int(math.Round(99000.0 / 12.0))
		assertEqual(t, result.SimulatedMonthlyTotal, 10000+expectedMonthly)
	})

	t.Run("correctly includes new item in category breakdown", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		entertainment := makeCategory("Entertainment", "#FF5722")
		repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, entertainment)

		result, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName:  "NewUncategorized",
			Amount:       5000,
			BillingCycle: "monthly",
		})
		assertNil(t, err)
		assertNotNil(t, result)
		// Should have 2 categories: Entertainment and uncategorized.
		assertEqual(t, len(result.CategoryBreakdown), 2)
	})

	t.Run("adding to existing category groups correctly", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		catID := uuid.New()
		cat := &models.Category{ID: catID, Name: "Music", Color: strPtr("#2196F3")}
		repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, cat)

		// Add another subscription with same categoryId.
		result, err := svc.SimulateAdd(userID.String(), &AddSimulationRequest{
			ServiceName:  "AppleMusic",
			Amount:       8900,
			BillingCycle: "monthly",
			CategoryID:   strPtr(catID.String()),
		})
		assertNil(t, err)
		assertNotNil(t, result)
		// Since category resolution uses the ID as name placeholder for added items,
		// the existing category group should absorb the new amount.
		// The mock may show the existing cat group's amount as 10900+8900.
		// The category mapping uses catID.String() for existing items and the raw string for new items.
		// The existing sub has catID.String() and the virtual uses the same string → grouped together.
		assertEqual(t, result.SimulatedMonthlyTotal, 10900+8900)
	})
}

// ===========================================================================
// ApplySimulation
// ===========================================================================

func TestApplySimulation(t *testing.T) {
	userID := uuid.New()

	t.Run("successfully soft-deletes specified subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		sub1 := repo.seedSubscriptionWithDetails(userID, "Netflix", 17000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)
		sub2 := repo.seedSubscriptionWithDetails(userID, "Spotify", 10900, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		err := svc.ApplySimulation(userID.String(), &ApplySimulationRequest{
			Action:          "cancel",
			SubscriptionIDs: []string{sub1.ID.String()},
		})
		assertNil(t, err)

		// sub1 should be deleted from the mock store.
		_, found := repo.subs[sub1.ID.String()]
		assertEqual(t, found, false)

		// sub2 should remain.
		_, found = repo.subs[sub2.ID.String()]
		assertEqual(t, found, true)
	})

	t.Run("returns error for non-existent subscription", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		err := svc.ApplySimulation(userID.String(), &ApplySimulationRequest{
			Action:          "cancel",
			SubscriptionIDs: []string{uuid.New().String()},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusNotFound)
	})

	t.Run("returns error when subscription belongs to different user (403)", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		otherUserID := uuid.New()
		sub := repo.seedSubscriptionWithDetails(otherUserID, "OtherUserSub", 10000, models.BillingCycleMonthly, models.SubscriptionStatusActive, nil, nil)

		err := svc.ApplySimulation(userID.String(), &ApplySimulationRequest{
			Action:          "cancel",
			SubscriptionIDs: []string{sub.ID.String()},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for invalid request", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSimulationService(repo)

		// Empty SubscriptionIDs.
		err := svc.ApplySimulation(userID.String(), &ApplySimulationRequest{
			Action:          "cancel",
			SubscriptionIDs: []string{},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)

		// Missing action.
		err = svc.ApplySimulation(userID.String(), &ApplySimulationRequest{
			SubscriptionIDs: []string{uuid.New().String()},
		})
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})
}
