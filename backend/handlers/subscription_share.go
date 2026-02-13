package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// SubscriptionShareHandler handles subscription share-related HTTP requests.
type SubscriptionShareHandler struct {
	service *services.SubscriptionShareService
}

// NewSubscriptionShareHandler creates a new SubscriptionShareHandler.
func NewSubscriptionShareHandler(service *services.SubscriptionShareService) *SubscriptionShareHandler {
	return &SubscriptionShareHandler{service: service}
}

// Link handles POST /api/v1/subscriptions/:id/share.
func (h *SubscriptionShareHandler) Link(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subscriptionID := c.Params("id")
	if subscriptionID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	var req services.LinkShareRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("구독 공유 연결 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	// Override subscriptionID from URL path.
	req.SubscriptionID = subscriptionID

	share, svcErr := h.service.LinkSubscriptionToShareGroup(userID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Created(c, share)
}

// Update handles PUT /api/v1/subscription-shares/:id.
func (h *SubscriptionShareHandler) Update(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	shareID := c.Params("id")
	if shareID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 공유 ID가 필요합니다"))
	}

	var req services.UpdateShareRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		slog.Debug("구독 공유 수정 요청 파싱 실패", "error", parseErr)
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	share, svcErr := h.service.UpdateSubscriptionShare(userID, shareID, &req)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, share)
}

// Unlink handles DELETE /api/v1/subscription-shares/:id.
func (h *SubscriptionShareHandler) Unlink(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	shareID := c.Params("id")
	if shareID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 공유 ID가 필요합니다"))
	}

	if svcErr := h.service.UnlinkSubscriptionShare(userID, shareID); svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.NoContent(c)
}

// GetBySubscription handles GET /api/v1/subscriptions/:id/share.
func (h *SubscriptionShareHandler) GetBySubscription(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.Error(c, err.(*utils.AppError))
	}

	subscriptionID := c.Params("id")
	if subscriptionID == "" {
		return utils.Error(c, utils.ErrBadRequest("구독 ID가 필요합니다"))
	}

	share, svcErr := h.service.GetSubscriptionShare(userID, subscriptionID)
	if svcErr != nil {
		if appErr, ok := svcErr.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal(""))
	}

	return utils.Success(c, share)
}
