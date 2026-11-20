import { describe, expect, it } from 'vitest';
import {
  computeLegacyWebhookSignature,
  computeStandardWebhookSignature,
  constantTimeCompare,
  parseSecretBytes,
  verifyLegacyWebhook,
  verifyStandardWebhook,
  verifyWebhookSignature,
} from '../src/crypto/webhook.js';
import { WebhookVerificationError } from '../src/errors.js';

describe('Webhook Crypto & Verification', () => {
  const secretKey = 'whsec_MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE';
  const rawSecret = 'super_secret_signing_key_12345';
  const deliveryId = 'whd_01j9k2m3n4p5';
  const samplePayload = JSON.stringify({
    event: 'message.received',
    delivery_id: deliveryId,
    timestamp: 1710000000000,
    data: {
      message_id: 'msg_987654321',
      sender: '6281234567890',
      text: 'Hello from WhatsApp test',
    },
  });

  describe('parseSecretBytes', () => {
    it('strips whsec_ prefix and handles base64', () => {
      const buf = parseSecretBytes('whsec_dGVzdF9zZWNyZXQ=');
      expect(buf.toString('utf-8')).toBe('test_secret');
    });

    it('falls back to utf-8 when base64 is invalid', () => {
      const buf = parseSecretBytes('whsec_not-valid-base64!!!');
      expect(buf.toString('utf-8')).toBe('not-valid-base64!!!');
    });

    it('throws error when secret is empty', () => {
      expect(() => parseSecretBytes('')).toThrow(WebhookVerificationError);
    });
  });

  describe('constantTimeCompare', () => {
    it('returns true for identical strings and buffers', () => {
      expect(constantTimeCompare('secret123', 'secret123')).toBe(true);
      expect(
        constantTimeCompare(Buffer.from('hello world'), Buffer.from('hello world'))
      ).toBe(true);
    });

    it('returns false for different strings or different lengths', () => {
      expect(constantTimeCompare('secret123', 'secret124')).toBe(false);
      expect(constantTimeCompare('short', 'longer_string')).toBe(false);
    });
  });

  describe('Standard Webhooks v1.0.0', () => {
    it('computes valid standard signature starting with v1,', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);
      expect(sig.startsWith('v1,')).toBe(true);
    });

    it('verifies valid standard webhook signature', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);

      const isValid = verifyStandardWebhook({
        id: deliveryId,
        timestamp: nowSec,
        payload: samplePayload,
        signature: sig,
        secret: secretKey,
        toleranceSeconds: 300,
      });

      expect(isValid).toBe(true);
    });

    it('verifies standard signature with multiple signatures in header', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const validSig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);
      const multiSigHeader = `v1,oldInvalidSig123 ${validSig} v1,anotherDummySig`;

      const isValid = verifyStandardWebhook({
        id: deliveryId,
        timestamp: nowSec,
        payload: samplePayload,
        signature: multiSigHeader,
        secret: secretKey,
      });

      expect(isValid).toBe(true);
    });

    it('rejects signature when secret does not match', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);

      expect(() =>
        verifyStandardWebhook({
          id: deliveryId,
          timestamp: nowSec,
          payload: samplePayload,
          signature: sig,
          secret: 'whsec_wrong_secret_key_12345678',
        })
      ).toThrow(WebhookVerificationError);
    });

    it('rejects payload when timestamp is outside tolerance (replay protection)', () => {
      const pastSec = Math.floor(Date.now() / 1000) - 305; // 5m 5s ago
      const sig = computeStandardWebhookSignature(deliveryId, pastSec, samplePayload, secretKey);

      expect(() =>
        verifyStandardWebhook({
          id: deliveryId,
          timestamp: pastSec,
          payload: samplePayload,
          signature: sig,
          secret: secretKey,
          toleranceSeconds: 300,
        })
      ).toThrowError(/outside tolerance/);
    });

    it('allows tolerance bypass when toleranceSeconds is 0', () => {
      const pastSec = Math.floor(Date.now() / 1000) - 10000;
      const sig = computeStandardWebhookSignature(deliveryId, pastSec, samplePayload, secretKey);

      const isValid = verifyStandardWebhook({
        id: deliveryId,
        timestamp: pastSec,
        payload: samplePayload,
        signature: sig,
        secret: secretKey,
        toleranceSeconds: 0,
      });

      expect(isValid).toBe(true);
    });
  });

  describe('Legacy X-Itssender-Signature', () => {
    it('computes and verifies sha256=<hex> legacy signature', () => {
      const sig = computeLegacyWebhookSignature(samplePayload, rawSecret);
      expect(sig.startsWith('sha256=')).toBe(true);

      const isValid = verifyLegacyWebhook({
        payload: samplePayload,
        signature: sig,
        secret: rawSecret,
      });

      expect(isValid).toBe(true);
    });

    it('verifies raw hex legacy signature without sha256= prefix', () => {
      const sig = computeLegacyWebhookSignature(samplePayload, rawSecret);
      const rawHex = sig.replace('sha256=', '');

      const isValid = verifyLegacyWebhook({
        payload: samplePayload,
        signature: rawHex,
        secret: rawSecret,
      });

      expect(isValid).toBe(true);
    });

    it('rejects invalid legacy signature', () => {
      expect(() =>
        verifyLegacyWebhook({
          payload: samplePayload,
          signature: 'sha256=deadbeefcafebabe00112233445566778899aabbccddeeff',
          secret: rawSecret,
        })
      ).toThrow(WebhookVerificationError);
    });
  });

  describe('verifyWebhookSignature (Unified Handler)', () => {
    it('successfully parses Standard Webhooks headers and returns typed envelope', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);

      const headers = {
        'webhook-id': deliveryId,
        'webhook-timestamp': String(nowSec),
        'webhook-signature': sig,
        'content-type': 'application/json',
      };

      const event = verifyWebhookSignature<{ message_id: string; sender: string; text: string }>({
        payload: samplePayload,
        headers,
        secret: secretKey,
      });

      expect(event.event).toBe('message.received');
      expect(event.delivery_id).toBe(deliveryId);
      expect(event.data.text).toBe('Hello from WhatsApp test');
    });

    it('works with Headers instance', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, samplePayload, secretKey);

      const headers = new Headers();
      headers.set('webhook-id', deliveryId);
      headers.set('webhook-timestamp', String(nowSec));
      headers.set('webhook-signature', sig);

      const event = verifyWebhookSignature({
        payload: samplePayload,
        headers,
        secret: secretKey,
      });

      expect(event.event).toBe('message.received');
      expect(event.delivery_id).toBe(deliveryId);
    });

    it('successfully parses legacy X-Itssender-Signature header', () => {
      const sig = computeLegacyWebhookSignature(samplePayload, rawSecret);

      const headers = {
        'X-Itssender-Signature': sig,
      };

      const event = verifyWebhookSignature({
        payload: samplePayload,
        headers,
        secret: rawSecret,
      });

      expect(event.event).toBe('message.received');
      expect(event.delivery_id).toBe(deliveryId);
    });

    it('throws MISSING_HEADERS when neither standard nor legacy headers exist', () => {
      expect(() =>
        verifyWebhookSignature({
          payload: samplePayload,
          headers: {},
          secret: secretKey,
        })
      ).toThrowError(/Missing webhook signature headers/);
    });

    it('throws INVALID_PAYLOAD on non-JSON payload', () => {
      const nowSec = Math.floor(Date.now() / 1000);
      const nonJson = 'not valid json text';
      const sig = computeStandardWebhookSignature(deliveryId, nowSec, nonJson, secretKey);

      const headers = {
        'webhook-id': deliveryId,
        'webhook-timestamp': String(nowSec),
        'webhook-signature': sig,
      };

      expect(() =>
        verifyWebhookSignature({
          payload: nonJson,
          headers,
          secret: secretKey,
        })
      ).toThrowError(/Invalid JSON payload/);
    });
  });
});
