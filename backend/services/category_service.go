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

// CreateCategoryRequest holds the body for creating a category.
type CreateCategoryRequest struct {
	Name      string  `json:"name" validate:"required,min=1,max=50"`
	Color     *string `json:"color" validate:"omitempty,hexcolor"`
	Icon      *string `json:"icon" validate:"omitempty,max=50"`
	SortOrder *int    `json:"sortOrder"`
}

// UpdateCategoryRequest holds the body for updating a category.
type UpdateCategoryRequest struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=50"`
	Color     *string `json:"color" validate:"omitempty,hexcolor"`
	Icon      *string `json:"icon" validate:"omitempty,max=50"`
	SortOrder *int    `json:"sortOrder"`
}

// CategoryService handles business logic for categories.
type CategoryService struct {
	repo repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetCategories returns system categories and user's custom categories.
func (s *CategoryService) GetCategories(userID string) ([]*models.Category, error) {
	cats, err := s.repo.FindByUserID(userID)
	if err != nil {
		slog.Error("카테고리 목록 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("카테고리 목록을 조회할 수 없습니다")
	}
	return cats, nil
}

// CreateCategory validates and creates a new custom category.
func (s *CategoryService) CreateCategory(userID string, req *CreateCategoryRequest) (*models.Category, error) {
	// Trim name.
	req.Name = strings.TrimSpace(req.Name)

	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Parse user ID.
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.ErrBadRequest("유효하지 않은 사용자 ID입니다")
	}

	cat := &models.Category{
		UserID:   &uid,
		Name:     req.Name,
		Color:    req.Color,
		Icon:     req.Icon,
		IsSystem: false,
	}

	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}

	if err := s.repo.Create(cat); err != nil {
		slog.Error("카테고리 생성 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("카테고리를 생성할 수 없습니다")
	}

	return cat, nil
}

// UpdateCategory validates ownership and applies partial updates to a category.
func (s *CategoryService) UpdateCategory(userID, categoryID string, req *UpdateCategoryRequest) (*models.Category, error) {
	// Validate request struct.
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	// Find category.
	cat, err := s.repo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("카테고리를 찾을 수 없습니다")
		}
		slog.Error("카테고리 조회 실패", "categoryID", categoryID, "error", err)
		return nil, utils.ErrInternal("카테고리를 조회할 수 없습니다")
	}

	// System categories cannot be modified.
	if cat.IsSystem {
		return nil, utils.ErrForbidden("시스템 카테고리는 수정할 수 없습니다")
	}

	// Verify ownership.
	if cat.UserID == nil || cat.UserID.String() != userID {
		return nil, utils.ErrForbidden("해당 카테고리에 대한 접근 권한이 없습니다")
	}

	// Apply partial updates.
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return nil, utils.ErrValidation("카테고리 이름은 비어있을 수 없습니다")
		}
		cat.Name = trimmed
	}

	if req.Color != nil {
		cat.Color = req.Color
	}

	if req.Icon != nil {
		cat.Icon = req.Icon
	}

	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(cat); err != nil {
		slog.Error("카테고리 수정 실패", "categoryID", categoryID, "error", err)
		return nil, utils.ErrInternal("카테고리를 수정할 수 없습니다")
	}

	return cat, nil
}

// DeleteCategory validates ownership and deletes a custom category.
// Subscriptions in the deleted category are reassigned to the system "기타" category.
func (s *CategoryService) DeleteCategory(userID, categoryID string) error {
	// Find category.
	cat, err := s.repo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("카테고리를 찾을 수 없습니다")
		}
		slog.Error("카테고리 조회 실패", "categoryID", categoryID, "error", err)
		return utils.ErrInternal("카테고리를 조회할 수 없습니다")
	}

	// System categories cannot be deleted.
	if cat.IsSystem {
		return utils.ErrForbidden("시스템 카테고리는 삭제할 수 없습니다")
	}

	// Verify ownership.
	if cat.UserID == nil || cat.UserID.String() != userID {
		return utils.ErrForbidden("해당 카테고리에 대한 접근 권한이 없습니다")
	}

	// Find system "기타" category for reassignment.
	systemCats, err := s.repo.FindSystemCategories()
	if err != nil {
		slog.Error("시스템 카테고리 조회 실패", "error", err)
		return utils.ErrInternal("카테고리를 삭제할 수 없습니다")
	}

	var targetCategoryID string
	for _, sc := range systemCats {
		if sc.Name == "기타" {
			targetCategoryID = sc.ID.String()
			break
		}
	}

	// Reassign subscriptions if target category exists.
	if targetCategoryID != "" {
		if err := s.repo.ReassignSubscriptions(categoryID, targetCategoryID); err != nil {
			slog.Error("구독 재배정 실패", "categoryID", categoryID, "targetCategoryID", targetCategoryID, "error", err)
			return utils.ErrInternal("카테고리를 삭제할 수 없습니다")
		}
	}

	if err := s.repo.Delete(categoryID); err != nil {
		slog.Error("카테고리 삭제 실패", "categoryID", categoryID, "error", err)
		return utils.ErrInternal("카테고리를 삭제할 수 없습니다")
	}

	return nil
}
