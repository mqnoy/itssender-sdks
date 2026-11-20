package itssender

import (
	"context"
	"net/http"
)

// AuthService handles user registration, authentication sessions, and profile retrieval.
type AuthService struct {
	client *Client
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	var resp AuthResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/auth/register", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetBearerToken(resp.Token)
	}
	return &resp, nil
}

// Login authenticates credentials and returns a JWT session token.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	var resp AuthResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/auth/login", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetBearerToken(resp.Token)
	}
	return &resp, nil
}

// Me retrieves the authenticated user's profile.
func (s *AuthService) Me(ctx context.Context) (*UserResponse, error) {
	var resp UserResponse
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/auth/me", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RefreshToken refreshes an active session token.
func (s *AuthService) RefreshToken(ctx context.Context) (*AuthResponse, error) {
	var resp AuthResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/auth/refresh", nil, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetBearerToken(resp.Token)
	}
	return &resp, nil
}

// GetConfig retrieves public authentication flags (e.g. whether self-registration is enabled).
func (s *AuthService) GetConfig(ctx context.Context) (*PublicConfigResponse, error) {
	var resp PublicConfigResponse
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/auth/config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
