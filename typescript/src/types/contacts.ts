import type { UnixMsTimestamp } from './common.js';

export interface CheckContactsRequest {
  phones: string[];
}

export interface ContactCheckResult {
  phone: string;
  is_registered: boolean;
  jid?: string;
}

export interface CheckContactsResponse {
  success: boolean;
  results: ContactCheckResult[];
  timestamp: UnixMsTimestamp;
}

export interface ContactAvatarResponse {
  success: boolean;
  url?: string;
  id?: string;
  timestamp: UnixMsTimestamp;
}

export interface UpdateStatusRequest {
  status: string;
}
