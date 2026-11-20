/**
 * Error hierarchy for the itssender TypeScript SDK.
 */

export class ItssenderError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ItssenderError';
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderAPIError extends ItssenderError {
  readonly statusCode: number;
  readonly code: string;
  readonly response?: unknown;
  readonly headers?: Headers;

  constructor(
    statusCode: number,
    code: string,
    message: string,
    response?: unknown,
    headers?: Headers
  ) {
    super(message || `API request failed with status ${statusCode} (${code})`);
    this.name = 'ItssenderAPIError';
    this.statusCode = statusCode;
    this.code = code;
    this.response = response;
    this.headers = headers;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderAuthenticationError extends ItssenderAPIError {
  constructor(
    statusCode: number,
    code: string,
    message: string,
    response?: unknown,
    headers?: Headers
  ) {
    super(statusCode, code, message, response, headers);
    this.name = 'ItssenderAuthenticationError';
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderNotFoundError extends ItssenderAPIError {
  constructor(
    code: string,
    message: string,
    response?: unknown,
    headers?: Headers
  ) {
    super(404, code, message, response, headers);
    this.name = 'ItssenderNotFoundError';
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderValidationError extends ItssenderAPIError {
  constructor(
    statusCode: number,
    code: string,
    message: string,
    response?: unknown,
    headers?: Headers
  ) {
    super(statusCode, code, message, response, headers);
    this.name = 'ItssenderValidationError';
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderRateLimitError extends ItssenderAPIError {
  readonly retryAfterSeconds?: number;

  constructor(
    code: string,
    message: string,
    retryAfterSeconds?: number,
    response?: unknown,
    headers?: Headers
  ) {
    super(429, code, message, response, headers);
    this.name = 'ItssenderRateLimitError';
    this.retryAfterSeconds = retryAfterSeconds;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderTimeoutError extends ItssenderError {
  readonly timeoutMs: number;

  constructor(timeoutMs: number, message?: string) {
    super(message || `Request timed out after ${timeoutMs}ms`);
    this.name = 'ItssenderTimeoutError';
    this.timeoutMs = timeoutMs;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class ItssenderNetworkError extends ItssenderError {
  readonly cause?: unknown;

  constructor(message: string, cause?: unknown) {
    super(message);
    this.name = 'ItssenderNetworkError';
    this.cause = cause;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

export class WebhookVerificationError extends ItssenderError {
  readonly reason:
    | 'MISSING_HEADERS'
    | 'INVALID_SIGNATURE'
    | 'TIMESTAMP_OUT_OF_TOLERANCE'
    | 'INVALID_PAYLOAD'
    | 'INVALID_SECRET'
    | string;

  constructor(
    message: string,
    reason:
      | 'MISSING_HEADERS'
      | 'INVALID_SIGNATURE'
      | 'TIMESTAMP_OUT_OF_TOLERANCE'
      | 'INVALID_PAYLOAD'
      | 'INVALID_SECRET'
      | string
  ) {
    super(message);
    this.name = 'WebhookVerificationError';
    this.reason = reason;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}
