package itssender

import (
	"context"
	"fmt"
	"net/http"
)

// MessagesService provides WhatsApp messaging operations (text, media, reactions, polls, locations).
type MessagesService struct {
	client *Client
}

// SendText dispatches a plain text message to a WhatsApp user or group chat.
func (s *MessagesService) SendText(ctx context.Context, deviceID string, req SendTextMessageRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/text", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendMedia dispatches an image, document, video, audio, or sticker attachment.
func (s *MessagesService) SendMedia(ctx context.Context, deviceID string, req SendMediaMessageRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/media", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Revoke deletes a previously sent message for everyone in the conversation.
func (s *MessagesService) Revoke(ctx context.Context, deviceID string, req RevokeMessageRequest) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/revoke", deviceID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendReaction adds or updates an emoji reaction on a target WhatsApp message.
func (s *MessagesService) SendReaction(ctx context.Context, deviceID string, req SendReactionRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/reaction", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendPoll dispatches an interactive multi-choice poll to a 1-on-1 or group chat.
func (s *MessagesService) SendPoll(ctx context.Context, deviceID string, req SendPollRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/poll", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendLocation sends GPS map coordinates with optional name and address labels.
func (s *MessagesService) SendLocation(ctx context.Context, deviceID string, req SendLocationMessageRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/location", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendContact dispatches a vCard contact card into the chat.
func (s *MessagesService) SendContact(ctx context.Context, deviceID string, req SendContactMessageRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/contact", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MarkRead marks messages in a conversation as read (sends blue double tick).
func (s *MessagesService) MarkRead(ctx context.Context, deviceID string, req MarkReadRequest) (*ActionSuccessResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/chat/read", deviceID)
	var resp ActionSuccessResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Edit modifies the content of a previously sent text message.
func (s *MessagesService) Edit(ctx context.Context, deviceID string, req EditMessageRequest) (*SendMessageResponse, error) {
	path := fmt.Sprintf("/api/v1/devices/%s/messages/edit", deviceID)
	var resp SendMessageResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
