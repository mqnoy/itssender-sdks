import { BaseResource } from './base.js';
import type {
  ListOTPParams,
  OTPVerificationDetail,
  OTPVerificationListResponse,
  ResendOTPRequest,
  SendOTPRequest,
  SendOTPResponse,
  VerifyOTPRequest,
  VerifyOTPResponse,
} from '../types/otp.js';

export class OTPResource extends BaseResource {
  /**
   * Generates and dispatches a high-priority transactional WhatsApp OTP.
   */
  public async send(req: SendOTPRequest): Promise<SendOTPResponse> {
    return this.http.post<SendOTPResponse>('/api/v1/otp/send', req);
  }

  /**
   * Performs constant-time verification of recipient-provided OTP digits against an active session.
   */
  public async verify(req: VerifyOTPRequest): Promise<VerifyOTPResponse> {
    return this.http.post<VerifyOTPResponse>('/api/v1/otp/verify', req);
  }

  /**
   * Resends an OTP. If given an existing `otp_id`, fetches session details and re-dispatches;
   * otherwise dispatches a new OTP request.
   */
  public async resend(req: ResendOTPRequest | SendOTPRequest): Promise<SendOTPResponse> {
    if ('otp_id' in req && req.otp_id && !('phone_number' in req && req.phone_number)) {
      const existing = await this.get(req.otp_id);
      return this.send({
        phone_number: existing.phone_number,
        template_id: existing.template_id ?? undefined,
        device_id: existing.device_id ?? undefined,
        channel: existing.channel,
        code_length: existing.code_length,
      });
    }

    return this.send(req as SendOTPRequest);
  }

  /**
   * Retrieves verification state, attempts count, channel, and expiration for an OTP session.
   */
  public async get(id: string): Promise<OTPVerificationDetail> {
    return this.http.get<OTPVerificationDetail>(`/api/v1/otp/${encodeURIComponent(id)}`);
  }

  /**
   * Lists historical OTP verification sessions.
   */
  public async list(params?: ListOTPParams): Promise<OTPVerificationListResponse> {
    return this.http.get<OTPVerificationListResponse>('/api/v1/otp', {
      params: {
        status: params?.status,
        phone_number: params?.phone_number,
        limit: params?.limit,
        offset: params?.offset,
      },
    });
  }
}
