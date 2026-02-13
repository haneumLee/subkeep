package services

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNotNil(t *testing.T, val interface{}) {
	t.Helper()
	if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		t.Fatalf("expected non-nil, got nil")
	}
}

func assertNil(t *testing.T, val interface{}) {
	t.Helper()
	if val != nil && !(reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		t.Errorf("expected nil, got %v", val)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func assertAppErrorCode(t *testing.T, err error, code int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected AppError with code %d, got nil", code)
	}
	appErr, ok := err.(*utils.AppError)
	if !ok {
		t.Fatalf("expected *utils.AppError, got %T: %v", err, err)
	}
	if appErr.Code != code {
		t.Errorf("expected error code %d, got %d (%s)", code, appErr.Code, appErr.Detail)
	}
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

type mockSubscriptionRepo struct {
	subs      map[string]*models.Subscription
	createErr error
	updateErr error
	deleteErr error
}

func newMockRepo() *mockSubscriptionRepo {
	return &mockSubscriptionRepo{subs: make(map[string]*models.Subscription)}
}

func (m *mockSubscriptionRepo) FindByID(id string) (*models.Subscription, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return sub, nil
}

func (m *mockSubscriptionRepo) FindByUserID(userID string, filter repositories.SubscriptionFilter) ([]*models.Subscription, int64, error) {
	var result []*models.Subscription
	for key, sub := range m.subs {
		// Skip soft-deleted entries.
		if len(key) > 8 && key[:8] == "deleted:" {
			continue
		}
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

func (m *mockSubscriptionRepo) Create(sub *models.Subscription) error {
	if m.createErr != nil {
		return m.createErr
	}
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	m.subs[sub.ID.String()] = sub
	return nil
}

func (m *mockSubscriptionRepo) Update(sub *models.Subscription) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	sub.UpdatedAt = time.Now()
	m.subs[sub.ID.String()] = sub
	return nil
}

func (m *mockSubscriptionRepo) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if sub, ok := m.subs[id]; ok {
		now := time.Now()
		sub.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
		delete(m.subs, id)
		m.subs["deleted:"+id] = sub
	}
	return nil
}

func (m *mockSubscriptionRepo) Restore(id string) error {
	key := "deleted:" + id
	sub, ok := m.subs[key]
	if !ok {
		return fmt.Errorf("subscription not found for restore: %s", id)
	}
	sub.DeletedAt = gorm.DeletedAt{Valid: false}
	delete(m.subs, key)
	m.subs[id] = sub
	return nil
}

func (m *mockSubscriptionRepo) CountByUserID(userID string) (int64, error) {
	var count int64
	for key, sub := range m.subs {
		if len(key) > 8 && key[:8] == "deleted:" {
			continue
		}
		if sub.UserID.String() == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockSubscriptionRepo) FindDuplicateName(userID, serviceName string) (bool, error) {
	for _, sub := range m.subs {
		if sub.UserID.String() == userID && sub.ServiceName == serviceName {
			return true, nil
		}
	}
	return false, nil
}

// seedSubscription inserts a subscription into the mock repo and returns it.
func (m *mockSubscriptionRepo) seedSubscription(userID uuid.UUID, name string, amount int, cycle models.BillingCycle) *models.Subscription {
	sub := &models.Subscription{
		ID:              uuid.New(),
		UserID:          userID,
		ServiceName:     name,
		Amount:          amount,
		BillingCycle:    cycle,
		Currency:        "KRW",
		NextBillingDate: time.Now().Add(30 * 24 * time.Hour),
		AutoRenew:       true,
		Status:          models.SubscriptionStatusActive,
		StartDate:       time.Now().Truncate(24 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	m.subs[sub.ID.String()] = sub
	return sub
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestGetSubscriptions(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("returns list of user subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)
		repo.seedSubscription(userID, "Spotify", 10900, models.BillingCycleMonthly)
		repo.seedSubscription(otherUserID, "YouTube", 14900, models.BillingCycleMonthly)

		subs, total, err := svc.GetSubscriptions(userID.String(), repositories.SubscriptionFilter{})
		assertNil(t, err)
		assertEqual(t, total, int64(2))
		assertEqual(t, len(subs), 2)
	})

	t.Run("returns empty list for user with no subscriptions", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		subs, total, err := svc.GetSubscriptions(uuid.New().String(), repositories.SubscriptionFilter{})
		assertNil(t, err)
		assertEqual(t, total, int64(0))
		assertEqual(t, len(subs), 0)
	})

	t.Run("respects status filter", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		activeSub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)
		activeSub.Status = models.SubscriptionStatusActive

		pausedSub := repo.seedSubscription(userID, "Spotify", 10900, models.BillingCycleMonthly)
		pausedSub.Status = models.SubscriptionStatusPaused

		subs, total, err := svc.GetSubscriptions(userID.String(), repositories.SubscriptionFilter{
			Status: "active",
		})
		assertNil(t, err)
		assertEqual(t, total, int64(1))
		assertEqual(t, len(subs), 1)
		assertEqual(t, subs[0].ServiceName, "Netflix")
	})
}

func TestGetSubscription(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("returns subscription when user is owner", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)

		got, err := svc.GetSubscription(userID.String(), sub.ID.String())
		assertNil(t, err)
		assertNotNil(t, got)
		assertEqual(t, got.ServiceName, "Netflix")
	})

	t.Run("returns ErrForbidden when user is not owner", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)

		_, err := svc.GetSubscription(otherUserID.String(), sub.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error when subscription not found", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		_, err := svc.GetSubscription(userID.String(), uuid.New().String())
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}

func TestCreateSubscription(t *testing.T) {
	userID := uuid.New()

	validReq := func() *CreateSubscriptionRequest {
		return &CreateSubscriptionRequest{
			ServiceName:     "Netflix",
			Amount:          17000,
			BillingCycle:    "monthly",
			NextBillingDate: time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
		}
	}

	t.Run("creates subscription with valid data", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.SatisfactionScore = intPtr(4)
		req.Note = strPtr("영화 보기")

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, sub)
		assertEqual(t, sub.ServiceName, "Netflix")
		assertEqual(t, sub.Amount, 17000)
		assertEqual(t, sub.BillingCycle, models.BillingCycleMonthly)
		assertEqual(t, sub.Currency, "KRW")
		assertEqual(t, sub.AutoRenew, true)
		assertEqual(t, sub.Status, models.SubscriptionStatusActive)
		assertEqual(t, *sub.SatisfactionScore, 4)
		assertEqual(t, *sub.Note, "영화 보기")
	})

	t.Run("creates subscription with minimum required fields", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		sub, err := svc.CreateSubscription(userID.String(), validReq())
		assertNil(t, err)
		assertNotNil(t, sub)
		assertEqual(t, sub.ServiceName, "Netflix")
	})

	t.Run("rejects empty service name", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.ServiceName = ""

		_, err := svc.CreateSubscription(userID.String(), req)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("rejects negative amount", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.Amount = -1

		_, err := svc.CreateSubscription(userID.String(), req)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("rejects invalid billing cycle", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.BillingCycle = "daily"

		_, err := svc.CreateSubscription(userID.String(), req)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("defaults status to active", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.Status = "" // not set

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertEqual(t, sub.Status, models.SubscriptionStatusActive)
	})

	t.Run("defaults autoRenew to true", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.AutoRenew = nil // not set

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertEqual(t, sub.AutoRenew, true)
	})

	t.Run("defaults currency to KRW", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		sub, err := svc.CreateSubscription(userID.String(), validReq())
		assertNil(t, err)
		assertEqual(t, sub.Currency, "KRW")
	})

	t.Run("sets startDate to today when not provided", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.StartDate = nil

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)

		today := time.Now().Truncate(24 * time.Hour)
		assertEqual(t, sub.StartDate.Format("2006-01-02"), today.Format("2006-01-02"))
	})

	t.Run("uses provided startDate", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.StartDate = strPtr("2025-01-15")

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertEqual(t, sub.StartDate.Format("2006-01-02"), "2025-01-15")
	})

	t.Run("allows autoRenew false", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.AutoRenew = boolPtr(false)

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertEqual(t, sub.AutoRenew, false)
	})

	t.Run("allows paused status", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		req := validReq()
		req.Status = "paused"

		sub, err := svc.CreateSubscription(userID.String(), req)
		assertNil(t, err)
		assertEqual(t, sub.Status, models.SubscriptionStatusPaused)
	})
}

