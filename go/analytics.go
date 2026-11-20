package itssender

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// AnalyticsService provides access to message metrics, quota counters, and dispatch audit logs.
type AnalyticsService struct {
	client *Client
}

// GetOverview returns aggregate metrics (messages sent, failed, storage, active devices, remaining quota).
func (s *AnalyticsService) GetOverview(ctx context.Context) (*UsageAnalytics, error) {
	var analytics UsageAnalytics
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/analytics/usage", nil, &analytics); err != nil {
		return nil, err
	}
	return &analytics, nil
}

// GetMessages retrieves recent message dispatch audit logs with optional pagination limit.
func (s *AnalyticsService) GetMessages(ctx context.Context, limit ...int64) ([]AnalyticsLog, error) {
	path := "/api/v1/analytics/logs"
	if len(limit) > 0 && limit[0] > 0 {
		path = buildURLWithQuery(path, map[string]string{"limit": strconv.FormatInt(limit[0], 10)})
	}
	var logs []AnalyticsLog
	if err := s.client.do(ctx, http.MethodGet, path, nil, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetDevices retrieves container runtime metrics (CPU, memory, goroutines) for a specific device.
func (s *AnalyticsService) GetDevices(ctx context.Context, id string) (*DeviceInstanceUsage, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/usage", id)
	var usage DeviceInstanceUsage
	if err := s.client.do(ctx, http.MethodGet, path, nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}
