package services

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/subkeep/backend/config"
	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/utils"
	"gorm.io/gorm"
)

// OAuthUser represents user information obtained from an OAuth provider.
type OAuthUser struct {
	Provider       string
	ProviderUserID string
	Email          string
	Nickname       string
	AvatarURL      string
}

// TokenPair holds the access and refresh JWT tokens.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// TokenClaims holds the custom JWT claims.
type TokenClaims struct {
	UserID    string `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

// AuthService handles authentication logic.
type AuthService struct {
	userRepo  repositories.UserRepository
	jwtConfig config.JWTConfig
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repositories.UserRepository, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

// HandleOAuthCallback processes an OAuth callback by finding or creating a user
// and generating a token pair.
func (s *AuthService) HandleOAuthCallback(provider string, oauthUser *OAuthUser) (*TokenPair, *models.User, error) {
	if !isValidProvider(provider) {
		return nil, nil, utils.ErrBadRequest("지원하지 않는 인증 제공자입니다: " + provider)
	}

	// Try to find existing user by provider ID.
	user, err := s.userRepo.FindByProviderID(provider, oauthUser.ProviderUserID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("failed to find user by provider id", "error", err, "provider", provider)
			return nil, nil, fmt.Errorf("find user by provider id: %w", err)
		}

		// User not found — create a new one.
		user = &models.User{
			Provider:       models.AuthProvider(provider),
			ProviderUserID: oauthUser.ProviderUserID,
		}
		if oauthUser.Email != "" {
			user.Email = &oauthUser.Email
		}
		if oauthUser.Nickname != "" {
			user.Nickname = &oauthUser.Nickname
		}
		if oauthUser.AvatarURL != "" {
			user.AvatarURL = &oauthUser.AvatarURL
		}

		if createErr := s.userRepo.Create(user); createErr != nil {
			slog.Error("failed to create user", "error", createErr, "provider", provider)
			return nil, nil, fmt.Errorf("create user: %w", createErr)
		}
		slog.Info("new user created via OAuth", "userID", user.ID, "provider", provider)
	} else {
		// Update existing user profile from OAuth data.
		updated := false
		if oauthUser.Email != "" && (user.Email == nil || *user.Email != oauthUser.Email) {
			user.Email = &oauthUser.Email
			updated = true
		}
		if oauthUser.Nickname != "" && (user.Nickname == nil || *user.Nickname != oauthUser.Nickname) {
			user.Nickname = &oauthUser.Nickname
			updated = true
		}
		if oauthUser.AvatarURL != "" && (user.AvatarURL == nil || *user.AvatarURL != oauthUser.AvatarURL) {
			user.AvatarURL = &oauthUser.AvatarURL
			updated = true
		}
		if updated {
			if updateErr := s.userRepo.Update(user); updateErr != nil {
				slog.Error("failed to update user profile", "error", updateErr, "userID", user.ID)
				return nil, nil, fmt.Errorf("update user profile: %w", updateErr)
			}
		}
	}

	// Update last login timestamp.
	if loginErr := s.userRepo.UpdateLastLogin(user.ID.String()); loginErr != nil {
		slog.Warn("failed to update last login", "error", loginErr, "userID", user.ID)
	}

	// Generate token pair.
	tokens, err := s.GenerateTokenPair(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, user, nil
}

// RefreshToken validates a refresh token and generates a new token pair.
func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, utils.ErrUnauthorized("유효하지 않은 리프레시 토큰입니다")
	}

	// Verify the user still exists.
	if _, findErr := s.userRepo.FindByID(claims.UserID); findErr != nil {
		slog.Warn("refresh token for non-existent user", "userID", claims.UserID)
		return nil, utils.ErrUnauthorized("사용자를 찾을 수 없습니다")
	}

	tokens, err := s.GenerateTokenPair(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, nil
}

// Logout handles user logout. Placeholder for future token blacklisting.
func (s *AuthService) Logout(userID string) error {
	slog.Info("user logged out", "userID", userID)
	// TODO: Add token blacklisting when Redis is integrated.
	return nil
}

// GetUserByID retrieves a user by their ID.
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// GenerateTokenPair creates a new access token (1h) and refresh token (7d).
func (s *AuthService) GenerateTokenPair(userID string) (*TokenPair, error) {
	now := time.Now()

	// Access token.
	accessExpiry := s.jwtConfig.Expiration
	if accessExpiry == 0 {
		accessExpiry = 1 * time.Hour
	}
	accessClaims := &TokenClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiry)),
			Issuer:    "subkeep",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Refresh token.
	refreshExpiry := s.jwtConfig.RefreshExpiration
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour
	}
	refreshClaims := &TokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiry)),
			Issuer:    "subkeep",
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshTokenObj.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}, nil
}

// ValidateAccessToken parses and validates an access token string.
func (s *AuthService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type: expected access, got %s", claims.TokenType)
	}
	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token string.
func (s *AuthService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected refresh, got %s", claims.TokenType)
	}
	return claims, nil
}

// parseToken parses and validates a JWT token string.
func (s *AuthService) parseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// isValidProvider checks whether the given provider is supported.
func isValidProvider(provider string) bool {
	switch provider {
	case "google", "apple", "naver", "kakao":
		return true
	default:
		return false
	}
}
