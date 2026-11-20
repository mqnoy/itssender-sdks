package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// CampaignsService manages bulk broadcast campaigns and delivery schedules.
type CampaignsService struct {
	client *Client
}

// List returns all campaigns created by the authenticated account.
func (s *CampaignsService) List(ctx context.Context) ([]Campaign, error) {
	var campaigns []Campaign
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/campaigns", nil, &campaigns); err != nil {
		return nil, err
	}
	return campaigns, nil
}

// Create schedules or provisions a new marketing bulk blast campaign.
func (s *CampaignsService) Create(ctx context.Context, req CreateCampaignRequest) (*Campaign, error) {
	var campaign Campaign
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/campaigns", req, &campaign); err != nil {
		return nil, err
	}
	return &campaign, nil
}

// Get retrieves detailed status and progress counters for a specific campaign.
func (s *CampaignsService) Get(ctx context.Context, id string) (*CampaignDetail, error) {
	path := fmt.Sprintf("/api/v1/campaigns/%s", id)
	var detail CampaignDetail
	if err := s.client.do(ctx, http.MethodGet, path, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// Delete archives/deletes a campaign by its unique ID.
func (s *CampaignsService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/campaigns/%s", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// Start triggers an immediate start of an scheduled or queued campaign.
func (s *CampaignsService) Start(ctx context.Context, id string) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/campaigns/%s/start", id)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses/archives execution of an in-progress campaign.
func (s *CampaignsService) Pause(ctx context.Context, id string) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/campaigns/%s/archive", id)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAnalytics retrieves message delivery progress and status distribution for a campaign.
func (s *CampaignsService) GetAnalytics(ctx context.Context, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/campaigns/%s/analytics", id)
	var resp map[string]interface{}
	err := s.client.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		// Fallback to campaign detail if custom analytics endpoint is absent
		detail, getErr := s.Get(ctx, id)
		if getErr != nil {
			return nil, err
		}
		return map[string]interface{}{
			"id":           detail.ID,
			"name":         detail.Name,
			"status":       detail.Status,
			"total_count":  detail.TotalCount,
			"sent_count":   detail.SentCount,
			"failed_count": detail.FailedCount,
		}, nil
	}
	return resp, nil
}
