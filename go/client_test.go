package itssender

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_OptionsAndHeaders(t *testing.T) {
	var capturedAuth, capturedAPIKey, capturedContentType, capturedUA string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedAPIKey = r.Header.Get("X-API-Key")
		capturedContentType = r.Header.Get("Content-Type")
		capturedUA = r.Header.Get("User-Agent")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":    "usr_123",
				"email": "test@example.com",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithAPIKey("its_live_testkey123"),
		WithBearerToken("jwt_test_token_456"),
		WithTimeout(5*time.Second),
		WithRetry(2),
	)

	user, err := client.Auth.Me(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != "usr_123" || user.Email != "test@example.com" {
		t.Errorf("unexpected user data: %+v", user)
	}
	if capturedAPIKey != "its_live_testkey123" {
		t.Errorf("expected X-API-Key 'its_live_testkey123', got '%s'", capturedAPIKey)
	}
	if capturedAuth != "Bearer jwt_test_token_456" {
		t.Errorf("expected Authorization 'Bearer jwt_test_token_456', got '%s'", capturedAuth)
	}
	if capturedContentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got '%s'", capturedContentType)
	}
	if capturedUA != "itssender-go/1.0.0" {
		t.Errorf("expected User-Agent itssender-go/1.0.0, got '%s'", capturedUA)
	}
}

func TestClient_RetryOnServerError(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":     "dev_123",
				"name":   "WhatsApp Primary",
				"status": "CONNECTED",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithRetry(3),
	)

	dev, err := client.Devices.Get(context.Background(), "dev_123")
	if err != nil {
		t.Fatalf("expected retry to succeed, got error: %v", err)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
	if dev.Name != "WhatsApp Primary" {
		t.Errorf("unexpected device name: %s", dev.Name)
	}
}

func TestClient_ErrorDecoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "DEVICE_NOT_FOUND",
				"message": "The requested device does not exist",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	_, err := client.Devices.Get(context.Background(), "dev_nonexistent")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != "DEVICE_NOT_FOUND" {
		t.Errorf("expected code DEVICE_NOT_FOUND, got %s", apiErr.Code)
	}
	if apiErr.Message != "The requested device does not exist" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
}

func TestAuthService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/auth/register":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"token": "reg_token_123",
					"user":  map[string]interface{}{"id": "u1", "email": "test@reg.com"},
				},
			})
		case "/api/v1/auth/login":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"token": "login_token_123",
					"user":  map[string]interface{}{"id": "u1", "email": "test@reg.com"},
				},
			})
		case "/api/v1/auth/config":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"registration_enabled": true,
					"version":              "1.0.0",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	reg, err := client.Auth.Register(context.Background(), RegisterRequest{Email: "test@reg.com", Password: "pass"})
	if err != nil || reg.Token != "reg_token_123" {
		t.Fatalf("register failed: %v", err)
	}

	login, err := client.Auth.Login(context.Background(), LoginRequest{Email: "test@reg.com", Password: "pass"})
	if err != nil || login.Token != "login_token_123" {
		t.Fatalf("login failed: %v", err)
	}

	cfg, err := client.Auth.GetConfig(context.Background())
	if err != nil || !cfg.RegistrationEnabled {
		t.Fatalf("getConfig failed: %v", err)
	}
}

func TestAPIKeysService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/user/keys" && r.Method == http.MethodGet {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "k1", "name": "Production Key", "key_prefix": "its_live_123"},
				},
			})
			return
		}
		if r.URL.Path == "/api/v1/user/keys" && r.Method == http.MethodPost {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":         "k2",
					"name":       "New Key",
					"key_prefix": "its_live_456",
					"key":        "its_live_456_fullplaintextsecret",
				},
			})
			return
		}
		if r.URL.Path == "/api/v1/user/keys/k1" && r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	keys, err := client.APIKeys.List(context.Background())
	if err != nil || len(keys) != 1 || keys[0].ID != "k1" {
		t.Fatalf("list API keys failed: %v", err)
	}

	created, err := client.APIKeys.Create(context.Background(), CreateApiKeyRequest{Name: "New Key"})
	if err != nil || created.Key == "" {
		t.Fatalf("create API key failed: %v", err)
	}

	if err := client.APIKeys.Delete(context.Background(), "k1"); err != nil {
		t.Fatalf("delete API key failed: %v", err)
	}
}

