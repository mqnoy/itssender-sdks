package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// ContactsService manages contact WhatsApp verification, profile avatars, and bio status.
type ContactsService struct {
	client *Client
}

// Check queries WhatsApp to verify whether one or more phone numbers are registered WhatsApp accounts.
func (s *ContactsService) Check(ctx context.Context, deviceID string, req CheckContactsRequest) (*CheckContactsResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/contacts/check", deviceID)
	var resp CheckContactsResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetProfilePicture retrieves the high-resolution profile picture/avatar URL for a contact or group.
func (s *ContactsService) GetProfilePicture(ctx context.Context, deviceID string, jid string) (*ContactAvatarResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/contacts/%s/avatar", deviceID, jid)
	var resp ContactAvatarResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStatus retrieves the WhatsApp profile status/about bio of the connected device.
func (s *ContactsService) GetStatus(ctx context.Context, deviceID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/profile/status", deviceID)
	var resp map[string]interface{}
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateStatus updates the WhatsApp status/about text of the device account.
func (s *ContactsService) UpdateStatus(ctx context.Context, deviceID string, req UpdateStatusRequest) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/profile/status", deviceID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
