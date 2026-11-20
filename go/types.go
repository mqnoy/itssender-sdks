package itssender

import (
	"fmt"
)

// Invariant 1 & 2 Enforced:
// All monetary quantities, counters, quotas, and Unix timestamps use int64.
// Millisecond Unix timestamps are used for all temporal fields in transit.

// APIError represents a structured error returned by the itssender API.
type APIError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	RawBody    string `json:"raw_body,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("itssender: API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("itssender: HTTP error %d: %s", e.StatusCode, e.Message)
}

// SuccessResponse represents a generic success response.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ActionSuccessResponse represents an action outcome response with timestamp.
type ActionSuccessResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"` // Unix timestamp in milliseconds
}

// ----------------- AUTH TYPES -----------------

// RegisterRequest payload for self-registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest payload for authenticating user sessions.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents user account details.
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsSuspended bool   `json:"is_suspended"`
	CreatedAt   int64  `json:"created_at"` // Unix timestamp in milliseconds
	UpdatedAt   int64  `json:"updated_at"` // Unix timestamp in milliseconds
}

// AuthResponse represents authentication response containing JWT token.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// PublicConfigResponse represents server public authentication flags.
type PublicConfigResponse struct {
	RegistrationEnabled bool   `json:"registration_enabled"`
	Version             string `json:"version,omitempty"`
}

// ----------------- API KEY TYPES -----------------

// CreateApiKeyRequest payload for generating an API key.
type CreateApiKeyRequest struct {
	Name string `json:"name"`
}

// ApiKey represents an API key metadata record.
type ApiKey struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	KeyPrefix  string `json:"key_prefix"`
	CreatedAt  int64  `json:"created_at"`             // Unix timestamp in milliseconds
	LastUsedAt *int64 `json:"last_used_at,omitempty"` // Unix timestamp in milliseconds
}

// ApiKeyCreated represents a freshly provisioned API key with full plaintext secret.
type ApiKeyCreated struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	KeyPrefix string `json:"key_prefix"`
	Key       string `json:"key"`        // Plaintext secret shown only once
	CreatedAt int64  `json:"created_at"` // Unix timestamp in milliseconds
}

// ----------------- DEVICE TYPES -----------------

// CreateDeviceRequest payload for provisioning a new WhatsApp instance.
type CreateDeviceRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number,omitempty"`
	ProxyURL    string `json:"proxy_url,omitempty"`
	NanoCPUs    int64  `json:"nano_cpus,omitempty"`
	MemoryMB    int64  `json:"memory_mb,omitempty"`
}

// UpdateDeviceRequest payload for modifying device settings.
type UpdateDeviceRequest struct {
	Name     *string `json:"name,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	NanoCPUs *int64  `json:"nano_cpus,omitempty"`
	MemoryMB *int64  `json:"memory_mb,omitempty"`
}

// Device represents a WhatsApp instance container record.
type Device struct {
	ID            string               `json:"id"`
	UserID        string               `json:"user_id"`
	Name          string               `json:"name"`
	PhoneNumber   string               `json:"phone_number,omitempty"`
	JID           string               `json:"jid,omitempty"`
	Status        string               `json:"status"`
	PushName      string               `json:"push_name,omitempty"`
	Platform      string               `json:"platform,omitempty"`
	ContainerID   string               `json:"container_id,omitempty"`
	Port          int64                `json:"port,omitempty"`
	NanoCPUs      int64                `json:"nano_cpus,omitempty"`
	MemoryMB      int64                `json:"memory_mb,omitempty"`
	CreatedAt     int64                `json:"created_at"` // Unix timestamp in milliseconds
	UpdatedAt     int64                `json:"updated_at"` // Unix timestamp in milliseconds
	InstanceUsage *DeviceInstanceUsage `json:"instance_usage,omitempty"`
}

// DeviceInstanceUsage represents runtime container metrics adhering to integer arithmetic.
type DeviceInstanceUsage struct {
	CPUUsagePct   int64  `json:"cpu_usage_pct,omitempty"`
	MemoryUsageMB int64  `json:"memory_usage_mb,omitempty"`
	MemoryLimitMB int64  `json:"memory_limit_mb,omitempty"`
	Goroutines    int64  `json:"goroutines,omitempty"`
	Status        string `json:"status,omitempty"`
	UptimeSeconds int64  `json:"uptime_seconds,omitempty"`
}

// DeviceLogsResponse represents log lines retrieved from a device container.
type DeviceLogsResponse struct {
	Logs []string `json:"logs"`
}

// QRCodeResponse represents a QR code for linking a WhatsApp device.
type QRCodeResponse struct {
	Code           string `json:"code,omitempty"`
	Raw            string `json:"raw,omitempty"`
	TimeoutSeconds int64  `json:"timeout_seconds,omitempty"`
}

// PairCodeRequest payload for requesting an 8-character phone pairing code.
type PairCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// PairCodeResponse represents the returned 8-digit WhatsApp pairing code.
type PairCodeResponse struct {
	PairCode       string `json:"pair_code"`
	TimeoutSeconds int64  `json:"timeout_seconds,omitempty"`
}

// DevicePresenceRequest payload for setting global device availability.
type DevicePresenceRequest struct {
	State string `json:"state"` // "available" | "unavailable"
}

// ChatPresenceRequest payload for sending chat state indicators (typing/recording).
type ChatPresenceRequest struct {
	To        string `json:"to,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	State     string `json:"state"` // "composing" | "recording" | "paused"
}

