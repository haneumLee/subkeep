package services

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// CreateSubscriptionRequest holds the body for creating a subscription.
type CreateSubscriptionRequest struct {
	ServiceName       string  `json:"serviceName" validate:"required,min=1,max=100"`
	CategoryID        *string `json:"categoryId" validate:"omitempty,uuid"`
	Amount            int     `json:"amount" validate:"required,gte=0,lte=9999999"`
	BillingCycle      string  `json:"billingCycle" validate:"required,oneof=weekly monthly yearly"`
	NextBillingDate   string  `json:"nextBillingDate" validate:"required"`
	AutoRenew         *bool   `json:"autoRenew"`
	Status            string  `json:"status" validate:"omitempty,oneof=active paused"`
	SatisfactionScore *int    `json:"satisfactionScore" validate:"omitempty,min=1,max=5"`
	Note              *string `json:"note" validate:"omitempty,max=500"`
	ServiceURL        *string `json:"serviceUrl" validate:"omitempty,url,max=255"`
	StartDate         *string `json:"startDate"`
}

// UpdateSubscriptionRequest holds the body for updating a subscription.
type UpdateSubscriptionRequest struct {
	ServiceName       *string `json:"serviceName" validate:"omitempty,min=1,max=100"`
	CategoryID        *string `json:"categoryId" validate:"omitempty,uuid"`
	Amount            *int    `json:"amount" validate:"omitempty,gte=0,lte=9999999"`
	BillingCycle      *string `json:"billingCycle" validate:"omitempty,oneof=weekly monthly yearly"`
	NextBillingDate   *string `json:"nextBillingDate"`
	AutoRenew         *bool   `json:"autoRenew"`
	Status            *string `json:"status" validate:"omitempty,oneof=active paused cancelled"`
	SatisfactionScore *int    `json:"satisfactionScore" validate:"omitempty,min=1,max=5"`
	Note              *string `json:"note" validate:"omitempty,max=500"`
	ServiceURL        *string `json:"serviceUrl" validate:"omitempty,url,max=255"`
}

// SubscriptionService handles business logic for subscriptions.
type SubscriptionService struct {
	repo repositories.SubscriptionRepository
}

// NewSubscriptionService creates a new SubscriptionService.
func NewSubscriptionService(repo repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// GetSubscriptions returns a paginated, filtered list of subscriptions for the user.
func (s *SubscriptionService) GetSubscriptions(userID string, filter repositories.SubscriptionFilter) ([]*models.Subscription, int64, error) {
	subs, total, err := s.repo.FindByUserID(userID, filter)
	if err != nil {
		slog.Error("구독 목록 조회 실패", "userID", userID, "error", err)
		return nil, 0, utils.ErrInternal("구독 목록을 조회할 수 없습니다")
	}
	return subs, total, nil
}

// GetSubscription returns a single subscription after verifying ownership.
func (s *SubscriptionService) GetSubscription(userID, subID string) (*models.Subscription, error) {
	sub, err := s.repo.FindByID(subID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("구독을 찾을 수 없습니다")
		}
		slog.Error("구독 조회 실패", "subID", subID, "error", err)
		return nil, utils.ErrInternal("구독을 조회할 수 없습니다")
	}

	if sub.UserID.String() != userID {
		return nil, utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
	}

	return sub, nil
}

// CreateSubscription validates and creates a new subscription.
func (s *SubscriptionService) CreateSubscription(userID string, req *CreateSubscriptionRequest) (*models.Subscription, error) {
	// Trim service name.
	req.ServiceName = strings.TrimSpace(req.ServiceName)

	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Parse next billing date.
	nextBillingDate, err := time.Parse("2006-01-02", req.NextBillingDate)
	if err != nil {
		return nil, utils.ErrValidation("다음 결제일 형식이 올바르지 않습니다 (YYYY-MM-DD)")
	}

	// Parse start date (defaults to today).
	startDate := time.Now().Truncate(24 * time.Hour)
	if req.StartDate != nil && *req.StartDate != "" {
		parsed, parseErr := time.Parse("2006-01-02", *req.StartDate)
		if parseErr != nil {
			return nil, utils.ErrValidation("시작일 형식이 올바르지 않습니다 (YYYY-MM-DD)")
		}
		startDate = parsed
	}

	// Warn for large amounts.
	if req.Amount > 1000000 {
		slog.Warn("높은 구독 금액 입력", "userID", userID, "serviceName", req.ServiceName, "amount", req.Amount)
	}

	// Parse user ID.
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.ErrBadRequest("유효하지 않은 사용자 ID입니다")
	}

	// Parse category ID.
	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		cid, parseErr := uuid.Parse(*req.CategoryID)
		if parseErr != nil {
			return nil, utils.ErrValidation("유효하지 않은 카테고리 ID입니다")
		}
		categoryID = &cid
	}

	// Determine auto-renew (default true).
	autoRenew := true
	if req.AutoRenew != nil {
		autoRenew = *req.AutoRenew
	}

	// Determine status (default active).
	status := models.SubscriptionStatusActive
	if req.Status != "" {
		status = models.SubscriptionStatus(req.Status)
	}

	sub := &models.Subscription{
		UserID:            uid,
		ServiceName:       req.ServiceName,
		CategoryID:        categoryID,
		Amount:            req.Amount,
		BillingCycle:      models.BillingCycle(req.BillingCycle),
		Currency:          "KRW",
		NextBillingDate:   nextBillingDate,
		AutoRenew:         autoRenew,
		Status:            status,
		SatisfactionScore: req.SatisfactionScore,
		Note:              req.Note,
		ServiceURL:        req.ServiceURL,
		StartDate:         startDate,
	}

	if err := s.repo.Create(sub); err != nil {
		slog.Error("구독 생성 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("구독을 생성할 수 없습니다")
	}

	// Re-fetch to preload associations.
	created, err := s.repo.FindByID(sub.ID.String())
	if err != nil {
		return sub, nil
	}

	return created, nil
}

