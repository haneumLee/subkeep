package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
)

// CreateShareGroupRequest holds the body for creating a share group.
type CreateShareGroupRequest struct {
	Name        string                     `json:"name" validate:"required,min=1,max=100"`
	Description *string                    `json:"description"`
	Members     []CreateShareMemberRequest `json:"members" validate:"required,min=1,dive"`
}

// CreateShareMemberRequest holds a member entry for group creation.
type CreateShareMemberRequest struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
}

// UpdateShareGroupRequest holds the body for updating a share group.
type UpdateShareGroupRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
}

// ShareGroupService handles business logic for share groups.
type ShareGroupService struct {
	repo repositories.ShareGroupRepository
}

// NewShareGroupService creates a new ShareGroupService.
func NewShareGroupService(repo repositories.ShareGroupRepository) *ShareGroupService {
	return &ShareGroupService{repo: repo}
}

// GetShareGroups returns all share groups owned by the user.
func (s *ShareGroupService) GetShareGroups(userID string) ([]*models.ShareGroup, error) {
	groups, err := s.repo.FindByOwnerID(userID)
	if err != nil {
		slog.Error("공유 그룹 목록 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("공유 그룹 목록을 조회할 수 없습니다")
	}
	return groups, nil
}

// GetShareGroup returns a single share group after verifying ownership.
func (s *ShareGroupService) GetShareGroup(userID, groupID string) (*models.ShareGroup, error) {
	group, err := s.repo.FindByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("공유 그룹을 찾을 수 없습니다")
		}
		slog.Error("공유 그룹 조회 실패", "groupID", groupID, "error", err)
		return nil, utils.ErrInternal("공유 그룹을 조회할 수 없습니다")
	}

	if group.OwnerUserID.String() != userID {
		return nil, utils.ErrForbidden("해당 공유 그룹에 대한 접근 권한이 없습니다")
	}

	return group, nil
}

// CreateShareGroup validates and creates a new share group with members.
// The owner is automatically added as a member with isOwner=true.
// Total members (owner + provided members) must be >= 2.
func (s *ShareGroupService) CreateShareGroup(userID string, req *CreateShareGroupRequest) (*models.ShareGroup, error) {
	// Trim name.
	req.Name = strings.TrimSpace(req.Name)

	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Total members = owner + provided members must be >= 2.
	totalMembers := 1 + len(req.Members) // 1 for owner
	if totalMembers < 2 {
		return nil, utils.ErrValidation("공유 그룹은 최소 2명의 멤버가 필요합니다")
	}

	// Parse user ID.
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.ErrBadRequest("유효하지 않은 사용자 ID입니다")
	}

	// Build member list: owner first, then provided members.
	members := make([]models.ShareMember, 0, totalMembers)
	members = append(members, models.ShareMember{
		Nickname: "나",
		IsOwner:  true,
	})

	for _, m := range req.Members {
		nickname := strings.TrimSpace(m.Nickname)
		if nickname == "" {
			return nil, utils.ErrValidation("멤버 닉네임은 비어있을 수 없습니다")
		}
		members = append(members, models.ShareMember{
			Nickname: nickname,
			IsOwner:  false,
		})
	}

	group := &models.ShareGroup{
		OwnerUserID: uid,
		Name:        req.Name,
		Description: req.Description,
		Members:     members,
	}

	if err := s.repo.Create(group); err != nil {
		slog.Error("공유 그룹 생성 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("공유 그룹을 생성할 수 없습니다")
	}

	// Re-fetch to preload associations.
	created, err := s.repo.FindByID(group.ID.String())
	if err != nil {
		return group, nil
	}

	return created, nil
}

// UpdateShareGroup validates ownership and applies partial updates.
func (s *ShareGroupService) UpdateShareGroup(userID, groupID string, req *UpdateShareGroupRequest) (*models.ShareGroup, error) {
	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Verify ownership.
	group, err := s.GetShareGroup(userID, groupID)
	if err != nil {
		return nil, err
	}

	// Apply partial updates.
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return nil, utils.ErrValidation("공유 그룹 이름은 비어있을 수 없습니다")
		}
		group.Name = trimmed
	}

	if req.Description != nil {
		group.Description = req.Description
	}

	if err := s.repo.Update(group); err != nil {
		slog.Error("공유 그룹 수정 실패", "groupID", groupID, "error", err)
		return nil, utils.ErrInternal("공유 그룹을 수정할 수 없습니다")
	}

	// Re-fetch to preload associations.
	updated, fetchErr := s.repo.FindByID(group.ID.String())
	if fetchErr != nil {
		return group, nil
	}

	return updated, nil
}

// DeleteShareGroup validates ownership, removes subscription shares, and
// soft-deletes the share group.
func (s *ShareGroupService) DeleteShareGroup(userID, groupID string) error {
	// Verify ownership.
	if _, err := s.GetShareGroup(userID, groupID); err != nil {
		return err
	}

	// Remove all subscription shares referencing this group.
	if err := s.repo.RemoveAllSubscriptionShares(groupID); err != nil {
		slog.Error("구독 공유 삭제 실패", "groupID", groupID, "error", err)
		return utils.ErrInternal("공유 그룹을 삭제할 수 없습니다")
	}

	if err := s.repo.Delete(groupID); err != nil {
		slog.Error("공유 그룹 삭제 실패", "groupID", groupID, "error", err)
		return utils.ErrInternal("공유 그룹을 삭제할 수 없습니다")
	}

	return nil
}
