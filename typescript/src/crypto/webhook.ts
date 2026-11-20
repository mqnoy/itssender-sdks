import { createHmac, timingSafeEqual } from 'node:crypto';
import { WebhookVerificationError } from '../errors.js';
import type { WebhookEnvelope } from '../types/webhooks.js';

export interface VerifyWebhookOptions {
  /**
   * The raw request payload body as a string, Buffer, or Uint8Array.
   */
  payload: string | Buffer | Uint8Array;

  /**
   * Request headers map or Headers object.
   */
  headers: Record<string, string | string[] | undefined> | Headers;

  /**
   * The webhook endpoint secret key (e.g. `whsec_...` or raw secret).
   */
  secret: string;

  /**
   * Maximum acceptable timestamp drift in seconds (replay attack protection).
   * Defaults to 300 seconds (5 minutes). Set to 0 to disable tolerance checking.
   */
  toleranceSeconds?: number;
}

export interface StandardWebhookVerifyParams {
  id: string;
  timestamp: string | number;
  payload: string | Buffer | Uint8Array;
  signature: string;
  secret: string;
  toleranceSeconds?: number;
}

export interface LegacyWebhookVerifyParams {
  payload: string | Buffer | Uint8Array;
  signature: string;
  secret: string;
}

/**
 * Parses secret bytes according to Standard Webhooks v1.0.0 specification.
 * Strips 'whsec_' prefix and attempts Base64 decoding; falls back to UTF-8 bytes.
 */
export function parseSecretBytes(secret: string): Buffer {
  if (!secret) {
    throw new WebhookVerificationError('Webhook secret cannot be empty', 'INVALID_SECRET');
  }

  const raw = secret.startsWith('whsec_') ? secret.slice(6) : secret;

  // Try standard base64 decoding
  try {
    const buf = Buffer.from(raw, 'base64');
    // Verify that re-encoding matches or raw was valid base64
    if (buf.length > 0 && buf.toString('base64').replace(/=+$/, '') === raw.replace(/=+$/, '')) {
      return buf;
    }
  } catch {
    // Fall back to raw utf-8 bytes
  }

  return Buffer.from(raw, 'utf-8');
}

/**
 * Computes a Standard Webhooks v1.0.0 HMAC-SHA256 signature string: `v1,<base64>`.
 *
 * @param webhookId - Unique delivery ID from `webhook-id` header
 * @param timestamp - Unix epoch timestamp in seconds from `webhook-timestamp` header
 * @param payload - Raw request body
 * @param secret - Webhook endpoint secret
 */
export function computeStandardWebhookSignature(
  webhookId: string,
  timestamp: string | number,
  payload: string | Buffer | Uint8Array,
  secret: string
): string {
  const secretBytes = parseSecretBytes(secret);
  const payloadStr = typeof payload === 'string' ? payload : Buffer.from(payload).toString('utf-8');
  const signedContent = `${webhookId}.${timestamp}.${payloadStr}`;

  const hmac = createHmac('sha256', secretBytes);
  hmac.update(signedContent, 'utf-8');
  const base64Digest = hmac.digest('base64');

  return `v1,${base64Digest}`;
}

/**
 * Computes a legacy HMAC-SHA256 signature string: `sha256=<hex>`.
 *
 * @param payload - Raw request body
 * @param secret - Webhook endpoint secret
 */
export function computeLegacyWebhookSignature(
  payload: string | Buffer | Uint8Array,
  secret: string
): string {
  const payloadBuf = typeof payload === 'string' ? Buffer.from(payload, 'utf-8') : Buffer.from(payload);
  const hmac = createHmac('sha256', Buffer.from(secret, 'utf-8'));
  hmac.update(payloadBuf);
  const hexDigest = hmac.digest('hex');

  return `sha256=${hexDigest}`;
}

/**
 * Constant-time string or buffer comparison to prevent timing attacks.
 */
export function constantTimeCompare(a: string | Buffer, b: string | Buffer): boolean {
  const bufA = typeof a === 'string' ? Buffer.from(a, 'utf-8') : a;
  const bufB = typeof b === 'string' ? Buffer.from(b, 'utf-8') : b;

  if (bufA.length !== bufB.length) {
    return false;
  }

  try {
    return timingSafeEqual(bufA, bufB);
  } catch {
    // Fallback constant-time comparison
    let result = 0;
    for (let i = 0; i < bufA.length; i++) {
      result |= bufA[i] ^ bufB[i];
    }
    return result === 0;
  }
}

/**
 * Verifies a Standard Webhooks v1.0.0 payload signature and timestamp tolerance.
 */
