# itssender Go SDK

The official Go SDK for the [itssender](https://github.com/mqnoy/itssender) WhatsApp Gateway & Omnichannel Messaging Platform.

Built with idiomatic Go standard library (`net/http`) with **zero external dependencies**.

---

## Architectural Invariants Enforced

- **Invariant 1 (Integer Precision):** Strict integer arithmetic for all quotas, message counters, and currency amounts (`int64`, zero floating-point arithmetic).
- **Invariant 2 (Hybrid Millisecond Timestamps):** All temporal fields and transit timestamps are represented as `int64` Unix millisecond timestamps.
- **Invariant 6 (Relative Path Portability):** Zero absolute filesystem paths.

---

## Installation

```bash
go get github.com/mqnoy/itssender-sdks/go

Consumers import the module with the package name `itssender`, for example: `import "github.com/mqnoy/itssender-sdks/go"` used as `itssender.NewClient(...)`.
```

---

## Quickstart

### 1. Initialize Client

The client supports functional configuration options:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mqnoy/itssender-sdks/go"
)

func main() {
	// Initialize client using an API Key (defaults to https://itssender.imzalab.com)
	client := itssender.NewClient(
		itssender.WithAPIKey("its_live_your_api_key_here"),
		itssender.WithBaseURL(itssender.DefaultBaseURL),
		itssender.WithTimeout(15*time.Second),
		itssender.WithRetry(3),
	)

	ctx := context.Background()

	// List provisioned WhatsApp devices
	devices, err := client.Devices.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	for _, d := range devices {
		fmt.Printf("Device ID: %s | Name: %s | Status: %s\n", d.ID, d.Name, d.Status)
	}
}
```

---

## Core Capabilities & Examples

### Messaging

#### Send Plain Text Message

```go
resp, err := client.Messages.SendText(ctx, "dev_01j9k2m3n4", itssender.SendTextMessageRequest{
	To:   "6281234567890",
	Text: "Hello from itssender Go SDK!",
})
if err != nil {
	log.Fatalf("Failed to send message: %v", err)
}
fmt.Printf("Message dispatched! ID: %s, Status: %s\n", resp.MessageID, resp.Status)
```

#### Send Media Attachment (Image, Document, Video, Audio)

```go
resp, err := client.Messages.SendMedia(ctx, "dev_01j9k2m3n4", itssender.SendMediaMessageRequest{
	To:        "6281234567890",
	MediaType: "image",
	MediaURL:  "https://example.com/invoice-receipt.png",
	Caption:   "Here is your monthly invoice receipt.",
})
```

#### Send Interactive Poll

```go
resp, err := client.Messages.SendPoll(ctx, "dev_01j9k2m3n4", itssender.SendPollRequest{
	To:       "6281234567890",
	Question: "Which day works best for our project sync?",
	Options:  []string{"Monday", "Wednesday", "Friday"},
})
```

#### Send Emoji Reaction

```go
_, err := client.Messages.SendReaction(ctx, "dev_01j9k2m3n4", itssender.SendReactionRequest{
	To:              "6281234567890",
	TargetMessageID: "3EB0C847E64B012",
	Emoji:           "🚀",
})
```

#### Send GPS Location

```go
_, err := client.Messages.SendLocation(ctx, "dev_01j9k2m3n4", itssender.SendLocationMessageRequest{
	To:        "6281234567890",
	Latitude:  -6.2088,
	Longitude: 106.8456,
	Name:      "Jakarta Central Hub",
	Address:   "Jl. Jend. Sudirman No. 1, Jakarta",
})
```

#### Revoke / Delete Message for Everyone

```go
_, err := client.Messages.Revoke(ctx, "dev_01j9k2m3n4", itssender.RevokeMessageRequest{
	TargetMessageID: "3EB0C847E64B012",
})
```

---

### Device Management

```go
// Provision a new WhatsApp device container
device, err := client.Devices.Create(ctx, itssender.CreateDeviceRequest{
	Name:     "Support Desk WhatsApp",
	NanoCPUs: 500000000, // 0.5 CPU
	MemoryMB: 256,
})

// Retrieve QR code for pairing
qr, err := client.Devices.GetQR(ctx, device.ID)
fmt.Println("Scan QR code:", qr.Code)

