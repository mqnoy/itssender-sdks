# @itssender/sdk

Official TypeScript and JavaScript SDK for **itssender**, the high-throughput notification and omnichannel messaging platform.

## Key Features

- **TypeScript-First**: Comprehensive types for all request payloads and response DTOs aligned with the OpenAPI 3.1.0 contract.
- **Dual Module Exports**: Native ESM (`import`) and CommonJS (`require`) builds with TypeScript `.d.ts` declaration maps.
- **Zero Runtime Dependencies**: Powered by native `fetch` and standard `node:crypto` APIs without external networking bloat.
- **Resilient Transport Engine**: Configurable request timeouts and automatic retry with exponential backoff and randomized jitter on network errors and HTTP 429/5xx responses.
- **Dual Authentication**: Seamless support for programmatic API Keys (`X-API-Key`) and tenant session Bearer Tokens (`Authorization: Bearer <token>`).
- **Standard Webhooks v1.0.0**: Built-in cryptographic verification for `webhook-id`, `webhook-timestamp`, and `webhook-signature` (`v1,<base64>`) with constant-time equality and replay defense.
- **Architectural Invariants**: Strict compliance with zero floating-point arithmetic (safe integers for currency and quotas) and hybrid millisecond Unix epoch timestamps.

---

## Installation

Install using `pnpm`:

```bash
pnpm add @itssender/sdk
```

---

## Quick Start

### 1. Initialization & Authentication

```typescript
import { ItssenderClient } from '@itssender/sdk';

// Initialize with an API Key
const client = new ItssenderClient({
  baseURL: 'https://itssender.imzalab.com', // Defaults to https://itssender.imzalab.com
  apiKey: 'its_live_a1b2c3d4e5f6g7h8...',
  timeout: 30000,                   // 30 seconds timeout
  maxRetries: 3,                    // Exponential jitter retries
});

// Or initialize with a Tenant Session JWT Bearer Token
const sessionClient = new ItssenderClient({
  baseURL: 'https://itssender.imzalab.com',
  bearerToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
});
```

Dynamic token updates are supported:

```typescript
client.setApiKey('its_live_new_key');
client.setBearerToken('new_jwt_session_token');
```

---

## Resource Reference & Examples

### 1. WhatsApp Devices (`client.devices`)

Manage WhatsApp client container instances, pairing, logs, and presence.

```typescript
// 1. List registered device instances
const devices = await client.devices.list();

// 2. Provision a new WhatsApp container
const newDevice = await client.devices.create({
  name: 'Primary Support Bot',
  nano_cpus: 500000000, // 0.5 CPU cores
  memory_mb: 256,       // 256 MB RAM
});

// 3. Fetch pairing QR code for WhatsApp Web linking
const qr = await client.devices.getQr(newDevice.id);
console.log('Scan QR Code:', qr.qr);

// 4. Request an 8-digit pairing code via phone number
const pairCode = await client.devices.pairCode(newDevice.id, {
  phone_number: '6281234567890',
});
console.log('Pairing Code:', pairCode.code);

// 5. Inspect container worker logs
const logs = await client.devices.getLogs(newDevice.id, { tail: 100 });
console.log(logs.logs);

// 6. Set presence (device online status or chat typing)
await client.devices.setDevicePresence(newDevice.id, { state: 'available' });
await client.devices.setChatPresence(newDevice.id, {
  to: '6281234567890',
  state: 'composing', // 'composing' | 'recording' | 'paused'
});

// 7. Restart or delete device
await client.devices.restart(newDevice.id);
await client.devices.delete(newDevice.id);
```

### 2. Direct Messaging (`client.messages`)

Dispatch messages, media attachments, reactions, polls, locations, and vCards.

```typescript
const deviceId = 'dev_01j9k2m3n4';

// 1. Send Text Message
const textMsg = await client.messages.sendText(deviceId, {
  to: '6281234567890',
  message: 'Your order #12345 has been confirmed.',
});

// 2. Send Media Message (Image, Document, Audio, Video, Sticker)
const mediaMsg = await client.messages.sendMedia(deviceId, {
  to: '6281234567890',
  media_type: 'image',
  url: 'https://example.com/receipt.jpg',
  caption: 'Transaction receipt',
});

// 3. Send Emoji Reaction
await client.messages.sendReaction(deviceId, {
  to: '6281234567890',
  target_message_id: textMsg.message_id,
  reaction: '👍',
});

// 4. Send Interactive Poll
await client.messages.sendPoll(deviceId, {
  to: '6281234567890',
  question: 'What time works best for delivery?',
  options: ['Morning (09:00 - 12:00)', 'Afternoon (13:00 - 17:00)'],
  max_answers: 1,
});

// 5. Send Location Pin
await client.messages.sendLocation(deviceId, {
  to: '6281234567890',
  latitude: -6.2088,
  longitude: 106.8456,
  name: 'Jakarta Headquarters',
  address: 'Jl. Jenderal Sudirman No. 1, Jakarta Pusat',
});

// 6. Send Contact vCard
await client.messages.sendContact(deviceId, {
  to: '6281234567890',
  display_name: 'Customer Support',
  vcard: 'BEGIN:VCARD\nVERSION:3.0\nFN:Customer Support\nTEL;TYPE=CELL:+628111222333\nEND:VCARD',
});

// 7. Edit or Revoke Message (Delete for Everyone)
await client.messages.edit(deviceId, {
  to: '6281234567890',
  target_message_id: textMsg.message_id,
  new_text: 'Your order #12345 has been completed.',
});

await client.messages.revoke(deviceId, {
  to: '6281234567890',
  target_message_id: textMsg.message_id,
});

// 8. Mark Chat as Read
await client.messages.markRead(deviceId, {
  chat_jid: '6281234567890@s.whatsapp.net',
});
```

