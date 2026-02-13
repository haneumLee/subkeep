package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/config"
	"github.com/subkeep/backend/models"
)

const testJWTSecret = "test-secret-key-for-unit-tests"

// ---------------------------------------------------------------------------
// Mock UserRepository
// ---------------------------------------------------------------------------

type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) FindByID(id string) (*models.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) FindByProviderID(provider, providerUserID string) (*models.User, error) {
	for _, u := range m.users {
		if string(u.Provider) == provider && u.ProviderUserID == providerUserID {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email != nil && *u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Create(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	m.users[user.ID.String()] = user
	return nil
}

func (m *mockUserRepo) Update(user *models.User) error {
	m.users[user.ID.String()] = user
	return nil
}

func (m *mockUserRepo) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) UpdateLastLogin(id string) error {
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestAuthService() (*AuthService, *mockUserRepo) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, config.JWTConfig{
		Secret:            testJWTSecret,
		Expiration:        15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})
	return svc, repo
}

// ---------------------------------------------------------------------------
// Tests – GenerateTokenPair
// ---------------------------------------------------------------------------

func TestGenerateTokenPair(t *testing.T) {
	svc, _ := newTestAuthService()
	userID := uuid.New().String()

	pair, err := svc.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("AccessToken is empty")
	}
	if pair.RefreshToken == "" {
		t.Error("RefreshToken is empty")
	}
	if pair.ExpiresIn <= 0 {
		t.Errorf("ExpiresIn = %d, want > 0", pair.ExpiresIn)
	}
}

// ---------------------------------------------------------------------------
// Tests – ValidateAccessToken
// ---------------------------------------------------------------------------

