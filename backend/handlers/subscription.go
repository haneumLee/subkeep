package handlers

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// SubscriptionResponse enriches a subscription with computed monetary fields.
type SubscriptionResponse struct {
	models.Subscription
	MonthlyAmount int `json:"monthlyAmount"`
	AnnualAmount  int `json:"annualAmount"`
}

// toSubscriptionResponse converts a Subscription to a SubscriptionResponse.
func toSubscriptionResponse(sub *models.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		Subscription:  *sub,
		MonthlyAmount: sub.MonthlyAmount(),
		AnnualAmount:  sub.AnnualAmount(),
	}
}

// toSubscriptionResponses converts a slice of subscriptions.
func toSubscriptionResponses(subs []*models.Subscription) []*SubscriptionResponse {
	responses := make([]*SubscriptionResponse, len(subs))
	for i, sub := range subs {
		responses[i] = toSubscriptionResponse(sub)
	}
	return responses
}

// SubscriptionHandler handles subscription-related HTTP requests.
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler.
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// getUserID extracts the authenticated user ID from Fiber context.
func getUserID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return "", utils.ErrUnauthorized("사용자 인증 정보를 확인할 수 없습니다")
	}
	return userID, nil
}

// GetAll handles GET /api/v1/subscriptions.
// Query params: status, categoryId, sortBy, sortOrder, page, perPage.
func (h *SubscriptionHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "20"))

	filter := repositories.SubscriptionFilter{
		Status:     c.Query("status"),
		CategoryID: c.Query("categoryId"),
		FolderID:   c.Query("folderId"),
		SortBy:     c.Query("sortBy"),
		SortOrder:  c.Query("sortOrder"),
		Page:       page,
		PerPage:    perPage,
	}

	subs, total, svcErr := h.service.GetSubscriptions(userID, filter)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	responses := toSubscriptionResponses(subs)

	// Apply filter defaults for pagination meta.
	filter.Defaults()

	return utils.Paginated(c, responses, filter.Page, filter.PerPage, total)
}

// Create handles POST /api/v1/subscriptions.
func (h *SubscriptionHandler) Create(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	var req services.CreateSubscriptionRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("구독 생성 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	sub, svcErr := h.service.CreateSubscription(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Created(c, toSubscriptionResponse(sub))
}

// GetByID handles GET /api/v1/subscriptions/:id.
func (h *SubscriptionHandler) GetByID(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subID := c.Params("id")
	if subID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	sub, svcErr := h.service.GetSubscription(userID, subID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, toSubscriptionResponse(sub))
}

// Update handles PUT /api/v1/subscriptions/:id.
func (h *SubscriptionHandler) Update(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subID := c.Params("id")
	if subID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	var req services.UpdateSubscriptionRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("구독 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	sub, svcErr := h.service.UpdateSubscription(userID, subID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, toSubscriptionResponse(sub))
}

// Delete handles DELETE /api/v1/subscriptions/:id.
func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subID := c.Params("id")
	if subID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	if svcErr := h.service.DeleteSubscription(userID, subID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.NoContent(c)
}

// CheckDuplicates handles GET /api/v1/subscriptions/duplicates.
func (h *SubscriptionHandler) CheckDuplicates(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	result, svcErr := h.service.CheckDuplicates(userID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, result)
}

// satisfactionRequest holds the body for updating satisfaction score.
type satisfactionRequest struct {
	Score int `json:"score"`
}

// UpdateSatisfaction handles PATCH /api/v1/subscriptions/:id/satisfaction.
func (h *SubscriptionHandler) UpdateSatisfaction(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subID := c.Params("id")
	if subID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	var req satisfactionRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("만족도 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	sub, svcErr := h.service.UpdateSatisfaction(userID, subID, req.Score)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, toSubscriptionResponse(sub))
}
