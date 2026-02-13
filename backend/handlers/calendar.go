package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// CalendarHandler handles calendar-related HTTP requests.
type CalendarHandler struct {
	service *services.CalendarService
}

// NewCalendarHandler creates a new CalendarHandler.
func NewCalendarHandler(service *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

// GetMonthlyCalendar handles GET /api/v1/calendar/monthly.
func (h *CalendarHandler) GetMonthlyCalendar(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Parse optional year parameter.
	if yearStr := c.Query("year"); yearStr != "" {
		parsed, parseErr := strconv.Atoi(yearStr)
		if parseErr != nil {
			return utils.Error(c, utils.ErrBadRequest("year는 숫자여야 합니다"))
		}
		year = parsed
	}

	// Parse optional month parameter.
	if monthStr := c.Query("month"); monthStr != "" {
		parsed, parseErr := strconv.Atoi(monthStr)
		if parseErr != nil {
			return utils.Error(c, utils.ErrBadRequest("month는 숫자여야 합니다"))
		}
		month = parsed
	}

	// Validate ranges.
	if year < 2000 || year > 2100 {
		return utils.Error(c, utils.ErrBadRequest("year는 2000~2100 범위여야 합니다"))
	}
	if month < 1 || month > 12 {
		return utils.Error(c, utils.ErrBadRequest("month는 1~12 범위여야 합니다"))
	}

	calendar, svcErr := h.service.GetMonthlyCalendar(userID, year, month)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("캘린더 데이터를 조회할 수 없습니다"))
	}

	return utils.Success(c, calendar)
}