// ----------------- MESSAGING TYPES -----------------

// SendTextMessageRequest payload for dispatching plain text messages.
type SendTextMessageRequest struct {
	To                string   `json:"to,omitempty"`
	Recipient         string   `json:"recipient,omitempty"`
	Message           string   `json:"message,omitempty"`
	Text              string   `json:"text,omitempty"`
	ReplyToID         string   `json:"reply_to_id,omitempty"`
	QuotedMessageID   string   `json:"quoted_message_id,omitempty"`
	QuotedParticipant string   `json:"quoted_participant,omitempty"`
	Mentions          []string `json:"mentions,omitempty"`
}

// SendMediaMessageRequest payload for dispatching media attachments.
type SendMediaMessageRequest struct {
	To          string `json:"to,omitempty"`
	Recipient   string `json:"recipient,omitempty"`
	MediaType   string `json:"media_type"` // "image" | "document" | "video" | "audio" | "sticker"
	MediaURL    string `json:"media_url,omitempty"`
	URL         string `json:"url,omitempty"`
	MediaBase64 string `json:"media_base64,omitempty"`
	DataBase64  string `json:"data_base64,omitempty"`
	Caption     string `json:"caption,omitempty"`
	Filename    string `json:"filename,omitempty"`
	MimeType    string `json:"mime_type,omitempty"`
	ReplyToID   string `json:"reply_to_id,omitempty"`
}

// SendReactionRequest payload for sending an emoji reaction to a message.
type SendReactionRequest struct {
	To              string `json:"to,omitempty"`
	Recipient       string `json:"recipient,omitempty"`
	TargetMessageID string `json:"target_message_id,omitempty"`
	MessageID       string `json:"message_id,omitempty"`
	Reaction        string `json:"reaction,omitempty"`
	Emoji           string `json:"emoji,omitempty"`
}

// SendPollRequest payload for sending an interactive poll.
type SendPollRequest struct {
	To            string   `json:"to,omitempty"`
	Recipient     string   `json:"recipient,omitempty"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	MaxAnswers    int64    `json:"max_answers,omitempty"`
	MaxSelectable int64    `json:"max_selectable,omitempty"`
}

// SendLocationMessageRequest payload for sending GPS coordinates.
type SendLocationMessageRequest struct {
	To        string  `json:"to,omitempty"`
	Recipient string  `json:"recipient,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// SendContactMessageRequest payload for sending a vCard contact.
type SendContactMessageRequest struct {
	To          string `json:"to,omitempty"`
	Recipient   string `json:"recipient,omitempty"`
	DisplayName string `json:"display_name"`
	VCard       string `json:"vcard"`
}

// EditMessageRequest payload for editing previously sent text.
type EditMessageRequest struct {
	To              string `json:"to,omitempty"`
	Recipient       string `json:"recipient,omitempty"`
	TargetMessageID string `json:"target_message_id"`
	NewText         string `json:"new_text"`
}

// RevokeMessageRequest payload for deleting a message for everyone.
type RevokeMessageRequest struct {
	To              string `json:"to,omitempty"`
	Recipient       string `json:"recipient,omitempty"`
	TargetMessageID string `json:"target_message_id"`
}

// MarkReadRequest payload for marking messages as read in a chat.
type MarkReadRequest struct {
	ChatJID    string   `json:"chat_jid"`
	MessageIDs []string `json:"message_ids,omitempty"`
}

// SendMessageResponse represents the outcome of a dispatched message.
type SendMessageResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp in milliseconds
}