### 3. WhatsApp Groups (`client.groups`)

```typescript
const deviceId = 'dev_01j9k2m3n4';

// 1. List groups
const groups = await client.groups.list(deviceId);

// 2. Create group
const newGroup = await client.groups.create(deviceId, {
  subject: 'VIP Customers Cohort',
  participants: ['6281234567890', '6289876543210'],
});

// 3. Update metadata
await client.groups.update(deviceId, newGroup.group_jid, {
  topic: 'Official discussion for VIP releases',
});

// 4. Manage participants (add, remove, promote, demote)
await client.groups.manageParticipants(deviceId, newGroup.group_jid, {
  action: 'promote',
  participants: ['6281234567890'],
});

// 5. Get and revoke invite link
const link = await client.groups.getInviteCode(deviceId, newGroup.group_jid);
console.log('Invite URL:', link.invite_url);

await client.groups.revokeInviteCode(deviceId, newGroup.group_jid);

// 6. Join or Leave group
await client.groups.join(deviceId, { code: 'ABC123XYZ' });
await client.groups.leave(deviceId, newGroup.group_jid);
```

### 4. Contacts (`client.contacts`)

```typescript
const deviceId = 'dev_01j9k2m3n4';

// Verify WhatsApp account registration
const check = await client.contacts.check(deviceId, {
  phones: ['6281234567890', '6280000000000'],
});
for (const item of check.results) {
  console.log(`${item.phone}: registered=${item.is_registered}`);
}

// Fetch avatar profile picture
const avatar = await client.contacts.getProfilePicture(deviceId, '6281234567890@s.whatsapp.net');
console.log('Avatar URL:', avatar.url);

// Update status bio
await client.contacts.updateStatus(deviceId, { status: 'Online via itssender SDK' });
```

### 5. High-Priority Transactional OTP (`client.otp`)

Bypasses broadcast delay jitter for sub-second transactional delivery.

```typescript
// 1. Dispatch OTP code
const otp = await client.otp.send({
  phone_number: '6281234567890',
  code_length: 6,
  ttl_ms: 300000,      // 5 minutes
  max_attempts: 3,
});
console.log(`OTP Session ID: ${otp.otp_id}, Expires At: ${otp.expires_at}`);

// 2. Verify OTP input with constant-time equality check
const verification = await client.otp.verify({
  otp_id: otp.otp_id,
  code: '654321',
});

if (verification.success) {
  console.log('OTP verified successfully!');
} else {
  console.log(`Failed. Remaining attempts: ${verification.remaining_attempts}`);
}

// 3. Resend OTP if needed
await client.otp.resend({ otp_id: otp.otp_id });
```

### 6. Message Templates & Rendering (`client.templates`)

```typescript
// 1. Create a rich-text template
const tpl = await client.templates.create({
  name: 'Security Alert',
  category: 'AUTHENTICATION',
  content_raw: 'Alert: Login from {{ip_address}} on {{device_model}}. If not you, code: {{code}}.',
  variables: ['ip_address', 'device_model', 'code'],
});

// 2. Render locally or from template ID
const renderedText = await client.templates.render(tpl.id, {
  ip_address: '198.51.100.4',
  device_model: 'MacBook Pro',
  code: '827104',
});
console.log(renderedText);
```

### 7. Bulk Marketing Campaigns (`client.campaigns`)

```typescript
// 1. Schedule a broadcast campaign with jitter controls
const campaign = await client.campaigns.create({
  name: 'Black Friday Broadcast',
  channel_type: 'WHATSMEOW',
  message_template: 'Flash Sale: 50% off storewide! Visit our site.',
  recipients: ['6281234567890', '6289876543210'],
  delay_min_seconds: 5,
  delay_max_seconds: 15,
});

// 2. Start campaign execution
await client.campaigns.start(campaign.id);

// 3. Fetch analytics and delivery status
const stats = await client.campaigns.getAnalytics(campaign.id);
console.log(`Sent: ${stats.sent_count}, Delivered: ${stats.delivered_count}`);
```

