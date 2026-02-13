package services

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// LinkShareRequest holds the body for linking a subscription to a share group.
type LinkShareRequest struct {
	SubscriptionID string   `json:"subscriptionId" validate:"required,uuid"`
	ShareGroupID   string   `json:"shareGroupId" validate:"required,uuid"`
	SplitType      string   `json:"splitType" validate:"required,oneof=equal custom_amount custom_ratio"`
	MyShareAmount  *int     `json:"myShareAmount" validate:"omitempty,gte=0"`
	MyShareRatio   *float64 `json:"myShareRatio" validate:"omitempty,gte=0,lte=1"`
}

// UpdateShareRequest holds the body for updating a subscription share.
type UpdateShareRequest struct {
	SplitType     *string  `json:"splitType" validate:"omitempty,oneof=equal custom_amount custom_ratio"`
	MyShareAmount *int     `json:"myShareAmount" validate:"omitempty,gte=0"`
	MyShareRatio  *float64 `json:"myShareRatio" validate:"omitempty,gte=0,lte=1"`
}

// SubscriptionShareService handles business logic for subscription shares.
type SubscriptionShareService struct {
	shareRepo      repositories.SubscriptionShareRepository
	subRepo        repositories.SubscriptionRepository
	shareGroupRepo repositories.ShareGroupRepository
}

// NewSubscriptionShareService creates a new SubscriptionShareService.
func NewSubscriptionShareService(
	shareRepo repositories.SubscriptionShareRepository,
	subRepo repositories.SubscriptionRepository,
	shareGroupRepo repositories.ShareGroupRepository,
) *SubscriptionShareService {
	return &SubscriptionShareService{
		shareRepo:      shareRepo,
		subRepo:        subRepo,
		shareGroupRepo: shareGroupRepo,
	}
}

// LinkSubscriptionToShareGroup links a subscription to a share group with a split configuration.
func (s *SubscriptionShareService) LinkSubscriptionToShareGroup(userID string, req *LinkShareRequest) (*models.SubscriptionShare, error) {
	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Verify subscription ownership.
	sub, err := s.subRepo.FindByID(req.SubscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("구독을 찾을 수 없습니다")
		}
		slog.Error("구독 조회 실패", "subscriptionID", req.SubscriptionID, "error", err)
		return nil, utils.ErrInternal("구독을 조회할 수 없습니다")
	}
	if sub.UserID.String() != userID {
		return nil, utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
	}

	// Verify share group ownership.
	group, err := s.shareGroupRepo.FindByID(req.ShareGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("공유 그룹을 찾을 수 없습니다")
		}
		slog.Error("공유 그룹 조회 실패", "shareGroupID", req.ShareGroupID, "error", err)
		return nil, utils.ErrInternal("공유 그룹을 조회할 수 없습니다")
	}
	if group.OwnerUserID.String() != userID {
		return nil, utils.ErrForbidden("해당 공유 그룹에 대한 접근 권한이 없습니다")
	}

	// Check for duplicate link.
	existing, err := s.shareRepo.FindBySubscriptionID(req.SubscriptionID)
	if err == nil && existing != nil {
		return nil, utils.ErrBadRequest("이미 공유 그룹이 연결된 구독입니다")
	}

	// Validate split type fields.
	splitType := models.SplitType(req.SplitType)
	if appErr := s.validateSplitFields(splitType, req.MyShareAmount, req.MyShareRatio); appErr != nil {
		return nil, appErr
	}

	// Compute total members snapshot.
	totalMembers := len(group.Members)
	if totalMembers == 0 {
		totalMembers = 1
	}

	// Parse UUIDs.
	subID, _ := uuid.Parse(req.SubscriptionID)
	groupID, _ := uuid.Parse(req.ShareGroupID)

	share := &models.SubscriptionShare{
		SubscriptionID:       subID,
		ShareGroupID:         groupID,
		SplitType:            splitType,
		MyShareAmount:        req.MyShareAmount,
		MyShareRatio:         req.MyShareRatio,
		TotalMembersSnapshot: totalMembers,
	}

	if err := s.shareRepo.Create(share); err != nil {
		slog.Error("구독 공유 생성 실패", "subscriptionID", req.SubscriptionID, "error", err)
		return nil, utils.ErrInternal("구독 공유를 생성할 수 없습니다")
	}

	// Re-fetch to preload associations.
	created, err := s.shareRepo.FindByID(share.ID.String())
	if err != nil {
		return share, nil
	}

	return created, nil
}

