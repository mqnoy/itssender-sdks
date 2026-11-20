import { BaseResource } from './base.js';
import type { SuccessResponse } from '../types/common.js';
import type {
  CreateWebhookRequest,
  ListWebhookDeliveriesParams,
  UpdateWebhookRequest,
  WebhookConfigResponse,
  WebhookDelivery,
  WebhookDeliveryLog,
  WebhookDeliveryListResponse,
} from '../types/webhooks.js';

export class WebhooksResource extends BaseResource {
  /**
   * Lists all active and configured webhook endpoints for the tenant.
   */
  public async list(): Promise<WebhookConfigResponse[]> {
    return this.http.get<WebhookConfigResponse[]>('/api/v1/webhooks');
  }

  /**
   * Registers a new webhook endpoint with Standard Webhooks HMAC-SHA256 signing secret.
   */
  public async create(req: CreateWebhookRequest): Promise<WebhookConfigResponse> {
    return this.http.post<WebhookConfigResponse>('/api/v1/webhooks', req);
  }

  /**
   * Retrieves a specific webhook configuration by ID.
   */
  public async get(id: string): Promise<WebhookConfigResponse> {
    try {
      return await this.http.get<WebhookConfigResponse>(
        `/api/v1/webhooks/${encodeURIComponent(id)}`
      );
    } catch {
      // Fallback: list all and find by ID if direct route is not mounted
      const all = await this.list();
      const found = all.find((w) => w.id === id);
      if (found) {
        return found;
      }
      throw new Error(`Webhook configuration with ID ${id} not found`);
    }
  }

  /**
   * Updates target URL, subscribed topics, active status, or PII masking preferences.
   */
  public async update(
    id: string,
    req: UpdateWebhookRequest
  ): Promise<WebhookConfigResponse> {
    return this.http.put<WebhookConfigResponse>(
      `/api/v1/webhooks/${encodeURIComponent(id)}`,
      req
    );
  }

  /**
   * Soft-archives / deactivates a webhook configuration.
   */
  public async delete(id: string): Promise<SuccessResponse> {
    return this.http.delete<SuccessResponse>(
      `/api/v1/webhooks/${encodeURIComponent(id)}`
    );
  }

  /**
   * Dispatches a synthetic signed test ping event (`webhook.test`) to verify connectivity.
   */
  public async test(id: string): Promise<WebhookDelivery> {
    return this.http.post<WebhookDelivery>(
      `/api/v1/webhooks/${encodeURIComponent(id)}/test`
    );
  }

  /**
   * Lists webhook delivery audit logs with statuses, HTTP response codes, and latencies.
   */
  public async listDeliveries(
    params?: ListWebhookDeliveriesParams
  ): Promise<WebhookDeliveryLog[]> {
    const res = await this.http.get<WebhookDeliveryLog[] | WebhookDeliveryListResponse>(
      '/api/v1/webhooks/deliveries',
      {
        params: {
          webhook_id: params?.webhook_id,
          event: params?.event,
          status: params?.status,
          limit: params?.limit,
          offset: params?.offset,
        },
      }
    );

    if (Array.isArray(res)) {
      return res;
    }
    if (res && 'deliveries' in res && Array.isArray(res.deliveries)) {
      return res.deliveries;
    }
    if (res && 'data' in res && Array.isArray(res.data)) {
      return res.data;
    }
    return [];
  }

  /**
   * Retries an outbound webhook event delivery.
   */
  public async retryDelivery(deliveryId: string): Promise<WebhookDelivery> {
    return this.http.post<WebhookDelivery>(
      `/api/v1/webhooks/deliveries/${encodeURIComponent(deliveryId)}/retry`
    );
  }
}