func TestValidateAccessToken(t *testing.T) {
	svc, _ := newTestAuthService()
	userID := uuid.New().String()

	t.Run("accepts valid access token", func(t *testing.T) {
		pair, err := svc.GenerateTokenPair(userID)
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}
		claims, err := svc.ValidateAccessToken(pair.AccessToken)
		if err != nil {
			t.Fatalf("ValidateAccessToken() error = %v", err)
		}
		if claims.UserID != userID {
			t.Errorf("UserID = %q, want %q", claims.UserID, userID)
		}
		if claims.TokenType != "access" {
			t.Errorf("TokenType = %q, want %q", claims.TokenType, "access")
		}
	})

	t.Run("rejects expired token", func(t *testing.T) {
		// Create a token that is already expired.
		now := time.Now()
		claims := &TokenClaims{
			UserID:    userID,
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
				Issuer:    "subkeep",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testJWTSecret))
		if err != nil {
			t.Fatalf("failed to sign expired token: %v", err)
		}

		_, err = svc.ValidateAccessToken(tokenString)
		if err == nil {
			t.Error("ValidateAccessToken() should reject expired token")
		}
	})

	t.Run("rejects refresh token as access token", func(t *testing.T) {
		pair, err := svc.GenerateTokenPair(userID)
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}
		_, err = svc.ValidateAccessToken(pair.RefreshToken)
		if err == nil {
			t.Error("ValidateAccessToken() should reject refresh token")
		}
	})

	t.Run("rejects malformed token", func(t *testing.T) {
		_, err := svc.ValidateAccessToken("not-a-valid-jwt")
		if err == nil {
			t.Error("ValidateAccessToken() should reject malformed token")
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – ValidateRefreshToken
// ---------------------------------------------------------------------------

func TestValidateRefreshToken(t *testing.T) {
	svc, _ := newTestAuthService()
	userID := uuid.New().String()

	t.Run("accepts valid refresh token", func(t *testing.T) {
		pair, err := svc.GenerateTokenPair(userID)
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}
		claims, err := svc.ValidateRefreshToken(pair.RefreshToken)
		if err != nil {
			t.Fatalf("ValidateRefreshToken() error = %v", err)
		}
		if claims.UserID != userID {
			t.Errorf("UserID = %q, want %q", claims.UserID, userID)
		}
		if claims.TokenType != "refresh" {
			t.Errorf("TokenType = %q, want %q", claims.TokenType, "refresh")
		}
	})

	t.Run("rejects access token as refresh token", func(t *testing.T) {
		pair, err := svc.GenerateTokenPair(userID)
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}
		_, err = svc.ValidateRefreshToken(pair.AccessToken)
		if err == nil {
			t.Error("ValidateRefreshToken() should reject access token")
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – HandleOAuthCallback
// ---------------------------------------------------------------------------

func TestHandleOAuthCallback(t *testing.T) {
	t.Run("creates new user when not found", func(t *testing.T) {
		svc, repo := newTestAuthService()
		oauthUser := &OAuthUser{
			Provider:       "google",
			ProviderUserID: "google-123",
			Email:          "new@example.com",
			Nickname:       "NewUser",
			AvatarURL:      "https://example.com/avatar.png",
		}

		tokens, user, err := svc.HandleOAuthCallback("google", oauthUser)
		if err != nil {
			t.Fatalf("HandleOAuthCallback() error = %v", err)
		}
		if tokens == nil {
			t.Fatal("tokens is nil")
		}
		if user == nil {
			t.Fatal("user is nil")
		}
		if user.ProviderUserID != "google-123" {
			t.Errorf("ProviderUserID = %q, want %q", user.ProviderUserID, "google-123")
		}

		// Verify user was stored.
		stored, err := repo.FindByID(user.ID.String())
		if err != nil {
			t.Fatalf("user not found in repo: %v", err)
		}
		if stored.Email == nil || *stored.Email != "new@example.com" {
			t.Error("stored user email mismatch")
		}
	})

	t.Run("returns existing user when found", func(t *testing.T) {
		svc, repo := newTestAuthService()

		// Pre-create user.
		existingEmail := "existing@example.com"
		existingNickname := "ExistingUser"
		existingUser := &models.User{
			ID:             uuid.New(),
			Provider:       models.AuthProviderGoogle,
			ProviderUserID: "google-456",
			Email:          &existingEmail,
			Nickname:       &existingNickname,
		}
		if err := repo.Create(existingUser); err != nil {
			t.Fatalf("failed to create existing user: %v", err)
		}

		oauthUser := &OAuthUser{
			Provider:       "google",
			ProviderUserID: "google-456",
			Email:          "existing@example.com",
			Nickname:       "ExistingUser",
		}

		tokens, user, err := svc.HandleOAuthCallback("google", oauthUser)
		if err != nil {
			t.Fatalf("HandleOAuthCallback() error = %v", err)
		}
		if tokens == nil {
			t.Fatal("tokens is nil")
		}
		if user.ID != existingUser.ID {
			t.Errorf("user ID = %s, want %s", user.ID, existingUser.ID)
		}
	})

	t.Run("rejects unsupported provider", func(t *testing.T) {
		svc, _ := newTestAuthService()
		oauthUser := &OAuthUser{
			Provider:       "unsupported",
			ProviderUserID: "xxx",
		}
		_, _, err := svc.HandleOAuthCallback("unsupported", oauthUser)
		if err == nil {
			t.Error("HandleOAuthCallback() should reject unsupported provider")
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – RefreshToken
// ---------------------------------------------------------------------------

func TestRefreshToken(t *testing.T) {
	t.Run("generates new pair for valid refresh token", func(t *testing.T) {
		svc, repo := newTestAuthService()

		// Create a user so FindByID succeeds during refresh.
		user := &models.User{
			ID:             uuid.New(),
			Provider:       models.AuthProviderGoogle,
			ProviderUserID: "google-refresh-test",
		}
		if err := repo.Create(user); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		pair, err := svc.GenerateTokenPair(user.ID.String())
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}

		newPair, err := svc.RefreshToken(pair.RefreshToken)
		if err != nil {
			t.Fatalf("RefreshToken() error = %v", err)
		}
		if newPair.AccessToken == "" {
			t.Error("new AccessToken is empty")
		}
		if newPair.RefreshToken == "" {
			t.Error("new RefreshToken is empty")
		}
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		svc, _ := newTestAuthService()
		_, err := svc.RefreshToken("invalid-token-string")
		if err == nil {
			t.Error("RefreshToken() should reject invalid token")
		}
	})

	t.Run("rejects access token used as refresh token", func(t *testing.T) {
		svc, repo := newTestAuthService()

		user := &models.User{
			ID:             uuid.New(),
			Provider:       models.AuthProviderGoogle,
			ProviderUserID: "google-access-as-refresh",
		}
		if err := repo.Create(user); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		pair, err := svc.GenerateTokenPair(user.ID.String())
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}

		_, err = svc.RefreshToken(pair.AccessToken)
		if err == nil {
			t.Error("RefreshToken() should reject access token")
		}
	})

	t.Run("rejects refresh token for non-existent user", func(t *testing.T) {
		svc, _ := newTestAuthService()

		// Generate a token for a user that doesn't exist in the repo.
		fakeUserID := uuid.New().String()
		pair, err := svc.GenerateTokenPair(fakeUserID)
		if err != nil {
			t.Fatalf("GenerateTokenPair() error = %v", err)
		}

		_, err = svc.RefreshToken(pair.RefreshToken)
		if err == nil {
			t.Error("RefreshToken() should reject token for non-existent user")
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – Token signed with wrong secret
// ---------------------------------------------------------------------------

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	svc, _ := newTestAuthService()
	userID := uuid.New().String()

	// Sign with a different secret.
	now := time.Now()
	claims := &TokenClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			Issuer:    "subkeep",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = svc.ValidateAccessToken(tokenString)
	if err == nil {
		t.Error("ValidateAccessToken() should reject token signed with wrong secret")
	}
}

// ---------------------------------------------------------------------------
// Tests – Logout (smoke test)
// ---------------------------------------------------------------------------

func TestLogout(t *testing.T) {
	svc, _ := newTestAuthService()
	err := svc.Logout(uuid.New().String())
	if err != nil {
		t.Errorf("Logout() error = %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests – GetUserByID
// ---------------------------------------------------------------------------

func TestGetUserByID(t *testing.T) {
	svc, repo := newTestAuthService()

	t.Run("returns user when found", func(t *testing.T) {
		user := &models.User{
			ID:             uuid.New(),
			Provider:       models.AuthProviderKakao,
			ProviderUserID: "kakao-999",
		}
		if err := repo.Create(user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		found, err := svc.GetUserByID(user.ID.String())
		if err != nil {
			t.Fatalf("GetUserByID() error = %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("ID = %s, want %s", found.ID, user.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		_, err := svc.GetUserByID(uuid.New().String())
		if err == nil {
			t.Error("GetUserByID() should return error for non-existent user")
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – NewAuthService
// ---------------------------------------------------------------------------

func TestNewAuthService(t *testing.T) {
	repo := newMockUserRepo()
	cfg := config.JWTConfig{Secret: testJWTSecret}
	svc := NewAuthService(repo, cfg)
	if svc == nil {
		t.Fatal("NewAuthService() returned nil")
	}

	// Verify defaults are used when Expiration/RefreshExpiration are zero.
	pair, err := svc.GenerateTokenPair(uuid.New().String())
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}
	// Default access expiry should be 1 hour = 3600 seconds.
	if pair.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want 3600 (default 1h)", pair.ExpiresIn)
	}

	// Verify generated tokens are valid.
	_, err = svc.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Errorf("ValidateAccessToken() on default-config token failed: %v", err)
	}
	_ = fmt.Sprintf("suppress unused import") // ensure fmt is used
}
