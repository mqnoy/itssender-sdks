import type { PaginationParams, SafeInteger, UnixMsTimestamp } from './common.js';

export type CampaignStatus =
  | 'DRAFT'
  | 'SCHEDULED'
  | 'RUNNING'
  | 'COMPLETED'
  | 'PAUSED'
  | 'FAILED'
  | string;

export interface CampaignItem {
  id: string;
  campaign_id: string;
  recipient: string;
  message_type: string;
  status: 'PENDING' | 'SENT' | 'DELIVERED' | 'FAILED' | string;
  error_message?: string | null;
  sent_at?: UnixMsTimestamp | null;
}

export interface Campaign {
  id: string;
  user_id: string;
  name: string;
  channel_type: 'WHATSMEOW' | 'META_CLOUD' | string;
  device_id?: string | null;
  status: CampaignStatus;
  total_recipients: SafeInteger;
  sent_count: SafeInteger;
  delivered_count: SafeInteger;
  failed_count: SafeInteger;
  delay_min_seconds: SafeInteger;
  delay_max_seconds: SafeInteger;
  is_archived: boolean;
  scheduled_at?: UnixMsTimestamp | null;
  created_at: UnixMsTimestamp;
  updated_at: UnixMsTimestamp;
}

export interface CampaignDetail extends Campaign {
  items?: CampaignItem[];
}

export interface CampaignAnalytics {
  campaign_id: string;
  total_recipients: SafeInteger;
  sent_count: SafeInteger;
  delivered_count: SafeInteger;
  failed_count: SafeInteger;
  delivery_rate_percent?: number;
  failure_rate_percent?: number;
}

export interface CreateCampaignRequest {
  name: string;
  channel_type: 'WHATSMEOW' | 'META_CLOUD';
  channel?: string;
  device_id?: string;
  meta_token?: string;
  meta_api_key?: string;
  meta_phone_id?: string;
  template_name?: string;
  message_template: string;
  message?: string;
  recipients: string[];
  delay_min_seconds: SafeInteger;
  delay_max_seconds: SafeInteger;
  message_type?: 'TEXT' | 'MEDIA';
  media_type?: 'image' | 'video' | 'audio' | 'document' | 'sticker';
  media_url?: string;
  media_data_base64?: string;
  media_data?: string;
  media_caption?: string;
  media_filename?: string;
  media_mime_type?: string;
  scheduled_at?: UnixMsTimestamp;
}

export interface ListCampaignsParams extends PaginationParams {
  status?: CampaignStatus;
  archived?: 'true' | 'false' | 'all';
}