func TestDevicesService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/devices" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "dev_1", "name": "Device 1", "status": "CONNECTED"},
				},
			})
		case r.URL.Path == "/api/v1/devices" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "dev_created", "name": "New Device", "status": "INITIALIZING",
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/qr":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"code": "2@qrCodeString123", "timeout_seconds": 60,
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/pair-code":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"pair_code": "ABCD-1234", "timeout_seconds": 120,
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/restart":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		case r.URL.Path == "/api/v1/devices/dev_1/logs":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"logs": []string{"starting container...", "connected"},
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/presence/device":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		case r.URL.Path == "/api/v1/devices/dev_1/presence/chat":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		case r.URL.Path == "/api/v1/devices/dev_1/usage":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"cpu_usage_pct": 12, "memory_usage_mb": 64, "status": "HEALTHY",
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1" && r.Method == http.MethodDelete:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	list, err := client.Devices.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("Devices.List failed: %v", err)
	}

	created, err := client.Devices.Create(context.Background(), CreateDeviceRequest{Name: "New Device"})
	if err != nil || created.ID != "dev_created" {
		t.Fatalf("Devices.Create failed: %v", err)
	}

	qr, err := client.Devices.GetQR(context.Background(), "dev_1")
	if err != nil || qr.Code != "2@qrCodeString123" {
		t.Fatalf("Devices.GetQR failed: %v", err)
	}

	pair, err := client.Devices.PairCode(context.Background(), "dev_1", PairCodeRequest{PhoneNumber: "6281234567890"})
	if err != nil || pair.PairCode != "ABCD-1234" {
		t.Fatalf("Devices.PairCode failed: %v", err)
	}

	if err := client.Devices.Restart(context.Background(), "dev_1"); err != nil {
		t.Fatalf("Devices.Restart failed: %v", err)
	}

	logs, err := client.Devices.GetLogs(context.Background(), "dev_1", 50)
	if err != nil || len(logs.Logs) != 2 {
		t.Fatalf("Devices.GetLogs failed: %v", err)
	}

	if err := client.Devices.SetPresence(context.Background(), "dev_1", DevicePresenceRequest{State: "available"}); err != nil {
		t.Fatalf("Devices.SetPresence failed: %v", err)
	}

	if err := client.Devices.SetChatPresence(context.Background(), "dev_1", ChatPresenceRequest{State: "composing"}); err != nil {
		t.Fatalf("Devices.SetChatPresence failed: %v", err)
	}

	state, err := client.Devices.GetState(context.Background(), "dev_1")
	if err != nil || state.CPUUsagePct != 12 {
		t.Fatalf("Devices.GetState failed: %v", err)
	}

	if err := client.Devices.Delete(context.Background(), "dev_1"); err != nil {
		t.Fatalf("Devices.Delete failed: %v", err)
	}
}