// Or link via 8-digit Phone Pairing Code
pair, err := client.Devices.PairCode(ctx, device.ID, itssender.PairCodeRequest{
	PhoneNumber: "6281234567890",
})
fmt.Println("Enter Pairing Code:", pair.PairCode)

// Fetch runtime container metrics
metrics, err := client.Devices.GetState(ctx, device.ID)
fmt.Printf("Memory usage: %d MB | Goroutines: %d\n", metrics.MemoryUsageMB, metrics.Goroutines)
```

---

### WhatsApp Groups

```go
// Create a new WhatsApp group
group, err := client.Groups.Create(ctx, "dev_01j9k2m3n4", itssender.CreateGroupRequest{
	Subject:      "Alpha Team Operations",
	Participants: []string{"6281234567890", "6289876543210"},
})

// Add or promote members
_, err = client.Groups.ManageParticipants(ctx, "dev_01j9k2m3n4", group.GroupJID, itssender.UpdateGroupParticipantsRequest{
	Action:       "promote",
	Participants: []string{"6281234567890"},
})

// Fetch invite link
link, err := client.Groups.GetInviteCode(ctx, "dev_01j9k2m3n4", group.GroupJID)
fmt.Println("Group Link:", link.InviteURL)
```

---

### Contact Verification

```go
// Check WhatsApp registration before sending messages
check, err := client.Contacts.Check(ctx, "dev_01j9k2m3n4", itssender.CheckContactsRequest{
	Phones: []string{"6281234567890", "6289999999999"},
})
for _, res := range check.Results {
	fmt.Printf("Phone %s is registered: %t\n", res.Phone, res.IsRegistered)
}
```

---

### High-Priority Transactional OTP

```go
// Dispatch OTP with 5-minute TTL
otpResp, err := client.OTP.Send(ctx, itssender.SendOTPRequest{
	PhoneNumber: "6281234567890",
	CodeLength:  6,
	TTLMS:       300000, // 5 minutes in milliseconds
	MaxAttempts: 3,
})

// Verify code provided by user
verifyResp, err := client.OTP.Verify(ctx, itssender.VerifyOTPRequest{
	OTPID: otpResp.OTPID,
	Code:  "491823",
})
if verifyResp.Success {
	fmt.Println("OTP successfully validated!")
}
```

---

### Rich-Text Message Templates

```go
// Render template with variable interpolation
rendered, err := client.Templates.Render(ctx, "tmpl_welcome", map[string]string{
	"name":     "Sarah",
	"order_id": "#10492",
})
fmt.Println(rendered)
```

---

## Webhook Verification

itssender dispatches outbound webhook notifications conforming to **Standard Webhooks v1.0.0** (`webhook-id`, `webhook-timestamp`, `webhook-signature` using HMAC-SHA256) and includes legacy headers for backwards compatibility (`X-Itssender-Signature`).

### Production HTTP Handler Example

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mqnoy/itssender-sdks/go"
)

const WebhookSecret = "whsec_C307e53aA+vL4o82pZ8W7aFjF6t1YQ=="

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 1. Verify Standard Webhooks v1.0.0 (HMAC-SHA256 + 5-minute replay window tolerance)
	valid, err := itssender.VerifyStandardWebhook(payload, r.Header, WebhookSecret)
	if err != nil || !valid {
		log.Printf("Webhook signature verification failed: %v", err)
		http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
		return
	}

	// 2. Process validated event payload
	fmt.Printf("Received authenticated event payload: %s\n", string(payload))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"received":true}`))
}

func main() {
	http.HandleFunc("/api/webhooks", webhookHandler)
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### Legacy Webhook Verification

If your pipeline relies on `X-Itssender-Signature`:

```go
sig := r.Header.Get("X-Itssender-Signature")
if !itssender.VerifyLegacyWebhook(payload, sig, WebhookSecret) {
	http.Error(w, "Invalid legacy signature", http.StatusUnauthorized)
	return
}
```

---

## Error Handling

All failed API requests return a structured `*itssender.APIError`:

```go
_, err := client.Devices.Get(ctx, "dev_invalid")
if err != nil {
	if apiErr, ok := err.(*itssender.APIError); ok {
		fmt.Printf("Status: %d | Code: %s | Message: %s\n", apiErr.StatusCode, apiErr.Code, apiErr.Message)
	} else {
		fmt.Printf("General error: %v\n", err)
	}
}
```

---

## License

Apache-2.0. See [LICENSE](../LICENSE) for details.
