import type { PaginationParams, SafeInteger, UnixMsTimestamp } from './common.js';

export type OTPStatus =
  | 'PENDING'
  | 'DISPATCHED'
  | 'VERIFIED'
  | 'FAILED'
  | 'EXPIRED'
  | 'MAX_ATTEMPTS_EXCEEDED'
  | string;

export interface SendOTPRequest {
  phone_number: string;
  template_id?: string;
  custom_code?: string;
  code_length?: SafeInteger;
  ttl_ms?: SafeInteger;
  max_attempts?: SafeInteger;
  device_id?: string;
  channel?: 'WHATSMEOW' | 'META_CLOUD' | 'AUTO' | string;
  metadata?: Record<string, unknown>;
}

export interface SendOTPResponse {
  otp_id: string;
  phone_number: string;
  status: OTPStatus;
  expires_at: UnixMsTimestamp;
  max_attempts: SafeInteger;
  created_at: UnixMsTimestamp;
}

export interface VerifyOTPRequest {
  otp_id: string;
  code: string;
  phone_number?: string;
}

export interface VerifyOTPResponse {
  success: boolean;
  otp_id: string;
  status: OTPStatus;
  verified_at?: UnixMsTimestamp | null;
  remaining_attempts?: SafeInteger | null;
}

export interface OTPVerificationDetail {
  id: string;
  user_id: string;
  template_id?: string | null;
  device_id?: string | null;
  phone_number: string;
  code_length: SafeInteger;
  status: OTPStatus;
  attempts: SafeInteger;
  max_attempts: SafeInteger;
  channel: string;
  message_log_id?: string | null;
  metadata?: Record<string, unknown>;
  expires_at: UnixMsTimestamp;
  verified_at?: UnixMsTimestamp | null;
  created_at: UnixMsTimestamp;
}

export interface OTPVerificationListResponse {
  sessions: OTPVerificationDetail[];
  total: SafeInteger;
  data?: OTPVerificationDetail[];
}

export interface ListOTPParams extends PaginationParams {
  status?: OTPStatus;
  phone_number?: string;
}

export interface ResendOTPRequest {
  otp_id?: string;
  phone_number?: string;
  template_id?: string;
  device_id?: string;
  channel?: string;
}
