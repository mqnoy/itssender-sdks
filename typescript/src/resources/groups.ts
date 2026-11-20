import { BaseResource } from './base.js';
import type { ActionSuccessResponse } from '../types/common.js';
import type {
  CreateGroupRequest,
  GroupCreateResponse,
  GroupInviteLinkResponse,
  GroupMetadata,
  GroupSummary,
  JoinGroupRequest,
  JoinGroupResponse,
  UpdateGroupMetadataRequest,
  UpdateGroupParticipantsRequest,
} from '../types/groups.js';

export class GroupsResource extends BaseResource {
  /**
   * Lists WhatsApp groups the device is currently a member of.
   */
  public async list(deviceId: string): Promise<GroupSummary[]> {
    return this.http.get<GroupSummary[]>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups`
    );
  }

  /**
   * Retrieves metadata and participant details for a specific WhatsApp group.
   */
  public async get(deviceId: string, jid: string): Promise<GroupMetadata> {
    return this.http.get<GroupMetadata>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/metadata`
    );
  }

  /**
   * Creates a new WhatsApp group chat.
   */
  public async create(
    deviceId: string,
    req: CreateGroupRequest
  ): Promise<GroupCreateResponse> {
    return this.http.post<GroupCreateResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/create`,
      req
    );
  }

  /**
   * Updates group metadata (subject title and/or topic description).
   */
  public async update(
    deviceId: string,
    jid: string,
    req: UpdateGroupMetadataRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.put<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/metadata`,
      req
    );
  }

  /**
   * Adds, removes, promotes, or demotes participants in a WhatsApp group.
   */
  public async manageParticipants(
    deviceId: string,
    jid: string,
    req: UpdateGroupParticipantsRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/participants`,
      req
    );
  }

  /**
   * Retrieves or generates a shareable WhatsApp group invite link.
   */
  public async getInviteCode(
    deviceId: string,
    jid: string
  ): Promise<GroupInviteLinkResponse> {
    return this.http.get<GroupInviteLinkResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/invite-link`
    );
  }

  /**
   * Revokes and invalidates the existing WhatsApp group invite link.
   */
  public async revokeInviteCode(
    deviceId: string,
    jid: string
  ): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/revoke-invite-link`
    );
  }

  /**
   * Exits and leaves the specified WhatsApp group.
   */
  public async leave(deviceId: string, jid: string): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/${encodeURIComponent(jid)}/leave`
    );
  }

  /**
   * Joins a WhatsApp group using an invite link code.
   */
  public async join(
    deviceId: string,
    req: JoinGroupRequest
  ): Promise<JoinGroupResponse> {
    return this.http.post<JoinGroupResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/groups/join`,
      req
    );
  }
}
