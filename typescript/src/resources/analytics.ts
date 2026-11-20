import { BaseResource } from './base.js';
import type {
  AnalyticsLogsParams,
  AnalyticsUsageParams,
  DeviceUsageAnalytics,
  MessageAuditLogDetail,
  MessageAuditLogsResponse,
  UsageAnalytics,
} from '../types/analytics.js';

export class AnalyticsResource extends BaseResource {
  /**
   * Queries message accounting logs and aggregations.
   */
  public async getOverview(params?: AnalyticsUsageParams): Promise<UsageAnalytics> {
    return this.http.get<UsageAnalytics>('/api/v1/analytics/usage', {
      params: {
        device_id: params?.device_id,
        from_timestamp: params?.from_timestamp,
        to_timestamp: params?.to_timestamp,
      },
    });
  }

  /**
   * Queries message audit transmission logs with direction, status, and search filters.
   */
  public async getMessages(params?: AnalyticsLogsParams): Promise<MessageAuditLogsResponse> {
    return this.http.get<MessageAuditLogsResponse>('/api/v1/analytics/logs', {
      params: {
        limit: params?.limit,
        offset: params?.offset,
        status: params?.status,
        direction: params?.direction,
        api_key_id: params?.api_key_id,
        q: params?.q,
      },
    });
  }

  /**
   * Retrieves full audit log detail for a specific message log ID.
   */
  public async getMessageDetail(id: string): Promise<MessageAuditLogDetail> {
    return this.http.get<MessageAuditLogDetail>(
      `/api/v1/analytics/logs/${encodeURIComponent(id)}`
    );
  }

  /**
   * Retrieves message usage statistics broken down by device, or for a specific device ID.
   */
  public async getDevices(deviceId?: string): Promise<DeviceUsageAnalytics | UsageAnalytics> {
    if (deviceId) {
      const usage = await this.getOverview({ device_id: deviceId });
      return {
        device_id: deviceId,
        usage,
      };
    }
    return this.getOverview();
  }
}