// ----------------- GROUP TYPES -----------------

// GroupParticipant represents a member in a WhatsApp group.
type GroupParticipant struct {
	JID          string `json:"jid"`
	IsAdmin      bool   `json:"is_admin"`
	IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
}

// GroupItem represents a WhatsApp group metadata record.
type GroupItem struct {
	JID          string             `json:"jid"`
	Subject      string             `json:"subject"`
	Owner        string             `json:"owner,omitempty"`
	Topic        string             `json:"topic,omitempty"`
	Participants []GroupParticipant `json:"participants,omitempty"`
}

// CreateGroupRequest payload for creating a new WhatsApp group.
type CreateGroupRequest struct {
	Subject      string   `json:"subject"`
	Participants []string `json:"participants"`
}

// GroupCreateResponse represents the result of group creation.
type GroupCreateResponse struct {
	Success   bool   `json:"success"`
	GroupJID  string `json:"group_jid"`
	OwnerJID  string `json:"owner_jid,omitempty"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp in milliseconds
}

// UpdateGroupMetadataRequest payload for updating group subject and topic.
type UpdateGroupMetadataRequest struct {
	Subject string `json:"subject,omitempty"`
	Topic   string `json:"topic,omitempty"`
}

// UpdateGroupParticipantsRequest payload for modifying group membership.
type UpdateGroupParticipantsRequest struct {
	Action       string   `json:"action"` // "add" | "remove" | "promote" | "demote"
	Participants []string `json:"participants"`
}

// GroupActionResponse represents the result of a group membership action.
type GroupActionResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"` // Unix timestamp in milliseconds
}

// GroupInviteLinkResponse represents group invitation code and URL.
type GroupInviteLinkResponse struct {
	Success    bool   `json:"success"`
	InviteCode string `json:"invite_code,omitempty"`
	InviteURL  string `json:"invite_url,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"` // Unix timestamp in milliseconds
}

// JoinGroupRequest payload for joining a group via invite link code.
type JoinGroupRequest struct {
	Code string `json:"code"`
}

// ----------------- CONTACT TYPES -----------------

// CheckContactsRequest payload for verifying WhatsApp registration.
type CheckContactsRequest struct {
	Phones []string `json:"phones"`
}

// ContactCheckResult represents registration state of a phone number.
type ContactCheckResult struct {
	Phone        string `json:"phone"`
	IsRegistered bool   `json:"is_registered"`
	JID          string `json:"jid,omitempty"`
}

// CheckContactsResponse represents contact verification outcomes.
type CheckContactsResponse struct {
	Success   bool                 `json:"success"`
	Results   []ContactCheckResult `json:"results"`
	Timestamp int64                `json:"timestamp"` // Unix timestamp in milliseconds
}