func TestMessagesService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/devices/dev_1/messages/text":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_txt_001", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/messages/media":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_media_002", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/messages/revoke":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true, "message": "revoked", "timestamp": int64(1710000000000),
			})
		case "/api/v1/devices/dev_1/messages/reaction":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_react_003", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/messages/poll":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_poll_004", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/messages/location":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_loc_005", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/messages/contact":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_con_006", "status": "SENT", "timestamp": int64(1710000000000),
				},
			})
		case "/api/v1/devices/dev_1/chat/read":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true, "message": "marked read", "timestamp": int64(1710000000000),
			})
		case "/api/v1/devices/dev_1/messages/edit":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"message_id": "msg_txt_001", "status": "EDITED", "timestamp": int64(1710000000000),
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	txt, err := client.Messages.SendText(context.Background(), "dev_1", SendTextMessageRequest{
		To: "6281234567890", Text: "Hello, World!",
	})
	if err != nil || txt.MessageID != "msg_txt_001" {
		t.Fatalf("Messages.SendText failed: %v", err)
	}

	media, err := client.Messages.SendMedia(context.Background(), "dev_1", SendMediaMessageRequest{
		To: "6281234567890", MediaType: "image", MediaURL: "https://example.com/img.png",
	})
	if err != nil || media.MessageID != "msg_media_002" {
		t.Fatalf("Messages.SendMedia failed: %v", err)
	}

	revoke, err := client.Messages.Revoke(context.Background(), "dev_1", RevokeMessageRequest{
		TargetMessageID: "msg_txt_001",
	})
	if err != nil || !revoke.Success {
		t.Fatalf("Messages.Revoke failed: %v", err)
	}

	react, err := client.Messages.SendReaction(context.Background(), "dev_1", SendReactionRequest{
		TargetMessageID: "msg_txt_001", Emoji: "👍",
	})
	if err != nil || react.MessageID != "msg_react_003" {
		t.Fatalf("Messages.SendReaction failed: %v", err)
	}

	poll, err := client.Messages.SendPoll(context.Background(), "dev_1", SendPollRequest{
		To: "6281234567890", Question: "Lunch preference?", Options: []string{"Pizza", "Burger"},
	})
	if err != nil || poll.MessageID != "msg_poll_004" {
		t.Fatalf("Messages.SendPoll failed: %v", err)
	}

	loc, err := client.Messages.SendLocation(context.Background(), "dev_1", SendLocationMessageRequest{
		To: "6281234567890", Latitude: -6.2088, Longitude: 106.8456, Name: "Jakarta",
	})
	if err != nil || loc.MessageID != "msg_loc_005" {
		t.Fatalf("Messages.SendLocation failed: %v", err)
	}

	con, err := client.Messages.SendContact(context.Background(), "dev_1", SendContactMessageRequest{
		To: "6281234567890", DisplayName: "John Doe", VCard: "BEGIN:VCARD...",
	})
	if err != nil || con.MessageID != "msg_con_006" {
		t.Fatalf("Messages.SendContact failed: %v", err)
	}

	read, err := client.Messages.MarkRead(context.Background(), "dev_1", MarkReadRequest{
		ChatJID: "6281234567890@s.whatsapp.net",
	})
	if err != nil || !read.Success {
		t.Fatalf("Messages.MarkRead failed: %v", err)
	}

	edit, err := client.Messages.Edit(context.Background(), "dev_1", EditMessageRequest{
		TargetMessageID: "msg_txt_001", NewText: "Updated message",
	})
	if err != nil || edit.Status != "EDITED" {
		t.Fatalf("Messages.Edit failed: %v", err)
	}
}

func TestGroupsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/devices/dev_1/groups" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"jid": "120363@g.us", "subject": "Dev Team"},
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/groups/120363@g.us/metadata":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"jid": "120363@g.us", "subject": "Dev Team", "topic": "Daily Standup",
					},
				})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
			}
		case r.URL.Path == "/api/v1/devices/dev_1/groups/create":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"group_jid": "120363@g.us", "owner_jid": "6281@s.whatsapp.net", "timestamp": int64(1710000000000),
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/groups/120363@g.us/participants":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "participants updated"})
		case r.URL.Path == "/api/v1/devices/dev_1/groups/120363@g.us/invite-link":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"invite_code": "ABCD123", "invite_url": "https://chat.whatsapp.com/ABCD123",
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/groups/join":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "joined"})
		case r.URL.Path == "/api/v1/devices/dev_1/groups/120363@g.us/leave":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "left"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	groups, err := client.Groups.List(context.Background(), "dev_1")
	if err != nil || len(groups) != 1 {
		t.Fatalf("Groups.List failed: %v", err)
	}

	created, err := client.Groups.Create(context.Background(), "dev_1", CreateGroupRequest{
		Subject: "Dev Team", Participants: []string{"6281234567890"},
	})
	if err != nil || created.GroupJID != "120363@g.us" {
		t.Fatalf("Groups.Create failed: %v", err)
	}

	grp, err := client.Groups.Get(context.Background(), "dev_1", "120363@g.us")
	if err != nil || grp.Subject != "Dev Team" {
		t.Fatalf("Groups.Get failed: %v", err)
	}

	upd, err := client.Groups.Update(context.Background(), "dev_1", "120363@g.us", UpdateGroupMetadataRequest{Topic: "New Topic"})
	if err != nil || !upd.Success {
		t.Fatalf("Groups.Update failed: %v", err)
	}

	parts, err := client.Groups.ManageParticipants(context.Background(), "dev_1", "120363@g.us", UpdateGroupParticipantsRequest{
		Action: "add", Participants: []string{"6289999999999"},
	})
	if err != nil || !parts.Success {
		t.Fatalf("Groups.ManageParticipants failed: %v", err)
	}

	link, err := client.Groups.GetInviteCode(context.Background(), "dev_1", "120363@g.us")
	if err != nil || link.InviteCode != "ABCD123" {
		t.Fatalf("Groups.GetInviteCode failed: %v", err)
	}

	join, err := client.Groups.Join(context.Background(), "dev_1", JoinGroupRequest{Code: "ABCD123"})
	if err != nil || !join.Success {
		t.Fatalf("Groups.Join failed: %v", err)
	}

	leave, err := client.Groups.Leave(context.Background(), "dev_1", "120363@g.us")
	if err != nil || !leave.Success {
		t.Fatalf("Groups.Leave failed: %v", err)
	}
}

func TestContactsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/devices/dev_1/contacts/check":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"success": true,
					"results": []map[string]interface{}{
						{"phone": "6281234567890", "is_registered": true, "jid": "6281234567890@s.whatsapp.net"},
					},
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/contacts/6281234567890@s.whatsapp.net/avatar":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"success": true, "url": "https://pps.whatsapp.net/avatar.jpg",
				},
			})
		case r.URL.Path == "/api/v1/devices/dev_1/profile/status":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true, "data": map[string]interface{}{"status": "Available"},
				})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "status updated"})
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	check, err := client.Contacts.Check(context.Background(), "dev_1", CheckContactsRequest{Phones: []string{"6281234567890"}})
	if err != nil || len(check.Results) != 1 || !check.Results[0].IsRegistered {
		t.Fatalf("Contacts.Check failed: %v", err)
	}

	avatar, err := client.Contacts.GetProfilePicture(context.Background(), "dev_1", "6281234567890@s.whatsapp.net")
	if err != nil || avatar.URL != "https://pps.whatsapp.net/avatar.jpg" {
		t.Fatalf("Contacts.GetProfilePicture failed: %v", err)
	}

	status, err := client.Contacts.GetStatus(context.Background(), "dev_1")
	if err != nil || status["status"] != "Available" {
		t.Fatalf("Contacts.GetStatus failed: %v", err)
	}

	upd, err := client.Contacts.UpdateStatus(context.Background(), "dev_1", UpdateStatusRequest{Status: "Busy"})
	if err != nil || !upd.Success {
		t.Fatalf("Contacts.UpdateStatus failed: %v", err)
	}
}

func TestCampaignsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/campaigns" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "cmp_1", "name": "Black Friday Blast", "status": "COMPLETED", "total_count": 100},
				},
			})
		case r.URL.Path == "/api/v1/campaigns" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "cmp_new", "name": "New Campaign", "status": "PENDING",
				},
			})
		case r.URL.Path == "/api/v1/campaigns/cmp_1" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "cmp_1", "name": "Black Friday Blast", "status": "COMPLETED", "sent_count": 98,
				},
			})
		case r.URL.Path == "/api/v1/campaigns/cmp_1/start":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "campaign started"})
		case r.URL.Path == "/api/v1/campaigns/cmp_1/archive":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "campaign paused"})
		case r.URL.Path == "/api/v1/campaigns/cmp_1" && r.Method == http.MethodDelete:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	list, err := client.Campaigns.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("Campaigns.List failed: %v", err)
	}

	created, err := client.Campaigns.Create(context.Background(), CreateCampaignRequest{
		Name: "New Campaign", ChannelType: "WHATSMEOW", MessageTemplate: "Hello {{name}}",
	})
	if err != nil || created.ID != "cmp_new" {
		t.Fatalf("Campaigns.Create failed: %v", err)
	}

	detail, err := client.Campaigns.Get(context.Background(), "cmp_1")
	if err != nil || detail.SentCount != 98 {
		t.Fatalf("Campaigns.Get failed: %v", err)
	}

	start, err := client.Campaigns.Start(context.Background(), "cmp_1")
	if err != nil || !start.Success {
		t.Fatalf("Campaigns.Start failed: %v", err)
	}

	pause, err := client.Campaigns.Pause(context.Background(), "cmp_1")
	if err != nil || !pause.Success {
		t.Fatalf("Campaigns.Pause failed: %v", err)
	}

	if err := client.Campaigns.Delete(context.Background(), "cmp_1"); err != nil {
		t.Fatalf("Campaigns.Delete failed: %v", err)
	}
}

func TestAnalyticsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/analytics/usage":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"messages_sent": 5000, "messages_failed": 12, "active_devices": 3, "quota_remaining": 45000,
				},
			})
		case "/api/v1/analytics/logs":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "log_1", "recipient": "6281234567890", "status": "DELIVERED"},
				},
			})
		case "/api/v1/devices/dev_1/usage":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"cpu_usage_pct": 5, "memory_usage_mb": 42,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	overview, err := client.Analytics.GetOverview(context.Background())
	if err != nil || overview.MessagesSent != 5000 {
		t.Fatalf("Analytics.GetOverview failed: %v", err)
	}

	logs, err := client.Analytics.GetMessages(context.Background(), 10)
	if err != nil || len(logs) != 1 || logs[0].Status != "DELIVERED" {
		t.Fatalf("Analytics.GetMessages failed: %v", err)
	}

	devUsage, err := client.Analytics.GetDevices(context.Background(), "dev_1")
	if err != nil || devUsage.MemoryUsageMB != 42 {
		t.Fatalf("Analytics.GetDevices failed: %v", err)
	}
}

func TestOTPService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/otp/send":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"otp_id": "otp_123", "phone_number": "6281234567890", "status": "DISPATCHED",
				},
			})
		case "/api/v1/otp/verify":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"success": true, "otp_id": "otp_123", "status": "VERIFIED",
				},
			})
		case "/api/v1/otp/otp_123":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "otp_123", "phone_number": "6281234567890", "status": "VERIFIED",
				},
			})
		case "/api/v1/otp":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "otp_123", "phone_number": "6281234567890", "status": "VERIFIED"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	sent, err := client.OTP.Send(context.Background(), SendOTPRequest{PhoneNumber: "6281234567890"})
	if err != nil || sent.OTPID != "otp_123" {
		t.Fatalf("OTP.Send failed: %v", err)
	}

	verified, err := client.OTP.Verify(context.Background(), VerifyOTPRequest{OTPID: "otp_123", Code: "123456"})
	if err != nil || !verified.Success {
		t.Fatalf("OTP.Verify failed: %v", err)
	}

	resend, err := client.OTP.Resend(context.Background(), SendOTPRequest{PhoneNumber: "6281234567890"})
	if err != nil || resend.OTPID != "otp_123" {
		t.Fatalf("OTP.Resend failed: %v", err)
	}

	detail, err := client.OTP.Get(context.Background(), "otp_123")
	if err != nil || detail.Status != "VERIFIED" {
		t.Fatalf("OTP.Get failed: %v", err)
	}

	list, err := client.OTP.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("OTP.List failed: %v", err)
	}
}

