package itssender

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// DevicesService manages WhatsApp client instances and Docker containers.
type DevicesService struct {
	client *Client
}

// List returns all WhatsApp devices provisioned for the authenticated account.
func (s *DevicesService) List(ctx context.Context) ([]Device, error) {
	var devices []Device
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/devices", nil, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

// Get retrieves details for a single WhatsApp device instance.
func (s *DevicesService) Get(ctx context.Context, id string) (*Device, error) {
	var device Device
	path := fmt.Sprintf("/api/v1/devices/%s", id)
	if err := s.client.do(ctx, http.MethodGet, path, nil, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

// Create provisions a new WhatsApp device container.
func (s *DevicesService) Create(ctx context.Context, req CreateDeviceRequest) (*Device, error) {
	var device Device
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/devices", req, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

// Delete terminates and removes a WhatsApp device container and session state.
func (s *DevicesService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/devices/%s", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// GetQR retrieves a QR code string or pairing code to link a WhatsApp session.
func (s *DevicesService) GetQR(ctx context.Context, id string) (*QRCodeResponse, error) {
	var qr QRCodeResponse
	path := fmt.Sprintf("/api/v1/devices/%s/qr", id)
	if err := s.client.do(ctx, http.MethodGet, path, nil, &qr); err != nil {
		return nil, err
	}
	return &qr, nil
}

// PairCode requests an 8-digit WhatsApp pairing code for linking via phone number.
func (s *DevicesService) PairCode(ctx context.Context, id string, req PairCodeRequest) (*PairCodeResponse, error) {
	var resp PairCodeResponse
	path := fmt.Sprintf("/api/v1/devices/%s/pair-code", id)
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Restart re-initializes or restarts a WhatsApp device instance container.
func (s *DevicesService) Restart(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/devices/%s/restart", id)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}

// GetLogs fetches container stdout/stderr log lines for debugging.
func (s *DevicesService) GetLogs(ctx context.Context, id string, lines ...int) (*DeviceLogsResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/logs", id)
	if len(lines) > 0 && lines[0] > 0 {
		path = buildURLWithQuery(path, map[string]string{"lines": strconv.Itoa(lines[0])})
	}
	var logs DeviceLogsResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &logs); err != nil {
		return nil, err
	}
	return &logs, nil
}

// UpdateSettings updates configuration parameters (e.g. proxy, CPU/memory limits) for a device.
func (s *DevicesService) UpdateSettings(ctx context.Context, id string, req UpdateDeviceRequest) (*Device, error) {
	path := fmt.Sprintf("/api/v1/devices/%s", id)
	var device Device
	if err := s.client.do(ctx, http.MethodPut, path, req, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

// SetPresence updates global online or offline presence for a device.
func (s *DevicesService) SetPresence(ctx context.Context, id string, req DevicePresenceRequest) error {
	path := fmt.Sprintf("/api/v1/devices/%s/presence/device", id)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}

// SetChatPresence sends real-time typing ("composing") or audio recording indicators.
func (s *DevicesService) SetChatPresence(ctx context.Context, id string, req ChatPresenceRequest) error {
	path := fmt.Sprintf("/api/v1/devices/%s/presence/chat", id)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}

// GetState retrieves runtime instance metrics and health status.
func (s *DevicesService) GetState(ctx context.Context, id string) (*DeviceInstanceUsage, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/usage", id)
	var usage DeviceInstanceUsage
	if err := s.client.do(ctx, http.MethodGet, path, nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}
