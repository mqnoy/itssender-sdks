import type { UnixMsTimestamp } from './common.js';

export interface SendTextMessageRequest {
  to: string;
  recipient?: string;
  message?: string;
  text?: string;
  reply_to_id?: string;
  quoted_message_id?: string;
  quoted_participant?: string;
  mentions?: string[];
}

export type MediaType = 'image' | 'document' | 'video' | 'audio' | 'sticker';

export interface SendMediaMessageRequest {
  to: string;
  recipient?: string;
  media_type: MediaType;
  data?: string;
  data_base64?: string;
  media_base64?: string;
  url?: string;
  media_url?: string;
  caption?: string;
  filename?: string;
  mime_type?: string;
  reply_to_id?: string;
}

export interface SendReactionRequest {
  to: string;
  recipient?: string;
  message_id?: string;
  target_message_id?: string;
  reaction?: string;
  emoji?: string;
}

export interface SendPollRequest {
  to: string;
  recipient?: string;
  question: string;
  options: string[];
  max_answers?: number;
  max_selectable?: number;
}

export interface SendLocationMessageRequest {
  to: string;
  recipient?: string;
  latitude: number;
  longitude: number;
  name?: string;
  address?: string;
}

export interface SendContactMessageRequest {
  to: string;
  recipient?: string;
  display_name: string;
  vcard: string;
}

export interface EditMessageRequest {
  to: string;
  recipient?: string;
  target_message_id: string;
  new_text: string;
}

export interface RevokeMessageRequest {
  to: string;
  recipient?: string;
  target_message_id: string;
}

export interface MarkReadRequest {
  chat_jid: string;
  message_ids?: string[];
}

export interface SendMessageResponse {
  message_id: string;
  status: 'PENDING' | 'SENT' | 'DELIVERED' | 'FAILED' | string;
  timestamp: UnixMsTimestamp;
  [key: string]: unknown;
}