// UpdateSubscription validates ownership and applies partial updates.
func (s *SubscriptionService) UpdateSubscription(userID, subID string, req *UpdateSubscriptionRequest) (*models.Subscription, error) {
	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Verify ownership.
	sub, err := s.GetSubscription(userID, subID)
	if err != nil {
		return nil, err
	}

	// Apply partial updates.
	if req.ServiceName != nil {
		trimmed := strings.TrimSpace(*req.ServiceName)
		if trimmed == "" {
			return nil, utils.ErrValidation("서비스 이름은 비어있을 수 없습니다")
		}
		sub.ServiceName = trimmed
	}

	if req.CategoryID != nil {
		if *req.CategoryID == "" {
			sub.CategoryID = nil
		} else {
			cid, parseErr := uuid.Parse(*req.CategoryID)
			if parseErr != nil {
				return nil, utils.ErrValidation("유효하지 않은 카테고리 ID입니다")
			}
			sub.CategoryID = &cid
		}
	}

	if req.Amount != nil {
		sub.Amount = *req.Amount
		if *req.Amount > 1000000 {
			slog.Warn("높은 구독 금액 수정", "userID", userID, "subID", subID, "amount", *req.Amount)
		}
	}

	if req.BillingCycle != nil {
		sub.BillingCycle = models.BillingCycle(*req.BillingCycle)
	}

	if req.NextBillingDate != nil {
		parsed, parseErr := time.Parse("2006-01-02", *req.NextBillingDate)
		if parseErr != nil {
			return nil, utils.ErrValidation("다음 결제일 형식이 올바르지 않습니다 (YYYY-MM-DD)")
		}
		sub.NextBillingDate = parsed
	}

	if req.AutoRenew != nil {
		sub.AutoRenew = *req.AutoRenew
	}

	if req.Status != nil {
		sub.Status = models.SubscriptionStatus(*req.Status)
	}

	if req.SatisfactionScore != nil {
		sub.SatisfactionScore = req.SatisfactionScore
	}

	if req.Note != nil {
		sub.Note = req.Note
	}

	if req.ServiceURL != nil {
		sub.ServiceURL = req.ServiceURL
	}

	if err := s.repo.Update(sub); err != nil {
		slog.Error("구독 수정 실패", "subID", subID, "error", err)
		return nil, utils.ErrInternal("구독을 수정할 수 없습니다")
	}

	// Re-fetch to preload associations.
	updated, fetchErr := s.repo.FindByID(sub.ID.String())
	if fetchErr != nil {
		return sub, nil
	}

	return updated, nil
}

// DeleteSubscription validates ownership and soft-deletes a subscription.
func (s *SubscriptionService) DeleteSubscription(userID, subID string) error {
	// Verify ownership.
	if _, err := s.GetSubscription(userID, subID); err != nil {
		return err
	}

	if err := s.repo.Delete(subID); err != nil {
		slog.Error("구독 삭제 실패", "subID", subID, "error", err)
		return utils.ErrInternal("구독을 삭제할 수 없습니다")
	}

	return nil
}

// UpdateSatisfaction updates only the satisfaction score of a subscription.
func (s *SubscriptionService) UpdateSatisfaction(userID, subID string, score int) (*models.Subscription, error) {
	if score < 1 || score > 5 {
		return nil, utils.ErrValidation("만족도 점수는 1에서 5 사이여야 합니다")
	}

	// Verify ownership.
	sub, err := s.GetSubscription(userID, subID)
	if err != nil {
		return nil, err
	}

	sub.SatisfactionScore = &score

	if err := s.repo.Update(sub); err != nil {
		slog.Error("만족도 점수 수정 실패", "subID", subID, "error", err)
		return nil, utils.ErrInternal("만족도 점수를 수정할 수 없습니다")
	}

	// Re-fetch to preload associations.
	updated, fetchErr := s.repo.FindByID(sub.ID.String())
	if fetchErr != nil {
		return sub, nil
	}

	return updated, nil
}