### 8. Analytics & Audit Logs (`client.analytics`)

```typescript
// Aggregate message metrics
const overview = await client.analytics.getOverview({
  from_timestamp: 1710000000000,
  to_timestamp: 1710100000000,
});
console.log(`Total messages processed: ${overview.total_messages}`);

// Paginated transmission audit logs
const logs = await client.analytics.getMessages({
  status: 'SENT',
  direction: 'OUTBOUND',
  limit: 50,
});
```

### 9. Webhooks Configuration (`client.webhooks`)

```typescript
// Register endpoint for real-time notifications
const webhook = await client.webhooks.create({
  url: 'https://api.yourdomain.com/webhooks/itssender',
  events: ['message.received', 'otp.verified', 'device.connected'],
  mask_pii: true,
});
console.log('Webhook Signing Secret:', webhook.secret);

// Trigger a test ping
const ping = await client.webhooks.test(webhook.id);
console.log('Ping status:', ping.status);
```

---

## Webhook Verification (Standard Webhooks v1.0.0)

Every outbound HTTP notification from itssender contains **Standard Webhooks v1.0.0** headers:
- `webhook-id`: Unique delivery ID
- `webhook-timestamp`: Unix epoch seconds integer string
- `webhook-signature`: Format `v1,<base64-encoded-hmac-sha256>`

### Express / Node.js Example

```typescript
import express from 'express';
import { verifyWebhookSignature, WebhookVerificationError } from '@itssender/sdk';

const app = express();
const WEBHOOK_SECRET = process.env.ITSSENDER_WEBHOOK_SECRET!; // 'whsec_...'

// Note: Use raw body buffer or string for HMAC computation
app.post('/webhook', express.raw({ type: 'application/json' }), (req, res) => {
  try {
    const rawPayload = req.body.toString('utf-8');

    // Automatically validates signature, timestamp drift (300s), and returns typed event
    const event = verifyWebhookSignature({
      payload: rawPayload,
      headers: req.headers,
      secret: WEBHOOK_SECRET,
      toleranceSeconds: 300, // 5 minutes replay protection
    });

    console.log(`Received event ${event.event} for delivery ${event.delivery_id}:`, event.data);

    // Process event...
    res.status(200).json({ received: true });
  } catch (err) {
    if (err instanceof WebhookVerificationError) {
      console.error(`Webhook verification failed [${err.reason}]:`, err.message);
      return res.status(401).json({ error: 'Invalid webhook signature' });
    }
    res.status(500).json({ error: 'Internal server error' });
  }
});
```

### Fastify Example

```typescript
fastify.post('/webhook', async (request, reply) => {
  try {
    const event = verifyWebhookSignature({
      payload: request.rawBody as string,
      headers: request.headers,
      secret: process.env.ITSSENDER_WEBHOOK_SECRET!,
    });

    return { success: true };
  } catch (err) {
    reply.status(401).send({ error: (err as Error).message });
  }
});
```

---

## Error Handling

The SDK exposes a structured, predictable error hierarchy:

```typescript
import {
  ItssenderError,
  ItssenderAPIError,
  ItssenderAuthenticationError,
  ItssenderNotFoundError,
  ItssenderValidationError,
  ItssenderRateLimitError,
  ItssenderTimeoutError,
  ItssenderNetworkError,
} from '@itssender/sdk';

try {
  await client.devices.get('dev_invalid');
} catch (error) {
  if (error instanceof ItssenderAuthenticationError) {
    console.error('Invalid API Key or Bearer Token');
  } else if (error instanceof ItssenderNotFoundError) {
    console.error('Resource not found (404)');
  } else if (error instanceof ItssenderRateLimitError) {
    console.error(`Rate limited! Retry after ${error.retryAfterSeconds} seconds`);
  } else if (error instanceof ItssenderTimeoutError) {
    console.error(`Request timed out after ${error.timeoutMs}ms`);
  } else if (error instanceof ItssenderAPIError) {
    console.error(`API Error ${error.statusCode} [${error.code}]: ${error.message}`);
  } else if (error instanceof ItssenderNetworkError) {
    console.error('Network connectivity error:', error.message);
  }
}
```

---

## Global Architectural Invariants

The SDK strictly enforces the platform's global invariants:

1. **Integer Precision for Currency and Quotas (Invariant 1)**: All quotas, delay durations, and monetary amounts are represented as safe integers (`number`). Zero floating-point operations.
2. **Hybrid Millisecond Unix Timestamps (Invariant 2)**: All date-time fields are millisecond timestamps (`number`, e.g. `1710000000000`).
3. **Zero Absolute Paths (Invariant 6)**: No local host filesystem paths are present in code, comments, or documentation.
4. **Strict Package Manager Hygiene (Invariant 7)**: Built and tested using `pnpm`.

---

## License

Apache-2.0
