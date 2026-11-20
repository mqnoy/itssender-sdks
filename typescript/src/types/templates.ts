import type { SafeInteger, UnixMsTimestamp } from './common.js';

export interface TemplateItem {
  id: string;
  user_id: string;
  name: string;
  category: string;
  description?: string | null;
  content_raw: string;
  content_html?: string;
  variables: string[];
  is_archived: boolean;
  archived_at?: UnixMsTimestamp | null;
  created_at: UnixMsTimestamp;
  updated_at: UnixMsTimestamp;
}

export interface TemplateCreateRequest {
  name: string;
  category: string;
  description?: string;
  content_raw: string;
  content_html?: string;
  variables?: string[];
}

export interface TemplateUpdateRequest {
  name?: string;
  category?: string;
  description?: string;
  content_raw?: string;
  content_html?: string;
  variables?: string[];
}

export interface TemplateListResponse {
  templates: TemplateItem[];
  total: SafeInteger;
  data?: TemplateItem[];
}

export interface RenderTemplateParams {
  variables: Record<string, string | number>;
}