export function verifyStandardWebhook(params: StandardWebhookVerifyParams): boolean {
  const { id, timestamp, payload, signature, secret, toleranceSeconds = 300 } = params;

  if (!id || !timestamp || !signature) {
    throw new WebhookVerificationError(
      'Missing required Standard Webhooks fields (id, timestamp, signature)',
      'MISSING_HEADERS'
    );
  }

  // Validate timestamp tolerance if enabled
  const tsNum = typeof timestamp === 'number' ? timestamp : parseInt(timestamp, 10);
  if (isNaN(tsNum)) {
    throw new WebhookVerificationError('Invalid webhook-timestamp header', 'TIMESTAMP_OUT_OF_TOLERANCE');
  }

  if (toleranceSeconds > 0) {
    const nowSec = Math.floor(Date.now() / 1000);
    const diff = Math.abs(nowSec - tsNum);
    if (diff > toleranceSeconds) {
      throw new WebhookVerificationError(
        `Webhook timestamp ${tsNum} is outside tolerance (${toleranceSeconds}s, diff was ${diff}s)`,
        'TIMESTAMP_OUT_OF_TOLERANCE'
      );
    }
  }

  const expectedSignature = computeStandardWebhookSignature(id, timestamp, payload, secret);
  const expectedSigValue = expectedSignature.slice(3); // strip 'v1,'

  // Standard Webhooks signature header may contain multiple signatures separated by space
  const candidates = signature.split(/\s+/).filter(Boolean);
  let matched = false;

  for (const candidate of candidates) {
    let cleanCandidate = candidate.trim();
    if (cleanCandidate.startsWith('v1,')) {
      cleanCandidate = cleanCandidate.slice(3);
    }

    if (constantTimeCompare(expectedSigValue, cleanCandidate)) {
      matched = true;
      break;
    }
  }

  if (!matched) {
    throw new WebhookVerificationError('Standard Webhook signature mismatch', 'INVALID_SIGNATURE');
  }

  return true;
}

/**
 * Verifies a legacy `X-Itssender-Signature` payload.
 */
export function verifyLegacyWebhook(params: LegacyWebhookVerifyParams): boolean {
  const { payload, signature, secret } = params;

  if (!signature) {
    throw new WebhookVerificationError('Missing legacy signature', 'MISSING_HEADERS');
  }

  const expectedSignature = computeLegacyWebhookSignature(payload, secret);
  let cleanExpected = expectedSignature;
  let cleanActual = signature.trim();

  if (cleanExpected.startsWith('sha256=')) {
    cleanExpected = cleanExpected.slice(7);
  }
  if (cleanActual.startsWith('sha256=')) {
    cleanActual = cleanActual.slice(7);
  }

  if (!constantTimeCompare(cleanExpected, cleanActual)) {
    throw new WebhookVerificationError('Legacy signature mismatch', 'INVALID_SIGNATURE');
  }

  return true;
}

/**
 * Helper to extract header value from Record or Headers object in case-insensitive manner.
 */
function getHeader(
  headers: Record<string, string | string[] | undefined> | Headers,
  headerName: string
): string | undefined {
  if (typeof (headers as Headers).get === 'function') {
    const val = (headers as Headers).get(headerName);
    return val ?? undefined;
  }

  const normalized = headerName.toLowerCase();
  for (const [key, value] of Object.entries(headers)) {
    if (key.toLowerCase() === normalized) {
      if (Array.isArray(value)) {
        return value[0];
      }
      return value;
    }
  }

  return undefined;
}

/**
 * Unified webhook verification function.
 * Automatically inspects headers for Standard Webhooks v1.0.0 (`webhook-id`, `webhook-timestamp`, `webhook-signature`)
 * or legacy `X-Itssender-Signature`. Returns the parsed typed event envelope or throws WebhookVerificationError.
 *
 * @param options - Webhook verification parameters
 * @returns Typed WebhookEnvelope containing the parsed JSON data
 * @throws WebhookVerificationError if signature verification fails or headers are missing
 */
export function verifyWebhookSignature<T = unknown>(
  options: VerifyWebhookOptions
): WebhookEnvelope<T> {
  const { payload, headers, secret, toleranceSeconds = 300 } = options;

  if (!secret) {
    throw new WebhookVerificationError('Webhook secret is required', 'INVALID_SECRET');
  }

  const payloadStr =
    typeof payload === 'string' ? payload : Buffer.from(payload).toString('utf-8');

  // Check for Standard Webhooks headers
  const stdId = getHeader(headers, 'webhook-id');
  const stdTimestamp = getHeader(headers, 'webhook-timestamp');
  const stdSignature = getHeader(headers, 'webhook-signature');

  if (stdId && stdTimestamp && stdSignature) {
    verifyStandardWebhook({
      id: stdId,
      timestamp: stdTimestamp,
      payload: payloadStr,
      signature: stdSignature,
      secret,
      toleranceSeconds,
    });
  } else {
    // Check for Legacy headers
    const legacySig =
      getHeader(headers, 'X-Itssender-Signature') || getHeader(headers, 'x-itssender-signature');
    if (legacySig) {
      verifyLegacyWebhook({
        payload: payloadStr,
        signature: legacySig,
        secret,
      });
    } else {
      throw new WebhookVerificationError(
        'Missing webhook signature headers (neither Standard Webhooks headers nor X-Itssender-Signature found)',
        'MISSING_HEADERS'
      );
    }
  }

  // Parse payload JSON
  try {
    return JSON.parse(payloadStr) as WebhookEnvelope<T>;
  } catch (err) {
    throw new WebhookVerificationError(
      `Invalid JSON payload: ${(err as Error)?.message || String(err)}`,
      'INVALID_PAYLOAD'
    );
  }
}
