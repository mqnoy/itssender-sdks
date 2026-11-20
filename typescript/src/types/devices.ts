import type { SafeInteger, UnixMsTimestamp } from './common.js';

export type DeviceStatus =
  | 'PROVISIONING'
  | 'WAITING_QR'
  | 'CONNECTING'
  | 'CONNECTED'
  | 'DISCONNECTED'
  | 'DEGRADED'
  | 'STOPPED'
  | string;

export interface DeviceInstanceUsage {
  sent_count?: SafeInteger;
  delivered_count?: SafeInteger;
  read_count?: SafeInteger;
  failed_count?: SafeInteger;
  total_messages?: SafeInteger;
  last_activity_at?: UnixMsTimestamp;
}

export interface Device {
  id: string;
  user_id: string;
  name: string;
  phone_number?: string | null;
  status: DeviceStatus;
  container_id?: string | null;
  internal_ip?: string | null;
  internal_token?: string;
  proxy_url?: string | null;
  cpu_limit?: number;
  nano_cpus?: SafeInteger;
  mem_limit_mb?: SafeInteger;
  memory_mb?: SafeInteger;
  instance_usage?: DeviceInstanceUsage;
  created_at: UnixMsTimestamp;
  updated_at: UnixMsTimestamp;
}

export interface CreateDeviceRequest {
  name: string;
  phone_number?: string;
  proxy_url?: string;
  cpu_limit?: number;
  nano_cpus?: SafeInteger;
  mem_limit_mb?: SafeInteger;
  memory_mb?: SafeInteger;
}

export interface CreateDeviceResponse {
  id: string;
  name: string;
  status: DeviceStatus;
  container_id?: string | null;
  ip_address?: string | null;
  created_at: UnixMsTimestamp;
}

export interface PairCodeRequest {
  phone_number: string;
}

export interface PairCodeResponse {
  success: boolean;
  code?: string;
  pairing_code?: string;
  expires_in?: SafeInteger;
  [key: string]: unknown;
}

export interface QRCodeResponse {
  success?: boolean;
  qr?: string;
  qr_code?: string;
  status?: string;
  [key: string]: unknown;
}

export interface DeviceLogsParams {
  tail?: SafeInteger;
  timestamps?: boolean;
}

export interface DeviceLogsResponse {
  device_id: string;
  container_id?: string;
  tail: SafeInteger;
  timestamps: boolean;
  stdout: string;
  stderr: string;
  logs: string;
}

export interface UpdateDeviceSettingsRequest {
  name?: string;
  proxy_url?: string;
  cpu_limit?: number;
  nano_cpus?: SafeInteger;
  mem_limit_mb?: SafeInteger;
  memory_mb?: SafeInteger;
}

export type DevicePresenceState = 'available' | 'unavailable';

export interface DevicePresenceRequest {
  state: DevicePresenceState;
}

export type ChatPresenceState = 'composing' | 'recording' | 'paused';

export interface ChatPresenceRequest {
  to?: string;
  recipient?: string;
  state: ChatPresenceState;
}
