import { describe, expect, it, vi } from 'vitest';
import { ItssenderClient } from '../src/client.js';
import {
  ItssenderAPIError,
  ItssenderAuthenticationError,
  ItssenderNotFoundError,
  ItssenderRateLimitError,
  ItssenderTimeoutError,
  ItssenderValidationError,
} from '../src/errors.js';

describe('ItssenderClient', () => {
  it('initializes with default baseURL https://itssender.imzalab.com', () => {
    const client = new ItssenderClient();
    expect(client.getBaseURL()).toBe('https://itssender.imzalab.com');
  });

  it('strips trailing slashes from custom baseURL', () => {
    const client = new ItssenderClient({ baseURL: 'https://gateway.example.com/' });
    expect(client.getBaseURL()).toBe('https://gateway.example.com');
  });

  it('correctly categorizes token shorthand', async () => {
    let capturedHeaders: Headers | undefined;
    const mockFetch = vi.fn().mockImplementation(async (_url: string, init: RequestInit) => {
      capturedHeaders = new Headers(init.headers as HeadersInit);
      return new Response(JSON.stringify({ success: true, data: { id: 'usr_123' } }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      });
    });

    // API Key token
    const clientApiKey = new ItssenderClient({
      token: 'its_live_secret_key_12345',
      fetch: mockFetch,
    });
    await clientApiKey.auth.me();
    expect(capturedHeaders?.get('X-API-Key')).toBe('its_live_secret_key_12345');
    expect(capturedHeaders?.get('Authorization')).toBeNull();

    // Bearer token
    const clientJwt = new ItssenderClient({
      token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy',
      fetch: mockFetch,
    });
    await clientJwt.auth.me();
    expect(capturedHeaders?.get('Authorization')).toBe(
      'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy'
    );
  });

  it('updates authentication dynamically with setApiKey and setBearerToken', async () => {
    let capturedHeaders: Headers | undefined;
    const mockFetch = vi.fn().mockImplementation(async (_url: string, init: RequestInit) => {
      capturedHeaders = new Headers(init.headers as HeadersInit);
      return new Response(JSON.stringify({ success: true, data: {} }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      });
    });

    const client = new ItssenderClient({ fetch: mockFetch });
    client.setApiKey('its_test_key_abc');
    await client.devices.list();
    expect(capturedHeaders?.get('X-API-Key')).toBe('its_test_key_abc');

    client.setBearerToken('new_jwt_session_token');
    await client.devices.list();
    expect(capturedHeaders?.get('Authorization')).toBe('Bearer new_jwt_session_token');
  });

  describe('HTTP Client Error Handling & Status Codes', () => {
    it('throws ItssenderAuthenticationError on 401 / 403', async () => {
      const mockFetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            success: false,
            error: { code: 'UNAUTHORIZED', message: 'Invalid credentials' },
          }),
          { status: 401, headers: { 'content-type': 'application/json' } }
        )
      );

      const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });
      await expect(client.auth.me()).rejects.toThrow(ItssenderAuthenticationError);
    });

    it('throws ItssenderNotFoundError on 404', async () => {
      const mockFetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            success: false,
            error: { code: 'NOT_FOUND', message: 'Device not found' },
          }),
          { status: 404, headers: { 'content-type': 'application/json' } }
        )
      );

      const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });
      await expect(client.devices.get('dev_missing')).rejects.toThrow(ItssenderNotFoundError);
    });

    it('throws ItssenderValidationError on 400', async () => {
      const mockFetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            success: false,
            error: { code: 'INVALID_REQUEST', message: 'Missing required field' },
          }),
          { status: 400, headers: { 'content-type': 'application/json' } }
        )
      );

      const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });
      await expect(
        client.devices.create({ name: '' })
      ).rejects.toThrow(ItssenderValidationError);
    });

    it('throws ItssenderRateLimitError on 429 and parses Retry-After header', async () => {
      const mockFetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            success: false,
            error: { code: 'RATE_LIMIT_EXCEEDED', message: 'Too many requests' },
          }),
          {
            status: 429,
            headers: { 'content-type': 'application/json', 'Retry-After': '12' },
          }
        )
      );

      const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });
      try {
        await client.devices.list();
        expect.unreachable();
      } catch (err) {
        expect(err).toBeInstanceOf(ItssenderRateLimitError);
        expect((err as ItssenderRateLimitError).statusCode).toBe(429);
        expect((err as ItssenderRateLimitError).retryAfterSeconds).toBe(12);
      }
    });

    it('throws ItssenderAPIError on enveloped failure even if HTTP status is 200', async () => {
      const mockFetch = vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({
            success: false,
            error: { code: 'DEVICE_OFFLINE', message: 'Target client disconnected' },
          }),
          { status: 200, headers: { 'content-type': 'application/json' } }
        )
      );

      const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });
      await expect(
        client.messages.sendText('dev_123', { to: '628123', message: 'Hi' })
      ).rejects.toThrow(ItssenderAPIError);
    });

    it('throws ItssenderTimeoutError when request times out', async () => {
      const mockFetch = vi.fn().mockImplementation((_url: string, init: RequestInit) => {
        return new Promise((_resolve, reject) => {
          init.signal?.addEventListener('abort', () => {
            const err = new Error('The operation was aborted');
            err.name = 'AbortError';
            reject(err);
          });
        });
      });

      const client = new ItssenderClient({
        fetch: mockFetch,
        timeout: 20, // 20ms timeout
        maxRetries: 0,
      });

      await expect(client.devices.list()).rejects.toThrow(ItssenderTimeoutError);
    });
  });

  describe('Automatic Retry with Jitter', () => {
    it('retries on 500 status code and succeeds on subsequent attempt', async () => {
      let callCount = 0;
      const mockFetch = vi.fn().mockImplementation(async () => {
        callCount++;
        if (callCount === 1) {
          return new Response(JSON.stringify({ error: 'Server error' }), { status: 500 });
        }
        return new Response(JSON.stringify({ success: true, data: [{ id: 'dev_ok' }] }), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        });
      });

      const client = new ItssenderClient({
        fetch: mockFetch,
        maxRetries: 2,
        baseDelayMs: 10,
      });

      const devices = await client.devices.list();
      expect(callCount).toBe(2);
      expect(devices).toEqual([{ id: 'dev_ok' }]);
    });

    it('fails after exhausting maxRetries on persistent 503', async () => {
      let callCount = 0;
      const mockFetch = vi.fn().mockImplementation(async () => {
        callCount++;
        return new Response(JSON.stringify({ error: 'Gateway unavailable' }), { status: 503 });
      });

      const client = new ItssenderClient({
        fetch: mockFetch,
        maxRetries: 2,
        baseDelayMs: 5,
      });

      await expect(client.devices.list()).rejects.toThrow(ItssenderAPIError);
      expect(callCount).toBe(3); // 1 initial + 2 retries
    });
  });

  describe('Modular Resource Endpoints', () => {
    let lastUrl = '';
    let lastMethod = '';
    let lastBody: unknown;

    const mockFetch = vi.fn().mockImplementation(async (url: string, init: RequestInit) => {
      lastUrl = url;
      lastMethod = init.method || 'GET';
      lastBody = init.body ? JSON.parse(init.body as string) : undefined;

      return new Response(
        JSON.stringify({
          success: true,
          data: { id: 'mock_resource_id', success: true },
        }),
        { status: 200, headers: { 'content-type': 'application/json' } }
      );
    });

    const client = new ItssenderClient({ fetch: mockFetch, maxRetries: 0 });

    it('auth: register, login, me, refreshToken, getConfig', async () => {
      await client.auth.register({ email: 'test@example.com', password: 'secretpassword' });
      expect(lastUrl).toContain('/api/v1/auth/register');
      expect(lastMethod).toBe('POST');

      await client.auth.login({ email: 'test@example.com', password: 'secretpassword' });
      expect(lastUrl).toContain('/api/v1/auth/login');

      await client.auth.me();
      expect(lastUrl).toContain('/api/v1/auth/me');
      expect(lastMethod).toBe('GET');

      await client.auth.refreshToken('mock_token');
      expect(lastUrl).toContain('/api/v1/auth/refresh');

      await client.auth.getConfig();
      expect(lastUrl).toContain('/api/v1/auth/config');
    });

    it('apiKeys: list, create, delete', async () => {
      await client.apiKeys.list();
      expect(lastUrl).toContain('/api/v1/user/keys');
      expect(lastMethod).toBe('GET');

      await client.apiKeys.create({ name: 'Integration Server' });
      expect(lastUrl).toContain('/api/v1/user/keys');
      expect(lastMethod).toBe('POST');
      expect(lastBody).toEqual({ name: 'Integration Server' });

      await client.apiKeys.delete('key_12345');
      expect(lastUrl).toContain('/api/v1/user/keys/key_12345');
      expect(lastMethod).toBe('DELETE');
    });

    it('devices: list, get, create, delete, getQr, pairCode, restart, getLogs, updateSettings, setPresence, getState', async () => {
      await client.devices.list();
      expect(lastUrl).toContain('/api/v1/devices');

      await client.devices.get('dev_001');
      expect(lastUrl).toContain('/api/v1/devices/dev_001');

      await client.devices.create({ name: 'Office Bot', nano_cpus: 500000000 });
      expect(lastUrl).toContain('/api/v1/devices');
      expect(lastMethod).toBe('POST');

      await client.devices.delete('dev_001');
      expect(lastUrl).toContain('/api/v1/devices/dev_001');
      expect(lastMethod).toBe('DELETE');

      await client.devices.getQr('dev_001');
      expect(lastUrl).toContain('/api/v1/devices/dev_001/qr');

      await client.devices.pairCode('dev_001', { phone_number: '628123456789' });
      expect(lastUrl).toContain('/api/v1/devices/dev_001/pair-code');
      expect(lastMethod).toBe('POST');

      await client.devices.restart('dev_001');
      expect(lastUrl).toContain('/api/v1/devices/dev_001/restart');
      expect(lastMethod).toBe('POST');

      await client.devices.getLogs('dev_001', { tail: 50 });
      expect(lastUrl).toContain('/api/v1/devices/dev_001/logs?tail=50');

      await client.devices.updateSettings('dev_001', { name: 'Renamed Bot' });
      expect(lastUrl).toContain('/api/v1/devices/dev_001');
      expect(lastMethod).toBe('PUT');

      // Device presence
      await client.devices.setPresence('dev_001', { state: 'available' });
      expect(lastUrl).toContain('/api/v1/devices/dev_001/presence/device');
      expect(lastMethod).toBe('POST');

      // Chat presence
      await client.devices.setPresence('dev_001', {
        to: '628123',
        state: 'composing',
      });
      expect(lastUrl).toContain('/api/v1/devices/dev_001/presence/chat');
      expect(lastMethod).toBe('POST');

      await client.devices.getState('dev_001');
      expect(lastUrl).toContain('/api/v1/devices/dev_001/usage');
    });

    it('messages: sendText, sendMedia, sendReaction, sendPoll, sendLocation, sendContact, edit, revoke, markRead', async () => {
      const dev = 'dev_msg_01';

      await client.messages.sendText(dev, { to: '628123', message: 'Hello' });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/text`);

      await client.messages.sendMedia(dev, {
        to: '628123',
        media_type: 'image',
        url: 'https://example.com/pic.png',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/media`);

      await client.messages.sendReaction(dev, {
        to: '628123',
        target_message_id: 'msg_1',
        reaction: '👍',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/reaction`);

      await client.messages.sendPoll(dev, {
        to: '628123',
        question: 'Lunch?',
        options: ['Pizza', 'Sushi'],
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/poll`);

      await client.messages.sendLocation(dev, {
        to: '628123',
        latitude: -6.2,
        longitude: 106.8,
        name: 'Office HQ',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/location`);

      await client.messages.sendContact(dev, {
        to: '628123',
        display_name: 'John Doe',
        vcard: 'BEGIN:VCARD...',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/contact`);

      await client.messages.edit(dev, {
        to: '628123',
        target_message_id: 'msg_1',
        new_text: 'Edited text',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/edit`);

      await client.messages.revoke(dev, {
        to: '628123',
        target_message_id: 'msg_1',
      });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/messages/revoke`);

      await client.messages.markRead(dev, { chat_jid: '628123@s.whatsapp.net' });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/chat/read`);
    });

    it('groups: list, get, create, update, manageParticipants, getInviteCode, revokeInviteCode, leave, join', async () => {
      const dev = 'dev_grp_01';
      const jid = '1203630123@g.us';

      await client.groups.list(dev);
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/groups`);

      await client.groups.get(dev, jid);
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/metadata`);

      await client.groups.create(dev, { subject: 'Team', participants: ['628123'] });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/groups/create`);

      await client.groups.update(dev, jid, { subject: 'New Subject' });
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/metadata`);
      expect(lastMethod).toBe('PUT');

      await client.groups.manageParticipants(dev, jid, {
        action: 'add',
        participants: ['628999'],
      });
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/participants`);

      await client.groups.getInviteCode(dev, jid);
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/invite-link`);

      await client.groups.revokeInviteCode(dev, jid);
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/revoke-invite-link`);

      await client.groups.leave(dev, jid);
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/groups/${jid}/leave`);

      await client.groups.join(dev, { code: 'ABC123XYZ' });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/groups/join`);
    });

    it('contacts: check, getProfilePicture, getStatus, updateStatus', async () => {
      const dev = 'dev_cnt_01';
      const jid = '628123@s.whatsapp.net';

      await client.contacts.check(dev, { phones: ['628123456789'] });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/contacts/check`);

      await client.contacts.getProfilePicture(dev, jid);
      expect(decodeURIComponent(lastUrl)).toContain(`/api/v1/devices/${dev}/contacts/${jid}/avatar`);

      await client.contacts.getStatus(dev);
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/profile/status`);
      expect(lastMethod).toBe('GET');

      await client.contacts.updateStatus(dev, { status: 'Available for chats' });
      expect(lastUrl).toContain(`/api/v1/devices/${dev}/profile/status`);
      expect(lastMethod).toBe('PUT');
    });

    it('campaigns: list, create, get, delete, start, pause, getAnalytics', async () => {
      await client.campaigns.list({ status: 'SCHEDULED' });
      expect(lastUrl).toContain('/api/v1/campaigns?status=SCHEDULED');

      await client.campaigns.create({
        name: 'Flash Sale',
        channel_type: 'WHATSMEOW',
        message_template: 'Discount 50%',
        recipients: ['628123'],
        delay_min_seconds: 5,
        delay_max_seconds: 15,
      });
      expect(lastUrl).toContain('/api/v1/campaigns');
      expect(lastMethod).toBe('POST');

      await client.campaigns.get('cmp_01');
      expect(lastUrl).toContain('/api/v1/campaigns/cmp_01');

      await client.campaigns.start('cmp_01');
      expect(lastUrl).toContain('/api/v1/campaigns/cmp_01/start');

      await client.campaigns.pause('cmp_01');
      expect(lastUrl).toContain('/api/v1/campaigns/cmp_01/archive');

      await client.campaigns.delete('cmp_01');
      expect(lastUrl).toContain('/api/v1/campaigns/cmp_01');
      expect(lastMethod).toBe('DELETE');

      // Mock campaign analytics return
      const mockCampaignDetail = {
        id: 'cmp_01',
        total_recipients: 100,
        sent_count: 95,
        delivered_count: 90,
        failed_count: 5,
      };
      const detailMockFetch = vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ success: true, data: mockCampaignDetail }), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        })
      );
      const campaignClient = new ItssenderClient({ fetch: detailMockFetch });
      const analytics = await campaignClient.campaigns.getAnalytics('cmp_01');
      expect(analytics.total_recipients).toBe(100);
      expect(analytics.delivery_rate_percent).toBe(90);
      expect(analytics.failure_rate_percent).toBe(5);
    });

    it('analytics: getOverview, getMessages, getDevices', async () => {
      await client.analytics.getOverview({ from_timestamp: 1710000000000 });
      expect(lastUrl).toContain('/api/v1/analytics/usage?from_timestamp=1710000000000');

      await client.analytics.getMessages({ status: 'SENT', limit: 25 });
      expect(lastUrl).toContain('/api/v1/analytics/logs?limit=25&status=SENT');

      await client.analytics.getMessageDetail('log_999');
      expect(lastUrl).toContain('/api/v1/analytics/logs/log_999');

      await client.analytics.getDevices('dev_001');
      expect(lastUrl).toContain('/api/v1/analytics/usage?device_id=dev_001');
    });

    it('otp: send, verify, resend, get, list', async () => {
      await client.otp.send({ phone_number: '628123456789' });
      expect(lastUrl).toContain('/api/v1/otp/send');
      expect(lastMethod).toBe('POST');

      await client.otp.verify({ otp_id: 'otp_01', code: '123456' });
      expect(lastUrl).toContain('/api/v1/otp/verify');
      expect(lastMethod).toBe('POST');

      await client.otp.resend({ phone_number: '628123456789' });
      expect(lastUrl).toContain('/api/v1/otp/send');

      await client.otp.get('otp_01');
      expect(lastUrl).toContain('/api/v1/otp/otp_01');

      await client.otp.list({ status: 'VERIFIED' });
      expect(lastUrl).toContain('/api/v1/otp?status=VERIFIED');
    });

    it('templates: list, create, get, update, delete, render', async () => {
      await client.templates.list();
      expect(lastUrl).toContain('/api/v1/templates');

      await client.templates.create({
        name: 'Welcome OTP',
        category: 'AUTHENTICATION',
        content_raw: 'Your code is {{code}} for account {{username}}',
      });
      expect(lastUrl).toContain('/api/v1/templates');
      expect(lastMethod).toBe('POST');

      await client.templates.get('tpl_01');
      expect(lastUrl).toContain('/api/v1/templates/tpl_01');

      await client.templates.update('tpl_01', { name: 'Updated Welcome' });
      expect(lastUrl).toContain('/api/v1/templates/tpl_01');
      expect(lastMethod).toBe('PUT');

      await client.templates.delete('tpl_01');
      expect(lastUrl).toContain('/api/v1/templates/tpl_01');
      expect(lastMethod).toBe('DELETE');

      // Local render with raw template string
      const renderedStr = await client.templates.render(
        'Hello {{name}}, your balance is {{balance}} IDR.',
        { name: 'Alice', balance: 50000 }
      );
      expect(renderedStr).toBe('Hello Alice, your balance is 50000 IDR.');

      // Render with TemplateItem object
      const renderedObj = await client.templates.render(
        {
          id: 'tpl_123',
          user_id: 'usr_1',
          name: 'Order',
          category: 'TRANSACTION',
          content_raw: 'Order #{{order_id}} has been shipped to {{city}}.',
          variables: ['order_id', 'city'],
          is_archived: false,
          created_at: 1710000000000,
          updated_at: 1710000000000,
        },
        { order_id: 'ORD-999', city: 'Jakarta' }
      );
      expect(renderedObj).toBe('Order #ORD-999 has been shipped to Jakarta.');
    });

    it('webhooks: list, create, get, update, delete, test, listDeliveries, retryDelivery', async () => {
      await client.webhooks.list();
      expect(lastUrl).toContain('/api/v1/webhooks');

      await client.webhooks.create({
        url: 'https://myserver.com/webhook',
        events: ['message.received', 'otp.verified'],
      });
      expect(lastUrl).toContain('/api/v1/webhooks');
      expect(lastMethod).toBe('POST');

      await client.webhooks.update('wh_01', { is_active: false });
      expect(lastUrl).toContain('/api/v1/webhooks/wh_01');
      expect(lastMethod).toBe('PUT');

      await client.webhooks.delete('wh_01');
      expect(lastUrl).toContain('/api/v1/webhooks/wh_01');
      expect(lastMethod).toBe('DELETE');

      await client.webhooks.test('wh_01');
      expect(lastUrl).toContain('/api/v1/webhooks/wh_01/test');
      expect(lastMethod).toBe('POST');

      await client.webhooks.listDeliveries({ status: 'SUCCESS' });
      expect(lastUrl).toContain('/api/v1/webhooks/deliveries?status=SUCCESS');

      await client.webhooks.retryDelivery('whd_999');
      expect(lastUrl).toContain('/api/v1/webhooks/deliveries/whd_999/retry');
      expect(lastMethod).toBe('POST');
    });
  });
});
