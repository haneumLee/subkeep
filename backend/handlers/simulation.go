package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// SimulationHandler handles simulation-related HTTP requests.
type SimulationHandler struct {
	service *services.SimulationService
}

// NewSimulationHandler creates a new SimulationHandler.
func NewSimulationHandler(service *services.SimulationService) *SimulationHandler {
	return &SimulationHandler{service: service}
}

// SimulateCancel handles POST /api/v1/simulation/cancel.
func (h *SimulationHandler) SimulateCancel(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.CancelSimulationRequest
	if err := c.BodyParser(&req); err != nil {
		slog.Warn("시뮬레이션 요청 파싱 실패", "error", err)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	result, svcErr := h.service.SimulateCancel(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("시뮬레이션을 수행할 수 없습니다"))
	}

	return utils.Success(c, result)
}

// SimulateAdd handles POST /api/v1/simulation/add.
func (h *SimulationHandler) SimulateAdd(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.AddSimulationRequest
	if err := c.BodyParser(&req); err != nil {
		slog.Warn("시뮬레이션 요청 파싱 실패", "error", err)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	result, svcErr := h.service.SimulateAdd(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("시뮬레이션을 수행할 수 없습니다"))
	}

	return utils.Success(c, result)
}

// ApplySimulation handles POST /api/v1/simulation/apply.
func (h *SimulationHandler) ApplySimulation(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.ApplySimulationRequest
	if err := c.BodyParser(&req); err != nil {
		slog.Warn("시뮬레이션 적용 요청 파싱 실패", "error", err)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	if svcErr := h.service.ApplySimulation(userID, &req); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("시뮬레이션을 적용할 수 없습니다"))
	}

	return utils.SuccessWithMessage(c, "시뮬레이션이 적용되었습니다", nil)
}

// UndoSimulation handles POST /api/v1/simulation/undo.
func (h *SimulationHandler) UndoSimulation(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	if svcErr := h.service.UndoSimulation(userID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("실행 취소에 실패했습니다"))
	}

	return utils.SuccessWithMessage(c, "시뮬레이션이 실행 취소되었습니다", nil)
}
