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

// CreateFolderRequest holds the body for creating a folder.
type CreateFolderRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=50"`
	SortOrder *int   `json:"sortOrder"`
}

// UpdateFolderRequest holds the body for updating a folder.
type UpdateFolderRequest struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=50"`
	SortOrder *int    `json:"sortOrder"`
}

// FolderService handles business logic for folders.
type FolderService struct {
	repo repositories.FolderRepository
}

// NewFolderService creates a new FolderService.
func NewFolderService(repo repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

// GetFolders returns all folders for a user.
func (s *FolderService) GetFolders(userID string) ([]*models.Folder, error) {
	folders, err := s.repo.FindByUserID(userID)
	if err != nil {
		slog.Error("폴더 목록 조회 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("폴더 목록을 조회할 수 없습니다")
	}
	return folders, nil
}

// CreateFolder validates and creates a new folder.
func (s *FolderService) CreateFolder(userID string, req *CreateFolderRequest) (*models.Folder, error) {
	req.Name = strings.TrimSpace(req.Name)

	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, utils.ErrBadRequest("유효하지 않은 사용자 ID입니다")
	}

	folder := &models.Folder{
		UserID: uid,
		Name:   req.Name,
	}

	if req.SortOrder != nil {
		folder.SortOrder = *req.SortOrder
	}

	if err := s.repo.Create(folder); err != nil {
		slog.Error("폴더 생성 실패", "userID", userID, "error", err)
		return nil, utils.ErrInternal("폴더를 생성할 수 없습니다")
	}

	return folder, nil
}

// UpdateFolder validates ownership and applies partial updates to a folder.
func (s *FolderService) UpdateFolder(userID, folderID string, req *UpdateFolderRequest) (*models.Folder, error) {
	if appErr := utils.ValidateStruct(req); appErr != nil {
		return nil, appErr
	}

	folder, err := s.repo.FindByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("폴더를 찾을 수 없습니다")
		}
		slog.Error("폴더 조회 실패", "folderID", folderID, "error", err)
		return nil, utils.ErrInternal("폴더를 조회할 수 없습니다")
	}

	if folder.UserID.String() != userID {
		return nil, utils.ErrForbidden("이 폴더를 수정할 권한이 없습니다")
	}

	if req.Name != nil {
		folder.Name = strings.TrimSpace(*req.Name)
	}
	if req.SortOrder != nil {
		folder.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(folder); err != nil {
		slog.Error("폴더 수정 실패", "folderID", folderID, "error", err)
		return nil, utils.ErrInternal("폴더를 수정할 수 없습니다")
	}

	return folder, nil
}

// DeleteFolder validates ownership and deletes a folder.
func (s *FolderService) DeleteFolder(userID, folderID string) error {
	folder, err := s.repo.FindByID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("폴더를 찾을 수 없습니다")
		}
		slog.Error("폴더 조회 실패", "folderID", folderID, "error", err)
		return utils.ErrInternal("폴더를 조회할 수 없습니다")
	}

	if folder.UserID.String() != userID {
		return utils.ErrForbidden("이 폴더를 삭제할 권한이 없습니다")
	}

	if err := s.repo.Delete(folderID); err != nil {
		slog.Error("폴더 삭제 실패", "folderID", folderID, "error", err)
		return utils.ErrInternal("폴더를 삭제할 수 없습니다")
	}

	return nil
}
