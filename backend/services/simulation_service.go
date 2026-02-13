package services

import (
	"errors"
	"log/slog"
	"math"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// CancelSimulationRequest holds the body for cancel simulation.
type CancelSimulationRequest struct {
	SubscriptionIDs []string `json:"subscriptionIds" validate:"required,min=1"`
}

// AddSimulationRequest holds the body for add simulation.
type AddSimulationRequest struct {
	ServiceName  string  `json:"serviceName" validate:"required,min=1,max=100"`
	Amount       int     `json:"amount" validate:"required,gte=0,lte=9999999"`
	BillingCycle string  `json:"billingCycle" validate:"required,oneof=weekly monthly yearly"`
	CategoryID   *string `json:"categoryId"`
}

// ApplySimulationRequest holds the body for applying a simulation.
type ApplySimulationRequest struct {
	Action          string   `json:"action" validate:"required,oneof=cancel"`
	SubscriptionIDs []string `json:"subscriptionIds" validate:"required,min=1"`
}

// SimulationResult holds the result of a simulation.
type SimulationResult struct {
	CurrentMonthlyTotal   int                 `json:"currentMonthlyTotal"`
	SimulatedMonthlyTotal int                 `json:"simulatedMonthlyTotal"`
	MonthlyDifference     int                 `json:"monthlyDifference"`
	AnnualDifference      int                 `json:"annualDifference"`
	CategoryBreakdown     []CategoryBreakdown `json:"categoryBreakdown"`
}

// undoEntry stores the subscription IDs that were soft-deleted by ApplySimulation.
type undoEntry struct {
	subscriptionIDs []string
	expiresAt       time.Time
}

// SimulationService handles simulation-related business logic.
type SimulationService struct {
	subRepo   repositories.SubscriptionRepository
	shareRepo repositories.SubscriptionShareRepository
	undoStore map[string]*undoEntry // key: userID
	undoMu    sync.Mutex
}

// NewSimulationService creates a new SimulationService.
func NewSimulationService(subRepo repositories.SubscriptionRepository, shareRepo repositories.SubscriptionShareRepository) *SimulationService {
	return &SimulationService{
		subRepo:   subRepo,
		shareRepo: shareRepo,
		undoStore: make(map[string]*undoEntry),
	}
}

// SimulateCancel simulates cancelling the given subscriptions and returns the impact.
func (s *SimulationService) SimulateCancel(userID string, req *CancelSimulationRequest) (*SimulationResult, error) {
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Fetch all active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("시뮬레이션 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("시뮬레이션 데이터를 조회할 수 없습니다")
	}

	// Build cancel ID set.
	cancelSet := make(map[string]bool, len(req.SubscriptionIDs))
	for _, id := range req.SubscriptionIDs {
		cancelSet[id] = true
	}

	// Validate that all requested IDs belong to the user.
	foundIDs := make(map[string]bool)
	for _, sub := range activeSubs {
		if cancelSet[sub.ID.String()] {
			foundIDs[sub.ID.String()] = true
		}
	}
	for _, id := range req.SubscriptionIDs {
		if !foundIDs[id] {
			return nil, utils.ErrNotFound("구독을 찾을 수 없습니다: " + id)
		}
	}

	// Fetch subscription shares for the user.
	shareMap := buildShareMap(s.shareRepo, userID)

	// Calculate current and simulated totals.
	currentTotal := 0
	simulatedTotal := 0

	categoryMap := make(map[string]*categoryGroupData)

	for _, sub := range activeSubs {
		monthly := sub.MonthlyAmount()
		personalMonthly := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personalMonthly = share.PersonalAmount(monthly)
		}
		currentTotal += personalMonthly

		if cancelSet[sub.ID.String()] {
			continue
		}

		simulatedTotal += personalMonthly
		addToCategoryGroup(categoryMap, sub, personalMonthly)
	}

	breakdown := buildCategoryBreakdown(categoryMap, simulatedTotal)
	diff := currentTotal - simulatedTotal

	return &SimulationResult{
		CurrentMonthlyTotal:   currentTotal,
		SimulatedMonthlyTotal: simulatedTotal,
		MonthlyDifference:     diff,
		AnnualDifference:      diff * 12,
		CategoryBreakdown:     breakdown,
	}, nil
}

