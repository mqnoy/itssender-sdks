package itssender

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestVerifyStandardWebhook_Success(t *testing.T) {
	secret := "whsec_C307e53aA+vL4o82pZ8W7aFjF6t1YQ=="
	webhookID := "del_01HXYZ123456"
	now := time.Now().Unix()
	payload := []byte(`{"event":"otp.verified","data":{"phone":"6281234567890"}}`)

	sig, err := ComputeStandardSignature(secret, webhookID, now, payload)
	if err != nil {
		t.Fatalf("unexpected error computing signature: %v", err)
	}

	headers := http.Header{}
	headers.Set("webhook-id", webhookID)
	headers.Set("webhook-timestamp", strconv.FormatInt(now, 10))
	headers.Set("webhook-signature", sig)

	valid, err := VerifyStandardWebhook(payload, headers, secret)
	if err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}
	if !valid {
		t.Errorf("expected valid to be true")
	}
}

func TestVerifyStandardWebhook_MultipleSignatures(t *testing.T) {
	secret := "whsec_C307e53aA+vL4o82pZ8W7aFjF6t1YQ=="
	webhookID := "del_01HXYZ123456"
	now := time.Now().Unix()
	payload := []byte(`{"event":"otp.verified"}`)

	validSig, _ := ComputeStandardSignature(secret, webhookID, now, payload)
	combinedSig := fmt.Sprintf("v1,invalidSig123 %s v1,anotherFakeSig", validSig)

	headers := http.Header{}
	headers.Set("webhook-id", webhookID)
	headers.Set("webhook-timestamp", strconv.FormatInt(now, 10))
	headers.Set("webhook-signature", combinedSig)

	valid, err := VerifyStandardWebhook(payload, headers, secret)
	if err != nil {
		t.Fatalf("expected verification to succeed with multiple signatures: %v", err)
	}
	if !valid {
		t.Errorf("expected valid to be true")
	}
}

func TestVerifyStandardWebhook_TamperedPayload(t *testing.T) {
	secret := "whsec_C307e53aA+vL4o82pZ8W7aFjF6t1YQ=="
	webhookID := "del_01HXYZ123456"
	now := time.Now().Unix()
	originalPayload := []byte(`{"amount":100}`)
	tamperedPayload := []byte(`{"amount":99999}`)

	sig, _ := ComputeStandardSignature(secret, webhookID, now, originalPayload)

	headers := http.Header{}
	headers.Set("webhook-id", webhookID)
	headers.Set("webhook-timestamp", strconv.FormatInt(now, 10))
	headers.Set("webhook-signature", sig)

	valid, err := VerifyStandardWebhook(tamperedPayload, headers, secret)
	if err == nil {
		t.Fatalf("expected error for tampered payload, got nil")
	}
	if valid {
		t.Errorf("expected valid to be false for tampered payload")
	}
}

func TestVerifyStandardWebhook_ExpiredTimestamp(t *testing.T) {
	secret := "whsec_C307e53aA+vL4o82pZ8W7aFjF6t1YQ=="
	webhookID := "del_01HXYZ123456"
	// 10 minutes in past (exceeds default 5 minutes tolerance)
	expiredTime := time.Now().Unix() - 600
	payload := []byte(`{"event":"ping"}`)

	sig, _ := ComputeStandardSignature(secret, webhookID, expiredTime, payload)

	headers := http.Header{}
	headers.Set("webhook-id", webhookID)
	headers.Set("webhook-timestamp", strconv.FormatInt(expiredTime, 10))
	headers.Set("webhook-signature", sig)

	valid, err := VerifyStandardWebhook(payload, headers, secret)
	if err == nil {
		t.Fatalf("expected error for expired timestamp, got nil")
	}
	if valid {
		t.Errorf("expected valid to be false for expired timestamp")
	}
}

func TestVerifyStandardWebhook_MissingHeaders(t *testing.T) {
	secret := "raw_secret_key_123"
	payload := []byte(`{"event":"ping"}`)

	// Missing webhook-id
	headers := http.Header{}
	headers.Set("webhook-timestamp", "1710000000")
	headers.Set("webhook-signature", "v1,abc")
	valid, err := VerifyStandardWebhook(payload, headers, secret)
	if err != ErrMissingWebhookID || valid {
		t.Errorf("expected ErrMissingWebhookID, got %v", err)
	}

	// Missing webhook-timestamp
	headers = http.Header{}
	headers.Set("webhook-id", "del_123")
	headers.Set("webhook-signature", "v1,abc")
	valid, err = VerifyStandardWebhook(payload, headers, secret)
	if err != ErrMissingWebhookTimestamp || valid {
		t.Errorf("expected ErrMissingWebhookTimestamp, got %v", err)
	}

	// Missing webhook-signature
	headers = http.Header{}
	headers.Set("webhook-id", "del_123")
	headers.Set("webhook-timestamp", "1710000000")
	valid, err = VerifyStandardWebhook(payload, headers, secret)
	if err != ErrMissingWebhookSignature || valid {
		t.Errorf("expected ErrMissingWebhookSignature, got %v", err)
	}

	// Invalid timestamp integer
	headers.Set("webhook-signature", "v1,abc")
	headers.Set("webhook-timestamp", "not-a-number")
	valid, err = VerifyStandardWebhook(payload, headers, secret)
	if err != ErrInvalidWebhookTimestamp || valid {
		t.Errorf("expected ErrInvalidWebhookTimestamp, got %v", err)
	}
}

func TestVerifyLegacyWebhook(t *testing.T) {
	secret := "test_legacy_secret"
	payload := []byte(`{"event":"message.sent","id":"msg_123"}`)

	sig := ComputeSignature(secret, payload) // produces "sha256=<hex>"

	if !VerifyLegacyWebhook(payload, sig, secret) {
		t.Errorf("expected VerifyLegacyWebhook to succeed with standard signature")
	}

	// Test without "sha256=" prefix
	rawHex := sig[7:]
	if !VerifyLegacyWebhook(payload, rawHex, secret) {
		t.Errorf("expected VerifyLegacyWebhook to succeed with raw hex signature")
	}

	// Tampered payload
	tampered := []byte(`{"event":"message.sent","id":"msg_999"}`)
	if VerifyLegacyWebhook(tampered, sig, secret) {
		t.Errorf("expected VerifyLegacyWebhook to fail with tampered payload")
	}

	// Wrong secret
	if VerifyLegacyWebhook(payload, sig, "wrong_secret") {
		t.Errorf("expected VerifyLegacyWebhook to fail with incorrect secret")
	}

	// Empty inputs
	if VerifyLegacyWebhook(payload, "", secret) || VerifyLegacyWebhook(payload, sig, "") {
		t.Errorf("expected VerifyLegacyWebhook to return false on empty signature or secret")
	}
}
