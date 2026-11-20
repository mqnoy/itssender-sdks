import { BaseResource } from './base.js';
import type {
  AuthResponse,
  LoginRequest,
  PublicConfigResponse,
  RegisterRequest,
  UserResponse,
} from '../types/auth.js';

export class AuthResource extends BaseResource {
  /**
   * Registers a new tenant user account.
   */
  public async register(req: RegisterRequest): Promise<AuthResponse> {
    const res = await this.http.post<AuthResponse>('/api/v1/auth/register', req);
    if (res?.token) {
      this.http.setBearerToken(res.token);
    }
    return res;
  }

  /**
   * Authenticates user with email and password, returning JWT bearer token and user profile.
   */
  public async login(req: LoginRequest): Promise<AuthResponse> {
    const res = await this.http.post<AuthResponse>('/api/v1/auth/login', req);
    if (res?.token) {
      this.http.setBearerToken(res.token);
    }
    return res;
  }

  /**
   * Fetches currently authenticated user identity and role.
   */
  public async me(): Promise<UserResponse> {
    return this.http.get<UserResponse>('/api/v1/auth/me');
  }

  /**
   * Refreshes JWT authentication token.
   */
  public async refreshToken(refreshToken?: string): Promise<AuthResponse> {
    const body = refreshToken ? { refresh_token: refreshToken } : undefined;
    const res = await this.http.post<AuthResponse>('/api/v1/auth/refresh', body);
    if (res?.token) {
      this.http.setBearerToken(res.token);
    }
    return res;
  }

  /**
   * Retrieves public platform authentication and registration configuration.
   */
  public async getConfig(): Promise<PublicConfigResponse> {
    return this.http.get<PublicConfigResponse>('/api/v1/auth/config');
  }
}
