package itssender

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// TemplatesService manages multi-tenant rich-text WhatsApp message templates.
type TemplatesService struct {
	client *Client
}

// List returns all templates belonging to the authenticated account.
func (s *TemplatesService) List(ctx context.Context) (*TemplateListResponse, error) {
	var resp TemplateListResponse
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/templates", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create stores a new rich-text message template.
func (s *TemplatesService) Create(ctx context.Context, req TemplateCreateRequest) (*TemplateItem, error) {
	var item TemplateItem
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/templates", req, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Get retrieves a specific template by its unique ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*TemplateItem, error) {
	path := fmt.Sprintf("/api/v1/templates/%s", id)
	var item TemplateItem
	if err := s.client.do(ctx, http.MethodGet, path, nil, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Update modifies content, variables, or metadata of an existing template.
func (s *TemplatesService) Update(ctx context.Context, id string, req TemplateUpdateRequest) (*TemplateItem, error) {
	path := fmt.Sprintf("/api/v1/templates/%s", id)
	var item TemplateItem
	if err := s.client.do(ctx, http.MethodPut, path, req, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Delete permanently archives/deletes a template by its ID.
func (s *TemplatesService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/templates/%s", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// Render retrieves the template and substitutes variables like {{key}} or {{.key}} with values.
func (s *TemplatesService) Render(ctx context.Context, id string, variables map[string]string) (string, error) {
	item, err := s.Get(ctx, id)
	if err != nil {
		return "", err
	}

	content := item.ContentRaw
	if content == "" {
		content = item.ContentHTML
	}

	for k, v := range variables {
		content = strings.ReplaceAll(content, "{{."+k+"}}", v)
		content = strings.ReplaceAll(content, "{{"+k+"}}", v)
		content = strings.ReplaceAll(content, "{{ "+k+" }}", v)
	}

	return content, nil
}
