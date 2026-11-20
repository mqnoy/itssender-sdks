import { BaseResource } from './base.js';
import type { ActionSuccessResponse } from '../types/common.js';
import type {
  Campaign,
  CampaignAnalytics,
  CampaignDetail,
  CreateCampaignRequest,
  ListCampaignsParams,
} from '../types/campaigns.js';

export class CampaignsResource extends BaseResource {
  /**
   * Lists scheduled, running, or archived broadcast campaigns.
   */
  public async list(params?: ListCampaignsParams): Promise<Campaign[]> {
    return this.http.get<Campaign[]>('/api/v1/campaigns', {
      params: {
        status: params?.status,
        archived: params?.archived,
        limit: params?.limit,
        offset: params?.offset,
      },
    });
  }

  /**
   * Creates and schedules a bulk broadcast campaign.
   */
  public async create(req: CreateCampaignRequest): Promise<Campaign> {
    return this.http.post<Campaign>('/api/v1/campaigns', req);
  }

  /**
   * Retrieves campaign configuration and recipient delivery status.
   */
  public async get(id: string): Promise<CampaignDetail> {
    return this.http.get<CampaignDetail>(`/api/v1/campaigns/${encodeURIComponent(id)}`);
  }

  /**
   * Deletes / archives a campaign.
   */
  public async delete(id: string): Promise<ActionSuccessResponse> {
    return this.http.delete<ActionSuccessResponse>(
      `/api/v1/campaigns/${encodeURIComponent(id)}`
    );
  }

  /**
   * Starts or resumes campaign execution.
   */
  public async start(id: string): Promise<Campaign> {
    return this.http.post<Campaign>(`/api/v1/campaigns/${encodeURIComponent(id)}/start`);
  }

  /**
   * Pauses an active broadcast campaign.
   */
  public async pause(id: string): Promise<Campaign> {
    return this.http.post<Campaign>(`/api/v1/campaigns/${encodeURIComponent(id)}/archive`);
  }

  /**
   * Archives a non-running campaign.
   */
  public async archive(id: string): Promise<Campaign> {
    return this.http.post<Campaign>(`/api/v1/campaigns/${encodeURIComponent(id)}/archive`);
  }

  /**
   * Restores an archived campaign to active listing.
   */
  public async unarchive(id: string): Promise<Campaign> {
    return this.http.post<Campaign>(`/api/v1/campaigns/${encodeURIComponent(id)}/unarchive`);
  }

  /**
   * Calculates delivery and performance analytics for a campaign.
   */
  public async getAnalytics(id: string): Promise<CampaignAnalytics> {
    const detail = await this.get(id);
    const total = detail.total_recipients || 0;
    const sent = detail.sent_count || 0;
    const delivered = detail.delivered_count || 0;
    const failed = detail.failed_count || 0;

    const deliveryRate = total > 0 ? Math.round((delivered / total) * 10000) / 100 : 0;
    const failureRate = total > 0 ? Math.round((failed / total) * 10000) / 100 : 0;

    return {
      campaign_id: detail.id,
      total_recipients: total,
      sent_count: sent,
      delivered_count: delivered,
      failed_count: failed,
      delivery_rate_percent: deliveryRate,
      failure_rate_percent: failureRate,
    };
  }
}