func TestUpdateSubscription(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	setup := func() (*mockSubscriptionRepo, *SubscriptionService, *models.Subscription) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)
		return repo, svc, sub
	}

	t.Run("updates service name", func(t *testing.T) {
		_, svc, sub := setup()

		updated, err := svc.UpdateSubscription(userID.String(), sub.ID.String(), &UpdateSubscriptionRequest{
			ServiceName: strPtr("Netflix Premium"),
		})
		assertNil(t, err)
		assertEqual(t, updated.ServiceName, "Netflix Premium")
	})

	t.Run("updates amount", func(t *testing.T) {
		_, svc, sub := setup()

		updated, err := svc.UpdateSubscription(userID.String(), sub.ID.String(), &UpdateSubscriptionRequest{
			Amount: intPtr(20000),
		})
		assertNil(t, err)
		assertEqual(t, updated.Amount, 20000)
	})

	t.Run("partial update keeps other fields", func(t *testing.T) {
		_, svc, sub := setup()

		updated, err := svc.UpdateSubscription(userID.String(), sub.ID.String(), &UpdateSubscriptionRequest{
			ServiceName: strPtr("Netflix 4K"),
		})
		assertNil(t, err)
		assertEqual(t, updated.ServiceName, "Netflix 4K")
		// Amount should remain unchanged.
		assertEqual(t, updated.Amount, 17000)
		// BillingCycle should remain unchanged.
		assertEqual(t, updated.BillingCycle, models.BillingCycleMonthly)
		// AutoRenew should remain unchanged.
		assertEqual(t, updated.AutoRenew, true)
	})

	t.Run("returns ErrForbidden when user is not owner", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSubscription(otherUserID.String(), sub.ID.String(), &UpdateSubscriptionRequest{
			ServiceName: strPtr("Hacked"),
		})
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error when subscription not found", func(t *testing.T) {
		_, svc, _ := setup()

		_, err := svc.UpdateSubscription(userID.String(), uuid.New().String(), &UpdateSubscriptionRequest{
			ServiceName: strPtr("Ghost"),
		})
		assertAppErrorCode(t, err, http.StatusNotFound)
	})

	t.Run("rejects empty service name on update", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSubscription(userID.String(), sub.ID.String(), &UpdateSubscriptionRequest{
			ServiceName: strPtr("   "),
		})
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})
}

