package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// OTPService handles high-priority transactional one-time password dispatch and verification.
type OTPService struct {
	client *Client
}

// Send dispatches a transactional OTP code over WhatsApp with configurable TTL and attempt limits.
func (s *OTPService) Send(ctx context.Context, req SendOTPRequest) (*SendOTPResponse, error) {
	var resp SendOTPResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/otp/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify validates an OTP verification code against an active OTP session.
func (s *OTPService) Verify(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	var resp VerifyOTPResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/otp/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resend dispatches a new OTP challenge for a phone number or session.
func (s *OTPService) Resend(ctx context.Context, req SendOTPRequest) (*SendOTPResponse, error) {
	return s.Send(ctx, req)
}

// Get retrieves full audit details for a specific OTP verification attempt.
func (s *OTPService) Get(ctx context.Context, id string) (*OTPVerificationDetail, error) {
	path := fmt.Sprintf("/api/v1/otp/%s", id)
	var detail OTPVerificationDetail
	if err := s.client.do(ctx, http.MethodGet, path, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// List returns a list of recent OTP verification sessions.
func (s *OTPService) List(ctx context.Context) ([]OTPVerificationDetail, error) {
	var list []OTPVerificationDetail
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/otp", nil, &list); err != nil {
		return nil, err
	}
	return list, nil
}
