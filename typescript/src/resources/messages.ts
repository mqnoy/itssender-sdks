import { BaseResource } from './base.js';
import type { ActionSuccessResponse } from '../types/common.js';
import type {
  EditMessageRequest,
  MarkReadRequest,
  RevokeMessageRequest,
  SendContactMessageRequest,
  SendLocationMessageRequest,
  SendMediaMessageRequest,
  SendMessageResponse,
  SendPollRequest,
  SendReactionRequest,
  SendTextMessageRequest,
} from '../types/messages.js';

export class MessagesResource extends BaseResource {
  /**
   * Dispatches a plaintext WhatsApp message to a recipient phone number or group JID.
   */
  public async sendText(
    deviceId: string,
    req: SendTextMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/text`,
      req
    );
  }

  /**
   * Dispatches a media attachment (image, document, video, audio, sticker) to WhatsApp.
   */
  public async sendMedia(
    deviceId: string,
    req: SendMediaMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/media`,
      req
    );
  }

  /**
   * Sends an emoji reaction to a previously sent or received WhatsApp message.
   */
  public async sendReaction(
    deviceId: string,
    req: SendReactionRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/reaction`,
      req
    );
  }

  /**
   * Sends an interactive multi-choice poll message to WhatsApp.
   */
  public async sendPoll(
    deviceId: string,
    req: SendPollRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/poll`,
      req
    );
  }

  /**
   * Sends a geographic location pin message with optional venue title and address.
   */
  public async sendLocation(
    deviceId: string,
    req: SendLocationMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/location`,
      req
    );
  }

  /**
   * Sends a contact business card (vCard) message to WhatsApp.
   */
  public async sendContact(
    deviceId: string,
    req: SendContactMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/contact`,
      req
    );
  }

  /**
   * Edits the text content of a previously dispatched WhatsApp message.
   */
  public async edit(
    deviceId: string,
    req: EditMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/edit`,
      req
    );
  }

  /**
   * Revokes (deletes for everyone) a previously sent WhatsApp message.
   */
  public async revoke(
    deviceId: string,
    req: RevokeMessageRequest
  ): Promise<SendMessageResponse> {
    return this.http.post<SendMessageResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/messages/revoke`,
      req
    );
  }

  /**
   * Marks incoming messages in a conversation as read.
   */
  public async markRead(
    deviceId: string,
    req: MarkReadRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(deviceId)}/chat/read`,
      req
    );
  }
}
