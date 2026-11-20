package itssender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrMissingWebhookID is returned when the webhook-id header is absent.
	ErrMissingWebhookID = errors.New("itssender: missing webhook-id header")
	// ErrMissingWebhookTimestamp is returned when the webhook-timestamp header is absent.
	ErrMissingWebhookTimestamp = errors.New("itssender: missing webhook-timestamp header")
	// ErrMissingWebhookSignature is returned when the webhook-signature header is absent.
	ErrMissingWebhookSignature = errors.New("itssender: missing webhook-signature header")
	// ErrInvalidWebhookTimestamp is returned when webhook-timestamp cannot be parsed as an integer.
	ErrInvalidWebhookTimestamp = errors.New("itssender: invalid webhook-timestamp format")
	// ErrWebhookTimestampExpired is returned when the timestamp is outside the replay tolerance window.
	ErrWebhookTimestampExpired = errors.New("itssender: webhook-timestamp outside replay tolerance window")
	// ErrInvalidWebhookSignature is returned when the computed signature does not match any provided signature.
	ErrInvalidWebhookSignature = errors.New("itssender: invalid webhook-signature")
)

const (
	// DefaultWebhookTolerance is the default replay window tolerance (5 minutes).
	DefaultWebhookTolerance = 5 * time.Minute
)

// VerifyStandardWebhook validates an incoming webhook request conforming to Standard Webhooks v1.0.0.
// It verifies:
// 1. Presence of `webhook-id`, `webhook-timestamp`, and `webhook-signature` headers.
// 2. That `webhook-timestamp` is within 5 minutes of current server time to prevent replay attacks.
// 3. Constant-time HMAC-SHA256 signature verification over "${webhook-id}.${webhook-timestamp}.${payload}".
//
// The secret parameter supports both "whsec_<base64>" and raw secret strings.
func VerifyStandardWebhook(payload []byte, headers http.Header, secret string) (bool, error) {
	return VerifyStandardWebhookWithTolerance(payload, headers, secret, DefaultWebhookTolerance)
}

// VerifyStandardWebhookWithTolerance validates a Standard Webhook with a custom replay tolerance.
// If tolerance <= 0, timestamp replay window validation is bypassed.
func VerifyStandardWebhookWithTolerance(payload []byte, headers http.Header, secret string, tolerance time.Duration) (bool, error) {
	webhookID := getHeaderCaseInsensitive(headers, "webhook-id")
	if webhookID == "" {
		return false, ErrMissingWebhookID
	}

	timestampStr := getHeaderCaseInsensitive(headers, "webhook-timestamp")
	if timestampStr == "" {
		return false, ErrMissingWebhookTimestamp
	}

	sigHeader := getHeaderCaseInsensitive(headers, "webhook-signature")
	if sigHeader == "" {
		return false, ErrMissingWebhookSignature
	}

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false, ErrInvalidWebhookTimestamp
	}

	// Validate timestamp replay tolerance if configured
	if tolerance > 0 {
		now := time.Now().Unix()
		sec := ts
		if sec > 100000000000 {
			// In case timestamp was provided in milliseconds
			sec = sec / 1000
		}
		diff := now - sec
		tolSec := int64(tolerance.Seconds())
		if diff < -tolSec || diff > tolSec {
			return false, fmt.Errorf("%w: drift of %d seconds exceeds tolerance of %d seconds", ErrWebhookTimestampExpired, diff, tolSec)
		}
	}

	secretBytes := parseSecretBytes(secret)
	signedContent := fmt.Sprintf("%s.%s.%s", webhookID, timestampStr, string(payload))

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signedContent))
	expectedBase64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// webhook-signature header can contain multiple space-delimited signatures (e.g. "v1,sig1 v1,sig2")
	signatures := strings.Fields(sigHeader)
	for _, sigToken := range signatures {
		sigToken = strings.TrimSpace(sigToken)
		if strings.HasPrefix(sigToken, "v1,") {
			candidate := strings.TrimPrefix(sigToken, "v1,")
			if hmac.Equal([]byte(candidate), []byte(expectedBase64)) {
				return true, nil
			}
		}
	}

	return false, ErrInvalidWebhookSignature
}

// VerifyLegacyWebhook validates an incoming webhook request signed with legacy X-Itssender-Signature.
// The signature format is expected to be "sha256=<hex>" or raw "<hex>".
// Returns true if valid, false otherwise.
func VerifyLegacyWebhook(payload []byte, signature string, secret string) bool {
	if signature == "" || secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedHex := hex.EncodeToString(mac.Sum(nil))

	cleanSig := strings.TrimSpace(signature)
	cleanSig = strings.TrimPrefix(cleanSig, "sha256=")

	return hmac.Equal([]byte(cleanSig), []byte(expectedHex))
}

// ComputeStandardSignature generates a Standard Webhooks v1.0.0 signature string formatted as "v1,<base64>".
func ComputeStandardSignature(secret string, webhookID string, timestampSec int64, payload []byte) (string, error) {
	secretBytes := parseSecretBytes(secret)
	signedContent := fmt.Sprintf("%s.%d.%s", webhookID, timestampSec, string(payload))

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signedContent))
	base64Sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return "v1," + base64Sig, nil
}

// ComputeSignature calculates the legacy HMAC-SHA256 signature formatted as "sha256=<hex>".
func ComputeSignature(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// parseSecretBytes decodes "whsec_<base64>" secrets or falls back to raw bytes.
func parseSecretBytes(secret string) []byte {
	rawSecret := strings.TrimPrefix(secret, "whsec_")
	if b, err := base64.StdEncoding.DecodeString(rawSecret); err == nil && len(b) > 0 {
		return b
	}
	if b, err := base64.RawStdEncoding.DecodeString(rawSecret); err == nil && len(b) > 0 {
		return b
	}
	return []byte(rawSecret)
}

// getHeaderCaseInsensitive searches for a header key ignoring case variations.
func getHeaderCaseInsensitive(headers http.Header, key string) string {
	if val := headers.Get(key); val != "" {
		return val
	}
	target := strings.ToLower(key)
	for k, v := range headers {
		if strings.ToLower(k) == target && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}
