package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// WebhooksService manages tenant webhook endpoints conforming to Standard Webhooks v1.0.0.
type WebhooksService struct {
	client *Client
}

// List returns all webhook configurations registered for the tenant.
func (s *WebhooksService) List(ctx context.Context) ([]WebhookConfigResponse, error) {
	var webhooks []WebhookConfigResponse
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/webhooks", nil, &webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

// Create registers a new webhook endpoint with subscribed event topics.
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest) (*WebhookConfigResponse, error) {
	var config WebhookConfigResponse
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/webhooks", req, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Get retrieves configuration details for a specific webhook endpoint.
func (s *WebhooksService) Get(ctx context.Context, id string) (*WebhookConfigResponse, error) {
	path := fmt.Sprintf("/api/v1/webhooks/%s", id)
	var config WebhookConfigResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Update modifies subscribed event topics, URL destination, or active state.
func (s *WebhooksService) Update(ctx context.Context, id string, req UpdateWebhookRequest) (*WebhookConfigResponse, error) {
	path := fmt.Sprintf("/api/v1/webhooks/%s", id)
	var config WebhookConfigResponse
	if err := s.client.do(ctx, http.MethodPut, path, req, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Delete archives and permanently removes a webhook endpoint.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/webhooks/%s", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// Test dispatches a sample ping payload to verify webhook receiver connectivity.
func (s *WebhooksService) Test(ctx context.Context, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/webhooks/%s/test", id)
	var resp map[string]interface{}
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListDeliveries retrieves audit logs of webhook dispatch attempts.
func (s *WebhooksService) ListDeliveries(ctx context.Context) (*WebhookDeliveryListResponse, error) {
	var resp WebhookDeliveryListResponse
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/webhooks/deliveries", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RetryDelivery retries an outbound webhook dispatch for a given delivery ID.
func (s *WebhooksService) RetryDelivery(ctx context.Context, deliveryID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/webhooks/deliveries/%s/retry", deliveryID)
	var resp map[string]interface{}
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