// ContactAvatarResponse represents avatar picture URL of a contact or group.
type ContactAvatarResponse struct {
	Success   bool   `json:"success"`
	URL       string `json:"url,omitempty"`
	ID        string `json:"id,omitempty"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp in milliseconds
}

// UpdateStatusRequest payload for updating WhatsApp profile status/about bio.
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// ----------------- CAMPAIGN TYPES -----------------

// CreateCampaignRequest payload for scheduling a bulk broadcast.
type CreateCampaignRequest struct {
	Name            string   `json:"name"`
	ChannelType     string   `json:"channel_type"` // "WHATSMEOW" | "META_CLOUD"
	Channel         string   `json:"channel,omitempty"`
	DeviceID        *string  `json:"device_id,omitempty"`
	MetaToken       *string  `json:"meta_token,omitempty"`
	MetaAPIKey      *string  `json:"meta_api_key,omitempty"`
	MetaPhoneID     *string  `json:"meta_phone_id,omitempty"`
	TemplateName    *string  `json:"template_name,omitempty"`
	MessageTemplate string   `json:"message_template"`
	Message         string   `json:"message,omitempty"`
	Recipients      []string `json:"recipients"`
	DelayMinSeconds int64    `json:"delay_min_seconds"`
	DelayMaxSeconds int64    `json:"delay_max_seconds"`
	MessageType     string   `json:"message_type,omitempty"` // "TEXT" | "MEDIA"
	MediaType       string   `json:"media_type,omitempty"`   // "image" | "video" | "audio" | "document" | "sticker"
	MediaURL        string   `json:"media_url,omitempty"`
	MediaDataBase64 string   `json:"media_data_base64,omitempty"`
	MediaCaption    string   `json:"media_caption,omitempty"`
	MediaFilename   string   `json:"media_filename,omitempty"`
	MediaMimeType   string   `json:"media_mime_type,omitempty"`
	ScheduledAt     int64    `json:"scheduled_at,omitempty"` // Unix timestamp in milliseconds
}

// Campaign represents a bulk messaging campaign record.
type Campaign struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Name            string  `json:"name"`
	ChannelType     string  `json:"channel_type"`
	DeviceID        *string `json:"device_id,omitempty"`
	MessageTemplate string  `json:"message_template"`
	MessageType     string  `json:"message_type,omitempty"`
	Status          string  `json:"status"`
	TotalCount      int64   `json:"total_count"`
	SentCount       int64   `json:"sent_count"`
	FailedCount     int64   `json:"failed_count"`
	ScheduledAt     int64   `json:"scheduled_at,omitempty"` // Unix timestamp in milliseconds
	CreatedAt       int64   `json:"created_at"`             // Unix timestamp in milliseconds
	UpdatedAt       int64   `json:"updated_at"`             // Unix timestamp in milliseconds
}

// CampaignDetail provides detailed campaign status including recipient stats.
type CampaignDetail struct {
	Campaign
	Recipients []string `json:"recipients,omitempty"`
}

// ----------------- ANALYTICS TYPES -----------------

// UsageAnalytics represents tenant-scoped metrics and counters.
type UsageAnalytics struct {
	MessagesSent     int64 `json:"messages_sent"`
	MessagesFailed   int64 `json:"messages_failed"`
	ActiveDevices    int64 `json:"active_devices"`
	TotalDevices     int64 `json:"total_devices"`
	StorageUsedBytes int64 `json:"storage_used_bytes"`
	QuotaLimit       int64 `json:"quota_limit"`
	QuotaRemaining   int64 `json:"quota_remaining"`
}

// AnalyticsLog represents an audit log entry for dispatched messages.
type AnalyticsLog struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	DeviceID     *string `json:"device_id,omitempty"`
	Recipient    string  `json:"recipient"`
	MessageType  string  `json:"message_type"`
	Status       string  `json:"status"`
	Direction    string  `json:"direction"`
	Timestamp    int64   `json:"timestamp"` // Unix timestamp in milliseconds
	ErrorMessage *string `json:"error_message,omitempty"`
}

// ----------------- OTP TYPES -----------------

// SendOTPRequest payload for dispatching high-priority transactional OTP.
type SendOTPRequest struct {
	PhoneNumber string                 `json:"phone_number"`
	TemplateID  *string                `json:"template_id,omitempty"`
	CustomCode  *string                `json:"custom_code,omitempty"`
	CodeLength  int64                  `json:"code_length,omitempty"`
	TTLMS       int64                  `json:"ttl_ms,omitempty"` // Lifetime in milliseconds
	MaxAttempts int64                  `json:"max_attempts,omitempty"`
	DeviceID    *string                `json:"device_id,omitempty"`
	Channel     string                 `json:"channel,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SendOTPResponse represents the outcome of OTP dispatch.
type SendOTPResponse struct {
	OTPID       string `json:"otp_id"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	ExpiresAt   int64  `json:"expires_at"` // Unix timestamp in milliseconds
	MaxAttempts int64  `json:"max_attempts"`
	CreatedAt   int64  `json:"created_at"` // Unix timestamp in milliseconds
}

// VerifyOTPRequest payload for validating an OTP code.
type VerifyOTPRequest struct {
	OTPID       string  `json:"otp_id"`
	Code        string  `json:"code"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// VerifyOTPResponse represents the outcome of an OTP verification attempt.
type VerifyOTPResponse struct {
	Success           bool   `json:"success"`
	OTPID             string `json:"otp_id"`
	Status            string `json:"status"`
	VerifiedAt        *int64 `json:"verified_at,omitempty"` // Unix timestamp in milliseconds
	RemainingAttempts *int64 `json:"remaining_attempts,omitempty"`
}

// OTPVerificationDetail represents full OTP record details for auditing.
type OTPVerificationDetail struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	TemplateID   *string                `json:"template_id,omitempty"`
	DeviceID     *string                `json:"device_id,omitempty"`
	PhoneNumber  string                 `json:"phone_number"`
	CodeLength   int64                  `json:"code_length"`
	Status       string                 `json:"status"`
	Attempts     int64                  `json:"attempts"`
	MaxAttempts  int64                  `json:"max_attempts"`
	Channel      string                 `json:"channel"`
	MessageLogID *string                `json:"message_log_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ExpiresAt    int64                  `json:"expires_at"`            // Unix timestamp in milliseconds
	VerifiedAt   *int64                 `json:"verified_at,omitempty"` // Unix timestamp in milliseconds
	CreatedAt    int64                  `json:"created_at"`            // Unix timestamp in milliseconds
}

// ----------------- TEMPLATE TYPES -----------------

// TemplateCreateRequest payload for authoring a rich-text template.
type TemplateCreateRequest struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description *string  `json:"description,omitempty"`
	ContentRaw  string   `json:"content_raw"`
	ContentHTML string   `json:"content_html"`
	Variables   []string `json:"variables,omitempty"`
}

// TemplateUpdateRequest payload for modifying a template.
type TemplateUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Description *string  `json:"description,omitempty"`
	ContentRaw  *string  `json:"content_raw,omitempty"`
	ContentHTML *string  `json:"content_html,omitempty"`
	Variables   []string `json:"variables,omitempty"`
}

// TemplateItem represents a stored rich-text message template.
type TemplateItem struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description *string  `json:"description,omitempty"`
	ContentRaw  string   `json:"content_raw"`
	ContentHTML string   `json:"content_html"`
	Variables   []string `json:"variables"`
	IsArchived  bool     `json:"is_archived"`
	ArchivedAt  *int64   `json:"archived_at,omitempty"` // Unix timestamp in milliseconds
	CreatedAt   int64    `json:"created_at"`            // Unix timestamp in milliseconds
	UpdatedAt   int64    `json:"updated_at"`            // Unix timestamp in milliseconds
}

// TemplateListResponse represents paginated list of templates.
type TemplateListResponse struct {
	Templates []TemplateItem `json:"templates"`
	Total     int64          `json:"total"`
}

// ----------------- WEBHOOK TYPES -----------------

// CreateWebhookRequest payload for registering a webhook endpoint.
type CreateWebhookRequest struct {
	URL      string   `json:"url"`
	Secret   *string  `json:"secret,omitempty"`
	Events   []string `json:"events"`
	IsActive *bool    `json:"is_active,omitempty"`
	MaskPII  *bool    `json:"mask_pii,omitempty"`
}

// UpdateWebhookRequest payload for updating an existing webhook configuration.
type UpdateWebhookRequest struct {
	URL      string   `json:"url,omitempty"`
	Secret   *string  `json:"secret,omitempty"`
	Events   []string `json:"events,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
	MaskPII  *bool    `json:"mask_pii,omitempty"`
}

// WebhookConfigResponse represents a tenant webhook configuration.
type WebhookConfigResponse struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	URL        string   `json:"url"`
	Secret     string   `json:"secret"`
	Events     []string `json:"events"`
	IsActive   bool     `json:"is_active"`
	IsArchived bool     `json:"is_archived"`
	MaskPII    bool     `json:"mask_pii"`
	CreatedAt  int64    `json:"created_at"` // Unix timestamp in milliseconds
	UpdatedAt  int64    `json:"updated_at"` // Unix timestamp in milliseconds
}

// WebhookDeliveryLog represents an audit record of an outbound webhook dispatch.
type WebhookDeliveryLog struct {
	ID            string                 `json:"id"`
	WebhookID     string                 `json:"webhook_id"`
	UserID        string                 `json:"user_id"`
	Event         string                 `json:"event"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Signature     string                 `json:"signature"`
	HTTPStatus    *int64                 `json:"http_status,omitempty"`
	ResponseBody  *string                `json:"response_body,omitempty"`
	DurationMS    int64                  `json:"duration_ms"`
	AttemptCount  int64                  `json:"attempt_count"`
	MaxRetries    int64                  `json:"max_retries"`
	NextRetryAtMS *int64                 `json:"next_retry_at_ms,omitempty"` // Unix timestamp in milliseconds
	Status        string                 `json:"status"`
	ErrorMessage  *string                `json:"error_message,omitempty"`
	CreatedAt     int64                  `json:"created_at"` // Unix timestamp in milliseconds
}

// WebhookDeliveryListResponse represents a list of webhook delivery attempts.
type WebhookDeliveryListResponse struct {
	Deliveries []WebhookDeliveryLog `json:"deliveries"`
	Total      int64                `json:"total"`
}
