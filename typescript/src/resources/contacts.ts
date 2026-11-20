import { BaseResource } from './base.js';
import type { ActionSuccessResponse } from '../types/common.js';
import type {
  CheckContactsRequest,
  CheckContactsResponse,
  ContactAvatarResponse,
  UpdateStatusRequest,
} from '../types/contacts.js';

export class ContactsResource extends BaseResource {
  /**
   * Verifies whether phone numbers are registered WhatsApp accounts.
   */
  public async check(
    deviceId: string,
    req: CheckContactsRequest
  ): Promise<CheckContactsResponse> {
    return this.http.post<CheckContactsResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/contacts/check`,
      req
    );
  }

  /**
   * Fetches the profile picture URL for a contact or group JID.
   */
  public async getProfilePicture(
    deviceId: string,
    jid: string
  ): Promise<ContactAvatarResponse> {
    return this.http.get<ContactAvatarResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/contacts/${encodeURIComponent(jid)}/avatar`
    );
  }

  /**
   * Retrieves the current device profile status bio.
   */
  public async getStatus(deviceId: string): Promise<UpdateStatusRequest> {
    return this.http.get<UpdateStatusRequest>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/profile/status`
    );
  }

  /**
   * Updates the device WhatsApp profile "About / Status" bio.
   */
  public async updateStatus(
    deviceId: string,
    req: UpdateStatusRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.put<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/profile/status`,
      req
    );
  }
}
