import { BaseResource } from './base.js';
import type { ActionSuccessResponse, DeleteResponse } from '../types/common.js';
import type {
  ChatPresenceRequest,
  CreateDeviceRequest,
  CreateDeviceResponse,
  Device,
  DeviceInstanceUsage,
  DeviceLogsParams,
  DeviceLogsResponse,
  DevicePresenceRequest,
  PairCodeRequest,
  PairCodeResponse,
  QRCodeResponse,
  UpdateDeviceSettingsRequest,
} from '../types/devices.js';

export class DevicesResource extends BaseResource {
  /**
   * Lists all WhatsApp device instances for the authenticated tenant.
   */
  public async list(): Promise<Device[]> {
    return this.http.get<Device[]>('/api/v1/devices');
  }

  /**
   * Retrieves detailed configuration and connection status for a specific device instance.
   */
  public async get(id: string): Promise<Device> {
    return this.http.get<Device>(`/api/v1/devices/${encodeURIComponent(id)}`);
  }

  /**
   * Provisions a new WhatsApp device worker container instance.
   */
  public async create(req: CreateDeviceRequest): Promise<CreateDeviceResponse> {
    return this.http.post<CreateDeviceResponse>('/api/v1/devices', req);
  }

  /**
   * Destroys and removes a device instance and its backing container.
   */
  public async delete(id: string): Promise<DeleteResponse> {
    return this.http.delete<DeleteResponse>(`/api/v1/devices/${encodeURIComponent(id)}`);
  }

  /**
   * Requests a pairing QR code for WhatsApp web client connection.
   */
  public async getQr(id: string): Promise<QRCodeResponse> {
    return this.http.get<QRCodeResponse>(`/api/v1/devices/${encodeURIComponent(id)}/qr`);
  }

  /**
   * Requests an 8-digit WhatsApp pairing code for linking via phone number.
   */
  public async pairCode(id: string, req: PairCodeRequest): Promise<PairCodeResponse> {
    return this.http.post<PairCodeResponse>(
      `/api/v1/devices/${encodeURIComponent(id)}/pair-code`,
      req
    );
  }

  /**
   * Restarts the WhatsApp device instance container worker.
   */
  public async restart(id: string): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(id)}/restart`
    );
  }

  /**
   * Fetches container logs for the device instance worker.
   */
  public async getLogs(id: string, params?: DeviceLogsParams): Promise<DeviceLogsResponse> {
    return this.http.get<DeviceLogsResponse>(`/api/v1/devices/${encodeURIComponent(id)}/logs`, {
      params: {
        tail: params?.tail,
        timestamps: params?.timestamps,
      },
    });
  }

  /**
   * Updates device settings, proxy, or resource limits.
   */
  public async updateSettings(id: string, req: UpdateDeviceSettingsRequest): Promise<Device> {
    return this.http.put<Device>(`/api/v1/devices/${encodeURIComponent(id)}`, req);
  }

  /**
   * Sets presence status. Automatically routes to device presence or chat presence
   * depending on whether `to` / `recipient` or chat states ('composing' | 'recording' | 'paused') are provided.
   */
  public async setPresence(
    id: string,
    req: DevicePresenceRequest | ChatPresenceRequest
  ): Promise<ActionSuccessResponse> {
    if (
      'to' in req ||
      'recipient' in req ||
      req.state === 'composing' ||
      req.state === 'recording' ||
      req.state === 'paused'
    ) {
      return this.setChatPresence(id, req as ChatPresenceRequest);
    }
    return this.setDevicePresence(id, req as DevicePresenceRequest);
  }

  /**
   * Sets global device presence ('available' | 'unavailable').
   */
  public async setDevicePresence(
    id: string,
    req: DevicePresenceRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(id)}/presence/device`,
      req
    );
  }

  /**
   * Sets real-time chat typing/recording presence ('composing' | 'recording' | 'paused').
   */
  public async setChatPresence(
    id: string,
    req: ChatPresenceRequest
  ): Promise<ActionSuccessResponse> {
    return this.http.post<ActionSuccessResponse>(
      `/api/v1/devices/${encodeURIComponent(id)}/presence/chat`,
      req
    );
  }

  /**
   * Gets runtime operational state and message telemetry for the device.
   */
  public async getState(id: string): Promise<DeviceInstanceUsage> {
    return this.http.get<DeviceInstanceUsage>(
      `/api/v1/devices/${encodeURIComponent(id)}/usage`
    );
  }
}
