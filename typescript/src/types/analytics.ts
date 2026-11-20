import type { PaginationParams, SafeInteger, UnixMsTimestamp } from './common.js';

export interface DeviceUsage {
  device_id: string;
  device_name?: string;
  sent_count: SafeInteger;
  delivered_count: SafeInteger;
  failed_count: SafeInteger;
  total_messages: SafeInteger;
}

export interface UsageAnalytics {
  total_sent: SafeInteger;
  total_delivered: SafeInteger;
  total_failed: SafeInteger;
  total_messages: SafeInteger;
  device_breakdown?: DeviceUsage[];
  period_start?: UnixMsTimestamp;
  period_end?: UnixMsTimestamp;
}

export interface AnalyticsUsageParams {
  device_id?: string;
  from_timestamp?: UnixMsTimestamp;
  to_timestamp?: UnixMsTimestamp;
}

export interface MessageAuditLog {
  id: string;
  device_id: string;
  device_name?: string;
  api_key_id?: string | null;
  api_key_name?: string | null;
  direction: 'OUTBOUND' | 'INBOUND' | string;
  channel: string;
  recipient: string;
  sender?: string;
  message_type: string;
  status: 'PENDING' | 'SENT' | 'DELIVERED' | 'READ' | 'FAILED' | string;
  error_message?: string | null;
  timestamp: UnixMsTimestamp;
  created_at: UnixMsTimestamp;
}

export interface MessageAuditLogsResponse {
  success: boolean;
  data: MessageAuditLog[];
  total?: SafeInteger;
  limit?: SafeInteger;
  offset?: SafeInteger;
}

export interface MessageAuditLogDetail {
  log: MessageAuditLog;
  recipient_logs?: unknown[];
  device_logs?: unknown[];
  total_recipient_logs?: SafeInteger;
  total_device_logs?: SafeInteger;
}

export interface AnalyticsLogsParams extends PaginationParams {
  status?: 'SENT' | 'DELIVERED' | 'READ' | 'FAILED' | string;
  direction?: 'OUTBOUND' | 'INBOUND' | string;
  api_key_id?: string;
  q?: string;
}

export interface DeviceUsageAnalytics {
  device_id: string;
  usage: DeviceUsage | UsageAnalytics;
}
