import { BaseResource } from './base.js';
import type { ApiKey, ApiKeyCreated, CreateApiKeyRequest } from '../types/api-keys.js';
import type { DeleteResponse } from '../types/common.js';

export class ApiKeysResource extends BaseResource {
  /**
   * Lists all active and registered API keys for the authenticated tenant.
   */
  public async list(): Promise<ApiKey[]> {
    return this.http.get<ApiKey[]>('/api/v1/user/keys');
  }

  /**
   * Generates a new API key for external integrations.
   * Note: The plaintext `raw_key` is only returned once upon creation.
   */
  public async create(req: CreateApiKeyRequest): Promise<ApiKeyCreated> {
    return this.http.post<ApiKeyCreated>('/api/v1/user/keys', req);
  }

  /**
   * Revokes and removes an existing API key.
   */
  public async delete(id: string): Promise<DeleteResponse> {
    return this.http.delete<DeleteResponse>(`/api/v1/user/keys/${encodeURIComponent(id)}`);
  }
}
