package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// ReportHandler handles report-related HTTP requests.
type ReportHandler struct {
	service *services.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetOverview handles GET /api/v1/reports/overview.
func (h *ReportHandler) GetOverview(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	overview, svcErr := h.service.GetOverview(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("리포트를 조회할 수 없습니다"))
	}

	return utils.Success(c, overview)
}
