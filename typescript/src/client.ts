import { HttpClient, type HttpClientConfig } from './http.js';
import { AnalyticsResource } from './resources/analytics.js';
import { ApiKeysResource } from './resources/api-keys.js';
import { AuthResource } from './resources/auth.js';
import { CampaignsResource } from './resources/campaigns.js';
import { ContactsResource } from './resources/contacts.js';
import { DevicesResource } from './resources/devices.js';
import { GroupsResource } from './resources/groups.js';
import { MessagesResource } from './resources/messages.js';
import { OTPResource } from './resources/otp.js';
import { TemplatesResource } from './resources/templates.js';
import { WebhooksResource } from './resources/webhooks.js';

export const DEFAULT_BASE_URL = 'https://itssender.imzalab.com';

export interface ItssenderClientOptions extends HttpClientConfig {
  /**
   * Shorthand token: if starts with 'its_' or 'its_live_' or 'its_test_', treated as apiKey;
   * otherwise treated as bearerToken unless specified explicitly.
   */
  token?: string;
}

export class ItssenderClient {
  private readonly http: HttpClient;

  public readonly auth: AuthResource;
  public readonly apiKeys: ApiKeysResource;
  public readonly devices: DevicesResource;
  public readonly messages: MessagesResource;
  public readonly groups: GroupsResource;
  public readonly contacts: ContactsResource;
  public readonly campaigns: CampaignsResource;
  public readonly analytics: AnalyticsResource;
  public readonly otp: OTPResource;
  public readonly templates: TemplatesResource;
  public readonly webhooks: WebhooksResource;

  constructor(options: ItssenderClientOptions = {}) {
    let apiKey = options.apiKey;
    let bearerToken = options.bearerToken;

    if (options.token) {
      if (
        options.token.startsWith('its_') ||
        options.token.startsWith('key_') ||
        options.token.includes('live') ||
        options.token.includes('test')
      ) {
        apiKey = apiKey || options.token;
      } else {
        bearerToken = bearerToken || options.token;
      }
    }

    this.http = new HttpClient({
      baseURL: options.baseURL ?? DEFAULT_BASE_URL,
      apiKey,
      bearerToken,
      timeout: options.timeout,
      maxRetries: options.maxRetries,
      baseDelayMs: options.baseDelayMs,
      maxDelayMs: options.maxDelayMs,
      fetch: options.fetch,
      defaultHeaders: options.defaultHeaders,
    });

    this.auth = new AuthResource(this.http);
    this.apiKeys = new ApiKeysResource(this.http);
    this.devices = new DevicesResource(this.http);
    this.messages = new MessagesResource(this.http);
    this.groups = new GroupsResource(this.http);
    this.contacts = new ContactsResource(this.http);
    this.campaigns = new CampaignsResource(this.http);
    this.analytics = new AnalyticsResource(this.http);
    this.otp = new OTPResource(this.http);
    this.templates = new TemplatesResource(this.http);
    this.webhooks = new WebhooksResource(this.http);
  }

  /**
   * Updates the API key used for external authentication (`X-API-Key`).
   */
  public setApiKey(apiKey: string | undefined): void {
    this.http.setApiKey(apiKey);
  }

  /**
   * Updates the Bearer JWT token used for session authentication (`Authorization: Bearer <token>`).
   */
  public setBearerToken(token: string | undefined): void {
    this.http.setBearerToken(token);
  }

  /**
   * Returns the current base URL for the gateway API.
   */
  public getBaseURL(): string {
    return this.http.getBaseURL();
  }
}
