package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/subkeep/backend/config"
	"github.com/subkeep/backend/utils"
)

// OAuth provider endpoints.
const (
	googleAuthURL    = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL   = "https://oauth2.googleapis.com/token"
	googleUserURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
	kakaoAuthURL     = "https://kauth.kakao.com/oauth/authorize"
	kakaoTokenURL    = "https://kauth.kakao.com/oauth/token"
	kakaoUserURL     = "https://kapi.kakao.com/v2/user/me"
	naverAuthURL     = "https://nid.naver.com/oauth2.0/authorize"
	naverTokenURL    = "https://nid.naver.com/oauth2.0/token"
	naverUserURL     = "https://openapi.naver.com/v1/nid/me"
)

// OAuthService handles OAuth provider integrations.
type OAuthService struct {
	oauthConfig config.OAuthConfig
	httpClient  *http.Client
}

// NewOAuthService creates a new OAuthService.
func NewOAuthService(oauthConfig config.OAuthConfig) *OAuthService {
	return &OAuthService{
		oauthConfig: oauthConfig,
		httpClient:  &http.Client{},
	}
}

// GetAuthorizationURL returns the OAuth authorization URL for the given provider.
func (s *OAuthService) GetAuthorizationURL(provider string, redirectURI string, state string) (string, error) {
	if !isValidProvider(provider) {
		return "", utils.ErrBadRequest("지원하지 않는 인증 제공자입니다: " + provider)
	}

	var authURL, clientID, scope string

	switch provider {
	case "google":
		authURL = googleAuthURL
		clientID = s.oauthConfig.Google.ClientID
		scope = "openid email profile"
	case "kakao":
		authURL = kakaoAuthURL
		clientID = s.oauthConfig.Kakao.ClientID
		scope = "profile_nickname profile_image account_email"
	case "naver":
		authURL = naverAuthURL
		clientID = s.oauthConfig.Naver.ClientID
		scope = ""
	default:
		return "", utils.ErrBadRequest("지원하지 않는 인증 제공자입니다: " + provider)
	}

	if clientID == "" {
		return "", utils.ErrBadRequest(provider + " OAuth가 구성되지 않았습니다")
	}

	params := url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"response_type": {"code"},
		"state":         {state},
	}
	if scope != "" {
		params.Set("scope", scope)
	}

	return authURL + "?" + params.Encode(), nil
}

// ExchangeCode exchanges an authorization code for user info from the given provider.
func (s *OAuthService) ExchangeCode(provider string, code string, redirectURI string) (*OAuthUser, error) {
	if !isValidProvider(provider) {
		return nil, utils.ErrBadRequest("지원하지 않는 인증 제공자입니다: " + provider)
	}

	switch provider {
	case "google":
		return s.exchangeGoogle(code, redirectURI)
	case "kakao":
		return s.exchangeKakao(code, redirectURI)
	case "naver":
		return s.exchangeNaver(code, redirectURI)
	default:
		return nil, utils.ErrBadRequest("지원하지 않는 인증 제공자입니다: " + provider)
	}
}

// exchangeGoogle handles Google OAuth code exchange and user info retrieval.
func (s *OAuthService) exchangeGoogle(code string, redirectURI string) (*OAuthUser, error) {
	cfg := s.oauthConfig.Google

	// Exchange code for access token.
	tokenData, err := s.exchangeCodeForToken(googleTokenURL, url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, fmt.Errorf("google token exchange: %w", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok || accessToken == "" {
		return nil, fmt.Errorf("google: missing access_token in response")
	}

	// Fetch user profile.
	userInfo, err := s.fetchUserInfo(googleUserURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("google user info: %w", err)
	}

	return &OAuthUser{
		Provider:       "google",
		ProviderUserID: getString(userInfo, "id"),
		Email:          getString(userInfo, "email"),
		Nickname:       getString(userInfo, "name"),
		AvatarURL:      getString(userInfo, "picture"),
	}, nil
}

// exchangeKakao handles Kakao OAuth code exchange and user info retrieval.
func (s *OAuthService) exchangeKakao(code string, redirectURI string) (*OAuthUser, error) {
	cfg := s.oauthConfig.Kakao

	// Exchange code for access token.
	tokenData, err := s.exchangeCodeForToken(kakaoTokenURL, url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, fmt.Errorf("kakao token exchange: %w", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok || accessToken == "" {
		return nil, fmt.Errorf("kakao: missing access_token in response")
	}

	// Fetch user profile.
	userInfo, err := s.fetchUserInfo(kakaoUserURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("kakao user info: %w", err)
	}

	oauthUser := &OAuthUser{
		Provider:       "kakao",
		ProviderUserID: fmt.Sprintf("%v", userInfo["id"]),
	}

	// Kakao nests profile data under kakao_account and properties.
	if account, ok := userInfo["kakao_account"].(map[string]interface{}); ok {
		oauthUser.Email = getString(account, "email")
		if profile, ok := account["profile"].(map[string]interface{}); ok {
			oauthUser.Nickname = getString(profile, "nickname")
			oauthUser.AvatarURL = getString(profile, "profile_image_url")
		}
	}

	return oauthUser, nil
}

// exchangeNaver handles Naver OAuth code exchange and user info retrieval.
func (s *OAuthService) exchangeNaver(code string, redirectURI string) (*OAuthUser, error) {
	cfg := s.oauthConfig.Naver

	// Exchange code for access token.
	tokenData, err := s.exchangeCodeForToken(naverTokenURL, url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, fmt.Errorf("naver token exchange: %w", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok || accessToken == "" {
		return nil, fmt.Errorf("naver: missing access_token in response")
	}

	// Fetch user profile.
	userInfo, err := s.fetchUserInfo(naverUserURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("naver user info: %w", err)
	}

	oauthUser := &OAuthUser{
		Provider: "naver",
	}

	// Naver nests user data under "response".
	if response, ok := userInfo["response"].(map[string]interface{}); ok {
		oauthUser.ProviderUserID = getString(response, "id")
		oauthUser.Email = getString(response, "email")
		oauthUser.Nickname = getString(response, "nickname")
		oauthUser.AvatarURL = getString(response, "profile_image")
	}

	return oauthUser, nil
}

// exchangeCodeForToken performs an HTTP POST to the token endpoint and returns
// the parsed JSON response.
func (s *OAuthService) exchangeCodeForToken(tokenURL string, params url.Values) (map[string]interface{}, error) {
	resp, err := s.httpClient.PostForm(tokenURL, params)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("token exchange failed",
			"statusCode", resp.StatusCode,
			"body", string(body),
			"url", tokenURL,
		)
		return nil, fmt.Errorf("token exchange returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse token response: %w", err)
	}

	return result, nil
}

// fetchUserInfo makes a GET request to the user info endpoint with the given
// access token in the Authorization header.
func (s *OAuthService) fetchUserInfo(userInfoURL string, accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create user info request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read user info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("user info request failed",
			"statusCode", resp.StatusCode,
			"body", string(body),
			"url", userInfoURL,
		)
		return nil, fmt.Errorf("user info returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse user info response: %w", err)
	}

	return result, nil
}

// getString safely retrieves a string value from a map.
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// Compile-time check to ensure the file compiles cleanly.
