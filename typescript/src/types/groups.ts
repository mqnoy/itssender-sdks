import type { UnixMsTimestamp } from './common.js';

export interface GroupParticipant {
  jid: string;
  is_admin?: boolean;
  is_super_admin?: boolean;
}

export interface GroupSummary {
  jid: string;
  name?: string;
  subject?: string;
  topic?: string;
  owner_jid?: string;
  participant_count?: number;
  participants?: GroupParticipant[] | string[];
}

export interface GroupMetadata {
  jid: string;
  owner_jid?: string;
  subject: string;
  topic?: string;
  participants?: GroupParticipant[];
  creation?: UnixMsTimestamp;
  [key: string]: unknown;
}

export interface CreateGroupRequest {
  subject: string;
  participants: string[];
}

export interface GroupCreateResponse {
  success: boolean;
  group_jid: string;
  owner_jid?: string;
  timestamp: UnixMsTimestamp;
}

export type GroupParticipantAction = 'add' | 'remove' | 'promote' | 'demote';

export interface UpdateGroupParticipantsRequest {
  action: GroupParticipantAction;
  participants: string[];
}

export interface UpdateGroupMetadataRequest {
  subject?: string;
  topic?: string;
}

export interface GroupInviteLinkResponse {
  success: boolean;
  invite_code?: string;
  invite_url?: string;
  timestamp: UnixMsTimestamp;
}

export interface JoinGroupRequest {
  code: string;
}

export interface JoinGroupResponse {
  success: boolean;
  group_jid?: string;
  timestamp: UnixMsTimestamp;
  [key: string]: unknown;
}
