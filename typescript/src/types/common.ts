/**
 * Common types and API response envelopes for itssender SDK.
 *
 * Invariant 1: Pure integer precision for currency and quotas (never floats).
 * Invariant 2: Hybrid Unix timestamps in milliseconds (e.g. 1710000000000).
 */

/**
 * Unix timestamp in milliseconds (Invariant 2).
 * Example: 1710000000000
 */
export type UnixMsTimestamp = number;

/**
 * Safe integer representation for quotas and monetary cents (Invariant 1).
 */
export type SafeInteger = number;

/**
 * Standard API error descriptor.
 */
export interface APIError {
  code: string;
  message: string;
}

/**
 * Standard API response envelope returned by the itssender gateway.
 */
export interface APIResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: APIError;
  timestamp?: UnixMsTimestamp;
}

/**
 * Generic success acknowledgement response.
 */
export interface SuccessResponse {
  success: boolean;
  message?: string;
  timestamp?: UnixMsTimestamp;
}

/**
 * Generic action outcome response.
 */
export interface ActionSuccessResponse {
  success: boolean;
  timestamp?: UnixMsTimestamp;
  [key: string]: unknown;
}

/**
 * Resource deletion confirmation.
 */
export interface DeleteResponse {
  deleted: boolean;
  id?: string;
  timestamp?: UnixMsTimestamp;
}

/**
 * Standard pagination query parameters.
 */
export interface PaginationParams {
  limit?: SafeInteger;
  offset?: SafeInteger;
}
