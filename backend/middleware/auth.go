package middleware

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

// AuthMiddleware validates the JWT access token and sets the userID in Locals.
// Rejects requests without a valid token.
func AuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c)
		if tokenString == "" {
			return utils.Error(c, utils.ErrUnauthorized("인증 토큰이 필요합니다"))
		}

		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			slog.Debug("invalid access token", "error", err, "path", c.Path())
			return utils.Error(c, utils.ErrUnauthorized("유효하지 않은 인증 토큰입니다"))
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

// OptionalAuthMiddleware attempts to validate the JWT access token and sets
// the userID in Locals if present. Does not reject requests without a token.
func OptionalAuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c)
		if tokenString == "" {
			return c.Next()
		}

		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			slog.Debug("invalid access token in optional auth", "error", err, "path", c.Path())
			return c.Next()
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

// extractToken retrieves the JWT token from the Authorization header or
// the access_token cookie.
func extractToken(c *fiber.Ctx) string {
	// Try Authorization header first.
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			token := strings.TrimSpace(parts[1])
			if token != "" {
				return token
			}
		}
	}

	// Fall back to cookie.
	cookie := c.Cookies("access_token")
	if cookie != "" {
		return cookie
	}

	return ""
}