// SimulateAdd simulates adding a new subscription and returns the impact.
func (s *SimulationService) SimulateAdd(userID string, req *AddSimulationRequest) (*SimulationResult, error) {
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Fetch all active subscriptions.
	activeSubs, _, err := s.subRepo.FindByUserID(userID, repositories.SubscriptionFilter{
		Status:  "active",
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		slog.Error("시뮬레이션 구독 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("시뮬레이션 데이터를 조회할 수 없습니다")
	}

	// Fetch subscription shares for the user.
	shareMap := buildShareMap(s.shareRepo, userID)

	// Calculate current total and category breakdown.
	currentTotal := 0
	categoryMap := make(map[string]*categoryGroupData)

	for _, sub := range activeSubs {
		monthly := sub.MonthlyAmount()
		personalMonthly := monthly
		if share, ok := shareMap[sub.ID.String()]; ok {
			personalMonthly = share.PersonalAmount(monthly)
		}
		currentTotal += personalMonthly
		addToCategoryGroup(categoryMap, sub, personalMonthly)
	}

	// Calculate the virtual item's monthly amount.
	virtualMonthly := calcMonthlyAmount(req.Amount, models.BillingCycle(req.BillingCycle))

	// Add virtual item to category breakdown.
	catID := "uncategorized"
	catName := "미분류"
	catColor := "#9E9E9E"
	if req.CategoryID != nil && *req.CategoryID != "" {
		catID = *req.CategoryID
		// We don't have the category name from the repo here, use ID as name placeholder.
		// The category may be resolved on the frontend side.
		catName = catID
	}

	if g, ok := categoryMap[catID]; ok {
		g.amount += virtualMonthly
		g.count++
	} else {
		categoryMap[catID] = &categoryGroupData{
			categoryID:   catID,
			categoryName: catName,
			color:        catColor,
			amount:       virtualMonthly,
			count:        1,
		}
	}

	simulatedTotal := currentTotal + virtualMonthly
	diff := currentTotal - simulatedTotal // negative = cost increase

	breakdown := buildCategoryBreakdown(categoryMap, simulatedTotal)

	return &SimulationResult{
		CurrentMonthlyTotal:   currentTotal,
		SimulatedMonthlyTotal: simulatedTotal,
		MonthlyDifference:     diff,
		AnnualDifference:      diff * 12,
		CategoryBreakdown:     breakdown,
	}, nil
}

// ApplySimulation applies a simulation by actually performing the action (cancel).
func (s *SimulationService) ApplySimulation(userID string, req *ApplySimulationRequest) error {
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return appErr
	}

	// Validate all subscription IDs belong to the user.
	for _, id := range req.SubscriptionIDs {
		sub, err := s.subRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrNotFound("구독을 찾을 수 없습니다: " + id)
			}
			slog.Error("시뮬레이션 적용 구독 조회 실패", "subID", id, "error", err)
			return utils.ErrInternal("시뮬레이션을 적용할 수 없습니다")
		}
		if sub.UserID.String() != userID {
			return utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
		}
	}

	// 적용 전 상태를 undo 스토어에 저장 (30초 TTL).
	s.undoMu.Lock()
	s.undoStore[userID] = &undoEntry{
		subscriptionIDs: make([]string, len(req.SubscriptionIDs)),
		expiresAt:       time.Now().Add(30 * time.Second),
	}
	copy(s.undoStore[userID].subscriptionIDs, req.SubscriptionIDs)
	s.undoMu.Unlock()

	// Soft-delete all selected subscriptions.
	for _, id := range req.SubscriptionIDs {
		if err := s.subRepo.Delete(id); err != nil {
			slog.Error("시뮬레이션 적용 구독 삭제 실패", "subID", id, "error", err)
			return utils.ErrInternal("구독 해지에 실패했습니다: " + id)
		}
	}

	slog.Info("시뮬레이션 적용 완료", "userID", userID, "action", req.Action, "count", len(req.SubscriptionIDs))
	return nil
}

// UndoSimulation reverses the last ApplySimulation for the given user within 30 seconds.
func (s *SimulationService) UndoSimulation(userID string) error {
	s.undoMu.Lock()
	entry, ok := s.undoStore[userID]
	if ok {
		delete(s.undoStore, userID)
	}
	s.undoMu.Unlock()

	if !ok {
		return utils.ErrNotFound("실행 취소할 작업이 없습니다")
	}

	if time.Now().After(entry.expiresAt) {
		return utils.ErrBadRequest("실행 취소 기간이 만료되었습니다")
	}

	for _, id := range entry.subscriptionIDs {
		if err := s.subRepo.Restore(id); err != nil {
			slog.Error("시뮬레이션 실행 취소 복원 실패", "subID", id, "error", err)
			return utils.ErrInternal("실행 취소에 실패했습니다: " + id)
		}
	}

	slog.Info("시뮬레이션 실행 취소 완료", "userID", userID, "count", len(entry.subscriptionIDs))
	return nil
}

// categoryGroupData is a helper for grouping subscriptions by category.
type categoryGroupData struct {
	categoryID   string
	categoryName string
	color        string
	amount       int
	count        int
}

// addToCategoryGroup adds a subscription's monthly amount to the category grouping map.
func addToCategoryGroup(m map[string]*categoryGroupData, sub *models.Subscription, monthly int) {
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

	if g, ok := m[catID]; ok {
		g.amount += monthly
		g.count++
	} else {
		m[catID] = &categoryGroupData{
			categoryID:   catID,
			categoryName: catName,
			color:        catColor,
			amount:       monthly,
			count:        1,
		}
	}
}

// buildCategoryBreakdown converts the category group map to a sorted slice.
func buildCategoryBreakdown(m map[string]*categoryGroupData, total int) []CategoryBreakdown {
	breakdown := make([]CategoryBreakdown, 0, len(m))
	for _, g := range m {
		pct := 0.0
		if total > 0 {
			pct = math.Round(float64(g.amount)/float64(total)*1000) / 10
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

// calcMonthlyAmount converts an amount with a billing cycle to a monthly equivalent.
func calcMonthlyAmount(amount int, cycle models.BillingCycle) int {
	switch cycle {
	case models.BillingCycleMonthly:
		return amount
	case models.BillingCycleYearly:
		return int(math.Round(float64(amount) / 12.0))
	case models.BillingCycleWeekly:
		return int(math.Round(float64(amount) * 52.0 / 12.0))
	default:
		return amount
	}
}
