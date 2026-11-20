package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// GroupsService provides WhatsApp group administration and membership management.
type GroupsService struct {
	client *Client
}

// List returns all WhatsApp groups the device is currently participating in.
func (s *GroupsService) List(ctx context.Context, deviceID string) ([]GroupItem, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups", deviceID)
	var groups []GroupItem
	if err := s.client.do(ctx, http.MethodGet, path, nil, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// Get retrieves metadata, subject, and participant list for a specific group.
func (s *GroupsService) Get(ctx context.Context, deviceID string, groupJID string) (*GroupItem, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/metadata", deviceID, groupJID)
	var group GroupItem
	if err := s.client.do(ctx, http.MethodGet, path, nil, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// Create provisions a new WhatsApp group chat with a subject and initial member list.
func (s *GroupsService) Create(ctx context.Context, deviceID string, req CreateGroupRequest) (*GroupCreateResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/create", deviceID)
	var resp GroupCreateResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update modifies a group's subject title or topic description.
func (s *GroupsService) Update(ctx context.Context, deviceID string, groupJID string, req UpdateGroupMetadataRequest) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/metadata", deviceID, groupJID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ManageParticipants modifies membership (add, remove, promote to admin, demote from admin).
func (s *GroupsService) ManageParticipants(ctx context.Context, deviceID string, groupJID string, req UpdateGroupParticipantsRequest) (*GroupActionResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/participants", deviceID, groupJID)
	var resp GroupActionResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInviteCode retrieves the active group invitation link and invite code.
func (s *GroupsService) GetInviteCode(ctx context.Context, deviceID string, groupJID string) (*GroupInviteLinkResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/invite-link", deviceID, groupJID)
	var resp GroupInviteLinkResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevokeInviteCode regenerates/revokes the active group invite link.
func (s *GroupsService) RevokeInviteCode(ctx context.Context, deviceID string, groupJID string) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/revoke-invite-link", deviceID, groupJID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Join accepts an invitation and adds the device to a group via invite code.
func (s *GroupsService) Join(ctx context.Context, deviceID string, req JoinGroupRequest) (*GroupActionResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/join", deviceID)
	var resp GroupActionResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Leave exits the device from the specified WhatsApp group.
func (s *GroupsService) Leave(ctx context.Context, deviceID string, groupJID string) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/groups/%s/leave", deviceID, groupJID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
