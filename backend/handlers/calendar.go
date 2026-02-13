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

// GetDayDetail handles GET /api/v1/calendar/daily.
func (h *CalendarHandler) GetDayDetail(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	// Parse required year parameter.
	yearStr := c.Query("year")
	if yearStr == "" {
		return utils.Error(c, utils.ErrBadRequest("year는 필수 파라미터입니다"))
	}
	year, parseErr := strconv.Atoi(yearStr)
	if parseErr != nil {
		return utils.Error(c, utils.ErrBadRequest("year는 숫자여야 합니다"))
	}

	// Parse required month parameter.
	monthStr := c.Query("month")
	if monthStr == "" {
		return utils.Error(c, utils.ErrBadRequest("month는 필수 파라미터입니다"))
	}
	month, parseErr := strconv.Atoi(monthStr)
	if parseErr != nil {
		return utils.Error(c, utils.ErrBadRequest("month는 숫자여야 합니다"))
	}

	// Parse required day parameter.
	dayStr := c.Query("day")
	if dayStr == "" {
		return utils.Error(c, utils.ErrBadRequest("day는 필수 파라미터입니다"))
	}
	day, parseErr := strconv.Atoi(dayStr)
	if parseErr != nil {
		return utils.Error(c, utils.ErrBadRequest("day는 숫자여야 합니다"))
	}

	// Validate ranges.
	if year < 2000 || year > 2100 {
		return utils.Error(c, utils.ErrBadRequest("year는 2000~2100 범위여야 합니다"))
	}
	if month < 1 || month > 12 {
		return utils.Error(c, utils.ErrBadRequest("month는 1~12 범위여야 합니다"))
	}
	if day < 1 || day > 31 {
		return utils.Error(c, utils.ErrBadRequest("day는 1~31 범위여야 합니다"))
	}

	detail, svcErr := h.service.GetDayDetail(userID, year, month, day)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("일별 결제 데이터를 조회할 수 없습니다"))
	}

	return utils.Success(c, detail)
}

// GetUpcomingPayments handles GET /api/v1/calendar/upcoming.
func (h *CalendarHandler) GetUpcomingPayments(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		parsed, parseErr := strconv.Atoi(daysStr)
		if parseErr != nil {
			return utils.Error(c, utils.ErrBadRequest("days는 숫자여야 합니다"))
		}
		if parsed < 1 || parsed > 90 {
			return utils.Error(c, utils.ErrBadRequest("days는 1~90 범위여야 합니다"))
		}
		days = parsed
	}

	payments, svcErr := h.service.GetUpcomingPayments(userID, days)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("예정 결제 데이터를 조회할 수 없습니다"))
	}

	return utils.Success(c, payments)
}
