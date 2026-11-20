import type { UnixMsTimestamp } from './common.js';

export interface CreateApiKeyRequest {
  name: string;
}

export interface ApiKey {
  id: string;
  user_id: string;
  name: string;
  is_active: boolean;
  created_at: UnixMsTimestamp;
}

export interface ApiKeyCreated {
  id: string;
  name: string;
  raw_key: string;
  created_at: UnixMsTimestamp;
}
