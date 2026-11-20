import type { UnixMsTimestamp } from './common.js';

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  id: string;
  email: string;
  role: 'admin' | 'user' | string;
  created_at: UnixMsTimestamp;
}

export interface UserResponse {
  id: string;
  email: string;
  role: 'admin' | 'user' | string;
  is_suspended?: boolean;
  created_at: UnixMsTimestamp;
  updated_at?: UnixMsTimestamp;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface PublicConfigResponse {
  registration_enabled: boolean;
  meta_cloud_enabled?: boolean;
  version?: string;
  [key: string]: unknown;
}
