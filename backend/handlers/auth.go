package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService  *services.AuthService
	oauthService *services.OAuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *services.AuthService, oauthService *services.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// oauthCallbackRequest represents the body of an OAuth callback request.
type oauthCallbackRequest struct {
	Code        string `json:"code" validate:"required"`
	RedirectURI string `json:"redirectUri" validate:"required,url"`
}

// refreshTokenRequest represents the body of a refresh token request.
type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// OAuthCallback handles POST /api/v1/auth/:provider/callback.
// It exchanges the authorization code for user info and returns tokens.
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	if provider == "" {
		return utils.Error(c, utils.ErrBadRequest("인증 제공자가 지정되지 않았습니다"))
	}

	var req oauthCallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, utils.ErrBadRequest("요청 본문을 파싱할 수 없습니다"))
	}

	if req.Code == "" {
		return utils.Error(c, utils.ErrBadRequest("인증 코드가 필요합니다"))
	}
	if req.RedirectURI == "" {
		return utils.Error(c, utils.ErrBadRequest("리다이렉트 URI가 필요합니다"))
	}

	// Exchange code for OAuth user info.
	oauthUser, err := h.oauthService.ExchangeCode(provider, req.Code, req.RedirectURI)
	if err != nil {
		slog.Error("oauth exchange failed", "error", err, "provider", provider)
		if appErr, ok := err.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("OAuth 인증에 실패했습니다"))
	}

	// Handle OAuth callback (find/create user + generate tokens).
	tokens, user, err := h.authService.HandleOAuthCallback(provider, oauthUser)
	if err != nil {
		slog.Error("oauth callback handling failed", "error", err, "provider", provider)
		if appErr, ok := err.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrInternal("인증 처리에 실패했습니다"))
	}

	// Set refresh token as httpOnly secure cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth",
		MaxAge:   7 * 24 * 60 * 60, // 7 days in seconds
	})

	return utils.Success(c, fiber.Map{
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"nickname":  user.Nickname,
			"avatarUrl": user.AvatarURL,
			"provider":  user.Provider,
			"createdAt": user.CreatedAt,
		},
		"tokens": tokens,
	})
}

// RefreshToken handles POST /api/v1/auth/refresh.
// It validates the refresh token and returns a new token pair.
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var refreshToken string

	// Try body first.
	var req refreshTokenRequest
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	}

	// Fall back to cookie.
	if refreshToken == "" {
		refreshToken = c.Cookies("refresh_token")
	}

	if refreshToken == "" {
		return utils.Error(c, utils.ErrBadRequest("리프레시 토큰이 필요합니다"))
	}

	tokens, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return utils.Error(c, appErr)
		}
		return utils.Error(c, utils.ErrUnauthorized("토큰 갱신에 실패했습니다"))
	}

	// Update refresh token cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return utils.Success(c, tokens)
}

// Logout handles POST /api/v1/auth/logout.
// It clears the authentication cookies.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)

	if err := h.authService.Logout(userID); err != nil {
		slog.Error("logout failed", "error", err, "userID", userID)
	}

	// Clear refresh token cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	// Clear access token cookie if set.
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	return utils.NoContent(c)
}

// GetMe handles GET /api/v1/auth/me.
// It returns the authenticated user's profile.
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.Error(c, utils.ErrUnauthorized("인증이 필요합니다"))
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		slog.Error("failed to get user", "error", err, "userID", userID)
		return utils.Error(c, utils.ErrNotFound("사용자를 찾을 수 없습니다"))
	}

	return utils.Success(c, fiber.Map{
		"id":        user.ID,
		"email":     user.Email,
		"nickname":  user.Nickname,
		"avatarUrl": user.AvatarURL,
		"provider":  user.Provider,
		"createdAt": user.CreatedAt,
	})
}