func TestTemplatesService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/templates" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"templates": []map[string]interface{}{
						{"id": "tmpl_1", "name": "Welcome Message", "content_raw": "Hi {{name}}, your order {{order_id}} is ready!"},
					},
					"total": 1,
				},
			})
		case r.URL.Path == "/api/v1/templates" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "tmpl_new", "name": "New Template", "content_raw": "Hello {{name}}",
				},
			})
		case r.URL.Path == "/api/v1/templates/tmpl_1" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "tmpl_1", "name": "Welcome Message", "content_raw": "Hi {{name}}, your order {{order_id}} is ready!",
				},
			})
		case r.URL.Path == "/api/v1/templates/tmpl_1" && r.Method == http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "tmpl_1", "name": "Welcome Message Updated",
				},
			})
		case r.URL.Path == "/api/v1/templates/tmpl_1" && r.Method == http.MethodDelete:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	list, err := client.Templates.List(context.Background())
	if err != nil || len(list.Templates) != 1 {
		t.Fatalf("Templates.List failed: %v", err)
	}

	created, err := client.Templates.Create(context.Background(), TemplateCreateRequest{
		Name: "New Template", Category: "UTILITY", ContentRaw: "Hello {{name}}",
	})
	if err != nil || created.ID != "tmpl_new" {
		t.Fatalf("Templates.Create failed: %v", err)
	}

	tmpl, err := client.Templates.Get(context.Background(), "tmpl_1")
	if err != nil || tmpl.ID != "tmpl_1" {
		t.Fatalf("Templates.Get failed: %v", err)
	}

	rendered, err := client.Templates.Render(context.Background(), "tmpl_1", map[string]string{
		"name":     "Alice",
		"order_id": "#8842",
	})
	if err != nil || rendered != "Hi Alice, your order #8842 is ready!" {
		t.Fatalf("Templates.Render failed: got '%s', err: %v", rendered, err)
	}

	newName := "Welcome Message Updated"
	upd, err := client.Templates.Update(context.Background(), "tmpl_1", TemplateUpdateRequest{Name: &newName})
	if err != nil || upd.Name != newName {
		t.Fatalf("Templates.Update failed: %v", err)
	}

	if err := client.Templates.Delete(context.Background(), "tmpl_1"); err != nil {
		t.Fatalf("Templates.Delete failed: %v", err)
	}
}

func TestWebhooksService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/webhooks" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": []map[string]interface{}{
					{"id": "wh_1", "url": "https://example.com/webhook", "events": []string{"message.received"}},
				},
			})
		case r.URL.Path == "/api/v1/webhooks" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "wh_new", "url": "https://example.com/webhook", "secret": "whsec_12345",
				},
			})
		case r.URL.Path == "/api/v1/webhooks/wh_1" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id": "wh_1", "url": "https://example.com/webhook", "secret": "whsec_12345",
				},
			})
		case r.URL.Path == "/api/v1/webhooks/wh_1/test":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true, "data": map[string]interface{}{"status": "delivered", "http_status": 200},
			})
		case r.URL.Path == "/api/v1/webhooks/deliveries":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"deliveries": []map[string]interface{}{
						{"id": "del_1", "event": "message.received", "status": "DELIVERED"},
					},
					"total": 1,
				},
			})
		case r.URL.Path == "/api/v1/webhooks/deliveries/del_1/retry":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true, "data": map[string]interface{}{"status": "QUEUED"},
			})
		case r.URL.Path == "/api/v1/webhooks/wh_1" && r.Method == http.MethodDelete:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	list, err := client.Webhooks.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("Webhooks.List failed: %v", err)
	}

	created, err := client.Webhooks.Create(context.Background(), CreateWebhookRequest{
		URL: "https://example.com/webhook", Events: []string{"message.received"},
	})
	if err != nil || created.ID != "wh_new" {
		t.Fatalf("Webhooks.Create failed: %v", err)
	}

	wh, err := client.Webhooks.Get(context.Background(), "wh_1")
	if err != nil || wh.Secret != "whsec_12345" {
		t.Fatalf("Webhooks.Get failed: %v", err)
	}

	testResult, err := client.Webhooks.Test(context.Background(), "wh_1")
	if err != nil || testResult["status"] != "delivered" {
		t.Fatalf("Webhooks.Test failed: %v", err)
	}

	deliveries, err := client.Webhooks.ListDeliveries(context.Background())
	if err != nil || len(deliveries.Deliveries) != 1 {
		t.Fatalf("Webhooks.ListDeliveries failed: %v", err)
	}

	retryResult, err := client.Webhooks.RetryDelivery(context.Background(), "del_1")
	if err != nil || retryResult["status"] != "QUEUED" {
		t.Fatalf("Webhooks.RetryDelivery failed: %v", err)
	}

	if err := client.Webhooks.Delete(context.Background(), "wh_1"); err != nil {
		t.Fatalf("Webhooks.Delete failed: %v", err)
	}
}
