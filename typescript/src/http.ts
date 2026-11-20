import {
  ItssenderAPIError,
  ItssenderAuthenticationError,
  ItssenderNetworkError,
  ItssenderNotFoundError,
  ItssenderRateLimitError,
  ItssenderTimeoutError,
  ItssenderValidationError,
} from './errors.js';
import type { APIResponse } from './types/common.js';

export interface HttpClientConfig {
  /**
   * Base URL for the itssender Gateway API.
   * Defaults to 'https://itssender.imzalab.com'.
   */
  baseURL?: string;

  /**
   * API Key for authentication (sent via `X-API-Key` header).
   */
  apiKey?: string;

  /**
   * Bearer JWT token for session authentication (sent via `Authorization: Bearer <token>`).
   */
  bearerToken?: string;

  /**
   * Request timeout in milliseconds.
   * Defaults to 30,000 ms (30 seconds).
   */
  timeout?: number;

  /**
   * Maximum number of retry attempts for network errors and 429/5xx responses.
   * Defaults to 3. Set to 0 to disable automatic retries.
   */
  maxRetries?: number;

  /**
   * Initial base backoff delay in milliseconds.
   * Defaults to 500 ms.
   */
  baseDelayMs?: number;

  /**
   * Maximum backoff delay cap in milliseconds.
   * Defaults to 10,000 ms.
   */
  maxDelayMs?: number;

  /**
   * Custom fetch implementation (optional, defaults to global fetch).
   */
  fetch?: typeof fetch;

  /**
   * Default headers included in every request.
   */
  defaultHeaders?: Record<string, string>;
}

export interface RequestOptions {
  headers?: Record<string, string | undefined>;
  params?: Record<string, string | number | boolean | undefined | null>;
  timeout?: number;
  skipRetry?: boolean;
}

export class HttpClient {
  private baseURL: string;
  private apiKey?: string;
  private bearerToken?: string;
  private timeout: number;
  private maxRetries: number;
  private baseDelayMs: number;
  private maxDelayMs: number;
  private customFetch: typeof fetch;
  private defaultHeaders: Record<string, string>;

  constructor(config: HttpClientConfig = {}) {
    this.baseURL = (config.baseURL || 'https://itssender.imzalab.com').replace(/\/+$/, '');
    this.apiKey = config.apiKey;
    this.bearerToken = config.bearerToken;
    this.timeout = config.timeout ?? 30000;
    this.maxRetries = config.maxRetries ?? 3;
    this.baseDelayMs = config.baseDelayMs ?? 500;
    this.maxDelayMs = config.maxDelayMs ?? 10000;
    this.customFetch = config.fetch || globalThis.fetch;
    this.defaultHeaders = config.defaultHeaders || {};
  }

  public setApiKey(apiKey: string | undefined): void {
    this.apiKey = apiKey;
  }

  public setBearerToken(token: string | undefined): void {
    this.bearerToken = token;
  }

  public getBaseURL(): string {
    return this.baseURL;
  }

