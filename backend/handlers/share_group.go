package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// ShareGroupHandler handles share group-related HTTP requests.
type ShareGroupHandler struct {
	service *services.ShareGroupService
}

// NewShareGroupHandler creates a new ShareGroupHandler.
func NewShareGroupHandler(service *services.ShareGroupService) *ShareGroupHandler {
	return &ShareGroupHandler{service: service}
}

// GetAll handles GET /api/v1/share-groups.
func (h *ShareGroupHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	groups, svcErr := h.service.GetShareGroups(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, groups)
}

// GetByID handles GET /api/v1/share-groups/:id.
func (h *ShareGroupHandler) GetByID(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	groupID := c.Params("id")
	if groupID == "" {
		return utils.Error(c, utils.ErrBadRequest("공유 그룹 ID가 필요합니다"))
	}

	group, svcErr := h.service.GetShareGroup(userID, groupID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, group)
}

// Create handles POST /api/v1/share-groups.
func (h *ShareGroupHandler) Create(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.CreateShareGroupRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("공유 그룹 생성 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	group, svcErr := h.service.CreateShareGroup(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Created(c, group)
}

// Update handles PUT /api/v1/share-groups/:id.
func (h *ShareGroupHandler) Update(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	groupID := c.Params("id")
	if groupID == "" {
		return utils.Error(c, utils.ErrBadRequest("공유 그룹 ID가 필요합니다"))
	}

	var req services.UpdateShareGroupRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("공유 그룹 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	group, svcErr := h.service.UpdateShareGroup(userID, groupID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, group)
}

// Delete handles DELETE /api/v1/share-groups/:id.
func (h *ShareGroupHandler) Delete(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	groupID := c.Params("id")
	if groupID == "" {
		return utils.Error(c, utils.ErrBadRequest("공유 그룹 ID가 필요합니다"))
	}

	if svcErr := h.service.DeleteShareGroup(userID, groupID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.NoContent(c)
}
