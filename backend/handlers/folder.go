package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// FolderHandler handles folder-related HTTP requests.
type FolderHandler struct {
	service *services.FolderService
}

// NewFolderHandler creates a new FolderHandler.
func NewFolderHandler(service *services.FolderService) *FolderHandler {
	return &FolderHandler{service: service}
}

// GetAll handles GET /api/v1/folders.
func (h *FolderHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	folders, svcErr := h.service.GetFolders(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, folders)
}

// Create handles POST /api/v1/folders.
func (h *FolderHandler) Create(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.CreateFolderRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("폴더 생성 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	folder, svcErr := h.service.CreateFolder(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Created(c, folder)
}

// Update handles PUT /api/v1/folders/:id.
func (h *FolderHandler) Update(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	folderID := c.Params("id")
	if folderID == "" {
		return utils.Error(c, utils.ErrBadRequest("폴더 ID가 필요합니다"))
	}

	var req services.UpdateFolderRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("폴더 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	folder, svcErr := h.service.UpdateFolder(userID, folderID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, folder)
}

// Delete handles DELETE /api/v1/folders/:id.
func (h *FolderHandler) Delete(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	folderID := c.Params("id")
	if folderID == "" {
		return utils.Error(c, utils.ErrBadRequest("폴더 ID가 필요합니다"))
	}

	if svcErr := h.service.DeleteFolder(userID, folderID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.NoContent(c)
}