func TestDeleteSubscription(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("soft deletes subscription", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)

		err := svc.DeleteSubscription(userID.String(), sub.ID.String())
		assertNil(t, err)

		// Verify it's no longer in the repo.
		_, findErr := repo.FindByID(sub.ID.String())
		if findErr == nil {
			t.Error("expected subscription to be deleted from repo")
		}
	})

	t.Run("returns ErrForbidden when user is not owner", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)

		err := svc.DeleteSubscription(otherUserID.String(), sub.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)

		// Verify it's still in the repo.
		got, findErr := repo.FindByID(sub.ID.String())
		assertNil(t, findErr)
		assertNotNil(t, got)
	})

	t.Run("returns error when subscription not found", func(t *testing.T) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)

		err := svc.DeleteSubscription(userID.String(), uuid.New().String())
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}

func TestUpdateSatisfaction(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	setup := func() (*mockSubscriptionRepo, *SubscriptionService, *models.Subscription) {
		repo := newMockRepo()
		svc := NewSubscriptionService(repo)
		sub := repo.seedSubscription(userID, "Netflix", 17000, models.BillingCycleMonthly)
		return repo, svc, sub
	}

	t.Run("valid scores", func(t *testing.T) {
		for score := 1; score <= 5; score++ {
			t.Run(fmt.Sprintf("updates score to %d", score), func(t *testing.T) {
				_, svc, sub := setup()

				updated, err := svc.UpdateSatisfaction(userID.String(), sub.ID.String(), score)
				assertNil(t, err)
				assertNotNil(t, updated)
				assertNotNil(t, updated.SatisfactionScore)
				assertEqual(t, *updated.SatisfactionScore, score)
			})
		}
	})

	t.Run("rejects score 0", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSatisfaction(userID.String(), sub.ID.String(), 0)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("rejects score 6", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSatisfaction(userID.String(), sub.ID.String(), 6)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("rejects negative score", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSatisfaction(userID.String(), sub.ID.String(), -1)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("returns ErrForbidden when user is not owner", func(t *testing.T) {
		_, svc, sub := setup()

		_, err := svc.UpdateSatisfaction(otherUserID.String(), sub.ID.String(), 3)
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error when subscription not found", func(t *testing.T) {
		_, svc, _ := setup()

		_, err := svc.UpdateSatisfaction(userID.String(), uuid.New().String(), 3)
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}