// UpdateSubscriptionShare updates the split configuration of a subscription share.
func (s *SubscriptionShareService) UpdateSubscriptionShare(userID string, shareID string, req *UpdateShareRequest) (*models.SubscriptionShare, error) {
	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Find existing share.
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("구독 공유를 찾을 수 없습니다")
		}
		slog.Error("구독 공유 조회 실패", "shareID", shareID, "error", err)
		return nil, utils.ErrInternal("구독 공유를 조회할 수 없습니다")
	}

	// Verify subscription ownership.
	sub, err := s.subRepo.FindByID(share.SubscriptionID.String())
	if err != nil {
		slog.Error("구독 조회 실패", "subscriptionID", share.SubscriptionID, "error", err)
		return nil, utils.ErrInternal("구독을 조회할 수 없습니다")
	}
	if sub.UserID.String() != userID {
		return nil, utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
	}

	// Determine the effective split type.
	effectiveSplitType := share.SplitType
	if req.SplitType != nil {
		effectiveSplitType = models.SplitType(*req.SplitType)
	}

	// Determine effective share amount and ratio.
	effectiveAmount := share.MyShareAmount
	if req.MyShareAmount != nil {
		effectiveAmount = req.MyShareAmount
	}
	effectiveRatio := share.MyShareRatio
	if req.MyShareRatio != nil {
		effectiveRatio = req.MyShareRatio
	}

	// Validate split type fields.
	if appErr := s.validateSplitFields(effectiveSplitType, effectiveAmount, effectiveRatio); appErr != nil {
		return nil, appErr
	}

	// Apply updates.
	if req.SplitType != nil {
		share.SplitType = models.SplitType(*req.SplitType)
	}
	if req.MyShareAmount != nil {
		share.MyShareAmount = req.MyShareAmount
	}
	if req.MyShareRatio != nil {
		share.MyShareRatio = req.MyShareRatio
	}

	if err := s.shareRepo.Update(share); err != nil {
		slog.Error("구독 공유 수정 실패", "shareID", shareID, "error", err)
		return nil, utils.ErrInternal("구독 공유를 수정할 수 없습니다")
	}

	// Re-fetch to preload associations.
	updated, fetchErr := s.shareRepo.FindByID(share.ID.String())
	if fetchErr != nil {
		return share, nil
	}

	return updated, nil
}

// UnlinkSubscriptionShare removes the share link from a subscription.
func (s *SubscriptionShareService) UnlinkSubscriptionShare(userID string, shareID string) error {
	// Find existing share.
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("구독 공유를 찾을 수 없습니다")
		}
		slog.Error("구독 공유 조회 실패", "shareID", shareID, "error", err)
		return utils.ErrInternal("구독 공유를 조회할 수 없습니다")
	}

	// Verify subscription ownership.
	sub, err := s.subRepo.FindByID(share.SubscriptionID.String())
	if err != nil {
		slog.Error("구독 조회 실패", "subscriptionID", share.SubscriptionID, "error", err)
		return utils.ErrInternal("구독을 조회할 수 없습니다")
	}
	if sub.UserID.String() != userID {
		return utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
	}

	if err := s.shareRepo.Delete(shareID); err != nil {
		slog.Error("구독 공유 삭제 실패", "shareID", shareID, "error", err)
		return utils.ErrInternal("구독 공유를 삭제할 수 없습니다")
	}

	return nil
}

// GetSubscriptionShare returns the share info for a given subscription.
func (s *SubscriptionShareService) GetSubscriptionShare(userID string, subscriptionID string) (*models.SubscriptionShare, error) {
	// Verify subscription ownership.
	sub, err := s.subRepo.FindByID(subscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("구독을 찾을 수 없습니다")
		}
		slog.Error("구독 조회 실패", "subscriptionID", subscriptionID, "error", err)
		return nil, utils.ErrInternal("구독을 조회할 수 없습니다")
	}
	if sub.UserID.String() != userID {
		return nil, utils.ErrForbidden("해당 구독에 대한 접근 권한이 없습니다")
	}

	share, err := s.shareRepo.FindBySubscriptionID(subscriptionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("구독 공유 정보를 찾을 수 없습니다")
		}
		slog.Error("구독 공유 조회 실패", "subscriptionID", subscriptionID, "error", err)
		return nil, utils.ErrInternal("구독 공유 정보를 조회할 수 없습니다")
	}

	return share, nil
}

// validateSplitFields validates that the required fields are present for the split type.
func (s *SubscriptionShareService) validateSplitFields(splitType models.SplitType, amount *int, ratio *float64) *utils.AppError {
	switch splitType {
	case models.SplitTypeCustomAmount:
		if amount == nil {
			return utils.ErrValidation("custom_amount 분담 방식에는 금액(myShareAmount)이 필수입니다")
		}
	case models.SplitTypeCustomRatio:
		if ratio == nil {
			return utils.ErrValidation("custom_ratio 분담 방식에는 비율(myShareRatio)이 필수입니다")
		}
	case models.SplitTypeEqual:
		// No additional fields required for equal split.
	}
	return nil
}
