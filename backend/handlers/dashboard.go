package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// DashboardHandler handles dashboard-related HTTP requests.
type DashboardHandler struct {
	service *services.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetSummary handles GET /api/v1/dashboard/summary.
func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	summary, svcErr := h.service.GetSummary(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("대시보드 요약을 조회할 수 없습니다"))
	}

	return utils.Success(c, summary)
}

// GetRecommendations handles GET /api/v1/dashboard/recommendations.
func (h *DashboardHandler) GetRecommendations(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	recommendations, svcErr := h.service.GetRecommendations(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("해지 추천을 조회할 수 없습니다"))
	}

	return utils.Success(c, recommendations)
}
