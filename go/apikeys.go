package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// APIKeysService handles management of developer API keys for programmatic access.
type APIKeysService struct {
	client *Client
}

// List returns all active API keys scoped to the authenticated user.
func (s *APIKeysService) List(ctx context.Context) ([]ApiKey, error) {
	var keys []ApiKey
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/user/keys", nil, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// Create generates a new API key. Note that the full key secret is returned only once.
func (s *APIKeysService) Create(ctx context.Context, req CreateApiKeyRequest) (*ApiKeyCreated, error) {
	var created ApiKeyCreated
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/user/keys", req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Delete permanently revokes an API key by its unique ID.
func (s *APIKeysService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/user/keys/%s", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