  public async get<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>('GET', path, undefined, options);
  }

  public async post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>('POST', path, body, options);
  }

  public async put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>('PUT', path, body, options);
  }

  public async delete<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>('DELETE', path, undefined, options);
  }

  public async patch<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>('PATCH', path, body, options);
  }

  public async request<T>(
    method: string,
    path: string,
    body?: unknown,
    options: RequestOptions = {}
  ): Promise<T> {
    const url = this.buildUrl(path, options.params);
    const headers = this.buildHeaders(options.headers, body);
    const timeoutMs = options.timeout ?? this.timeout;
    const maxRetries = options.skipRetry ? 0 : this.maxRetries;

    let attempt = 0;
    let lastError: unknown;

    while (attempt <= maxRetries) {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), timeoutMs);

      try {
        const fetchOptions: RequestInit = {
          method,
          headers,
          signal: controller.signal,
        };

        if (body !== undefined) {
          fetchOptions.body = typeof body === 'string' ? body : JSON.stringify(body);
        }

        const response = await this.customFetch(url, fetchOptions);
        clearTimeout(timer);

        if (!response.ok) {
          // Check retryable status codes
          if (attempt < maxRetries && this.isRetryableStatus(response.status)) {
            const delay = this.calculateDelay(attempt, response.headers);
            await this.sleep(delay);
            attempt++;
            continue;
          }

          await this.handleHttpError(response);
        }

        // Process successful response
        return await this.parseResponse<T>(response);
      } catch (err: unknown) {
        clearTimeout(timer);

        // Check if aborted due to timeout
        if (controller.signal.aborted) {
          throw new ItssenderTimeoutError(
            timeoutMs,
            `Request to ${method} ${path} timed out after ${timeoutMs}ms`
          );
        }

        // If it's already an ItssenderError, rethrow directly unless retryable
        if (err instanceof ItssenderAPIError) {
          throw err;
        }

        // Handle network/connection error with retry
        if (attempt < maxRetries && this.isNetworkError(err)) {
          const delay = this.calculateDelay(attempt);
          await this.sleep(delay);
          attempt++;
          lastError = err;
          continue;
        }

        if (err instanceof Error && err.name === 'AbortError') {
          throw new ItssenderTimeoutError(timeoutMs);
        }

        throw new ItssenderNetworkError(
          `Network request to ${method} ${path} failed: ${(err as Error)?.message || String(err)}`,
          err
        );
      }
    }

    throw lastError instanceof Error
      ? lastError
      : new ItssenderNetworkError(`Request to ${method} ${path} failed after ${maxRetries} retries`, lastError);
  }

  private buildUrl(
    path: string,
    params?: Record<string, string | number | boolean | undefined | null>
  ): string {
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    const url = new URL(`${this.baseURL}${cleanPath}`);

    if (params) {
      for (const [key, val] of Object.entries(params)) {
        if (val !== undefined && val !== null) {
          url.searchParams.set(key, String(val));
        }
      }
    }

    return url.toString();
  }

  private buildHeaders(
    customHeaders?: Record<string, string | undefined>,
    body?: unknown
  ): Record<string, string> {
    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...this.defaultHeaders,
    };

    if (this.apiKey) {
      headers['X-API-Key'] = this.apiKey;
    }

    if (this.bearerToken) {
      headers['Authorization'] = `Bearer ${this.bearerToken}`;
    }

    if (body !== undefined && typeof body !== 'string') {
      headers['Content-Type'] = 'application/json';
    } else if (typeof body === 'string' && !headers['Content-Type']) {
      headers['Content-Type'] = 'application/json';
    }

    if (customHeaders) {
      for (const [key, val] of Object.entries(customHeaders)) {
        if (val !== undefined) {
          headers[key] = val;
        }
      }
    }

    return headers;
  }

  private async parseResponse<T>(response: Response): Promise<T> {
    const text = await response.text();
    if (!text.trim()) {
      return {} as T;
    }

    let parsed: unknown;
    try {
      parsed = JSON.parse(text);
    } catch {
      return text as unknown as T;
    }

    // Check envelope
    if (parsed && typeof parsed === 'object' && 'success' in parsed) {
      const envelope = parsed as APIResponse<T>;
      if (envelope.success === false) {
        throw new ItssenderAPIError(
          response.status,
          envelope.error?.code || 'UNKNOWN_ERROR',
          envelope.error?.message || 'Operation failed',
          parsed,
          response.headers
        );
      }

      // If data is provided in envelope, unwrap it; otherwise return envelope
      if ('data' in envelope && envelope.data !== undefined) {
        return envelope.data as T;
      }
    }

    return parsed as T;
  }

  private async handleHttpError(response: Response): Promise<never> {
    const status = response.status;
    let code = 'HTTP_ERROR';
    let message = response.statusText || `Request failed with status ${status}`;
    let body: unknown;

    try {
      const text = await response.text();
      if (text) {
        try {
          body = JSON.parse(text);
          if (body && typeof body === 'object') {
            const b = body as Record<string, unknown>;
            if (b.error && typeof b.error === 'object') {
              const errObj = b.error as Record<string, unknown>;
              code = (errObj.code as string) || code;
              message = (errObj.message as string) || message;
            } else if (typeof b.message === 'string') {
              message = b.message;
            }
          }
        } catch {
          message = text;
        }
      }
    } catch {
      // Ignore body read errors
    }

    if (status === 401 || status === 403) {
      throw new ItssenderAuthenticationError(status, code, message, body, response.headers);
    }

    if (status === 404) {
      throw new ItssenderNotFoundError(code, message, body, response.headers);
    }

    if (status === 400 || status === 422) {
      throw new ItssenderValidationError(status, code, message, body, response.headers);
    }

    if (status === 429) {
      let retryAfter: number | undefined;
      const retryHeader = response.headers.get('Retry-After');
      if (retryHeader) {
        const parsed = parseInt(retryHeader, 10);
        if (!isNaN(parsed)) {
          retryAfter = parsed;
        }
      }
      throw new ItssenderRateLimitError(code, message, retryAfter, body, response.headers);
    }

    throw new ItssenderAPIError(status, code, message, body, response.headers);
  }

  private isRetryableStatus(status: number): boolean {
    return status === 429 || status === 500 || status === 502 || status === 503 || status === 504;
  }

  private isNetworkError(err: unknown): boolean {
    if (err instanceof TypeError) {
      return true; // Typical for fetch failed in Node/browser
    }
    if (err && typeof err === 'object' && 'code' in err) {
      const code = (err as { code: string }).code;
      return (
        code === 'ECONNRESET' ||
        code === 'ETIMEDOUT' ||
        code === 'ECONNREFUSED' ||
        code === 'EAI_AGAIN'
      );
    }
    return false;
  }

  private calculateDelay(attempt: number, headers?: Headers): number {
    if (headers) {
      const retryAfter = headers.get('Retry-After');
      if (retryAfter) {
        const parsed = parseInt(retryAfter, 10);
        if (!isNaN(parsed) && parsed > 0) {
          return parsed * 1000;
        }
      }
    }

    // Exponential backoff with full jitter:
    // delay = min(maxDelay, baseDelay * 2^attempt) + random_jitter
    const expDelay = Math.min(this.maxDelayMs, this.baseDelayMs * Math.pow(2, attempt));
    const jitter = Math.floor(Math.random() * (this.baseDelayMs / 2));
    return expDelay + jitter;
  }

  private sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}
