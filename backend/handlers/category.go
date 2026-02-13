package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// CategoryHandler handles category-related HTTP requests.
type CategoryHandler struct {
	service *services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll handles GET /api/v1/categories.
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	cats, svcErr := h.service.GetCategories(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, cats)
}

// Create handles POST /api/v1/categories.
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.CreateCategoryRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("카테고리 생성 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	cat, svcErr := h.service.CreateCategory(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Created(c, cat)
}

// Update handles PUT /api/v1/categories/:id.
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	categoryID := c.Params("id")
	if categoryID == "" {
		return utils.Error(c, utils.ErrBadRequest("카테고리 ID가 필요합니다"))
	}

	var req services.UpdateCategoryRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("카테고리 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	cat, svcErr := h.service.UpdateCategory(userID, categoryID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, cat)
}

// Delete handles DELETE /api/v1/categories/:id.
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	categoryID := c.Params("id")
	if categoryID == "" {
		return utils.Error(c, utils.ErrBadRequest("카테고리 ID가 필요합니다"))
	}

	if svcErr := h.service.DeleteCategory(userID, categoryID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.NoContent(c)
}
