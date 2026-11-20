import type { PaginationParams, SafeInteger, UnixMsTimestamp } from './common.js';

export type WebhookEventTopic =
  | '*'
  | 'device.*'
  | 'device.created'
  | 'device.updated'
  | 'device.deleted'
  | 'device.connected'
  | 'device.disconnected'
  | 'device.qr_generated'
  | 'device.pair_code_generated'
  | 'device.logged_out'
  | 'device.battery_changed'
  | 'message.*'
  | 'message.received'
  | 'message.sent'
  | 'message.delivered'
  | 'message.read'
  | 'message.failed'
  | 'message.revoked'
  | 'group.*'
  | 'group.created'
  | 'group.participants_updated'
  | 'group.left'
  | 'group.metadata_updated'
  | 'campaign.*'
  | 'campaign.created'
  | 'campaign.started'
  | 'campaign.completed'
  | 'campaign.paused'
  | 'otp.*'
  | 'otp.sent'
  | 'otp.verified'
  | 'otp.failed'
  | 'otp.expired'
  | 'webhook.test'
  | string;

export interface CreateWebhookRequest {
  url: string;
  secret?: string;
  events: WebhookEventTopic[];
  is_active?: boolean;
  mask_pii?: boolean;
}

export interface UpdateWebhookRequest {
  url?: string;
  secret?: string;
  events?: WebhookEventTopic[];
  is_active?: boolean;
  mask_pii?: boolean;
}

export interface WebhookConfigResponse {
  id: string;
  user_id: string;
  url: string;
  secret: string;
  events: WebhookEventTopic[];
  is_active: boolean;
  is_archived: boolean;
  mask_pii: boolean;
  created_at: UnixMsTimestamp;
  updated_at: UnixMsTimestamp;
}

export interface WebhookDeliveryLog {
  id: string;
  webhook_id: string;
  user_id: string;
  event: string;
  payload: Record<string, unknown>;
  signature: string;
  http_status?: SafeInteger | null;
  response_body?: string | null;
  duration_ms: SafeInteger;
  attempt_count: SafeInteger;
  max_retries: SafeInteger;
  next_retry_at_ms?: UnixMsTimestamp | null;
  status: 'SUCCESS' | 'FAILED' | 'RETRYING' | 'PENDING' | 'PENDING_RETRY' | string;
  error_message?: string | null;
  created_at: UnixMsTimestamp;
}

export type WebhookDelivery = WebhookDeliveryLog;

export interface WebhookDeliveryListResponse {
  deliveries: WebhookDeliveryLog[];
  total: SafeInteger;
  data?: WebhookDeliveryLog[];
}

export interface ListWebhookDeliveriesParams extends PaginationParams {
  webhook_id?: string;
  event?: WebhookEventTopic;
  status?: 'SUCCESS' | 'FAILED' | 'RETRYING' | 'PENDING' | 'PENDING_RETRY' | string;
}

export interface WebhookEnvelope<T = unknown> {
  event: WebhookEventTopic;
  delivery_id: string;
  timestamp: UnixMsTimestamp;
  data: T;
}
