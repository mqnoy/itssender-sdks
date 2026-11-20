import { BaseResource } from './base.js';
import type { ActionSuccessResponse } from '../types/common.js';
import type {
  TemplateCreateRequest,
  TemplateItem,
  TemplateListResponse,
  TemplateUpdateRequest,
} from '../types/templates.js';

export class TemplatesResource extends BaseResource {
  /**
   * Lists all message templates for the authenticated tenant.
   */
  public async list(): Promise<TemplateListResponse> {
    return this.http.get<TemplateListResponse>('/api/v1/templates');
  }

  /**
   * Creates a new multi-tenant rich-text WhatsApp message template.
   */
  public async create(req: TemplateCreateRequest): Promise<TemplateItem> {
    return this.http.post<TemplateItem>('/api/v1/templates', req);
  }

  /**
   * Retrieves a message template by ID.
   */
  public async get(id: string): Promise<TemplateItem> {
    return this.http.get<TemplateItem>(`/api/v1/templates/${encodeURIComponent(id)}`);
  }

  /**
   * Updates an existing message template.
   */
  public async update(id: string, req: TemplateUpdateRequest): Promise<TemplateItem> {
    return this.http.put<TemplateItem>(`/api/v1/templates/${encodeURIComponent(id)}`, req);
  }

  /**
   * Soft-archives / deletes a message template.
   */
  public async delete(id: string): Promise<ActionSuccessResponse> {
    return this.http.delete<ActionSuccessResponse>(
      `/api/v1/templates/${encodeURIComponent(id)}`
    );
  }

  /**
   * Renders a message template by substituting variable placeholders (`{{variable_name}}`)
   * with provided key-value mappings.
   *
   * @param templateOrId - A TemplateItem, a template ID string (e.g. 'tpl_...'), or raw template text with `{{var}}` placeholders.
   * @param variables - Map of variable names to their string or number replacements.
   * @returns The fully rendered text string.
   */
  public async render(
    templateOrId: string | TemplateItem,
    variables: Record<string, string | number> = {}
  ): Promise<string> {
    let rawContent: string;

    if (typeof templateOrId === 'object' && templateOrId !== null) {
      rawContent = templateOrId.content_raw || '';
    } else if (typeof templateOrId === 'string') {
      if (templateOrId.includes('{{') || templateOrId.includes(' ')) {
        // Raw template string
        rawContent = templateOrId;
      } else {
        // Template ID - fetch template
        const item = await this.get(templateOrId);
        rawContent = item.content_raw || '';
      }
    } else {
      rawContent = '';
    }

    return rawContent.replace(/\{\{\s*([a-zA-Z0-9_-]+)\s*\}\}/g, (match, key) => {
      if (Object.prototype.hasOwnProperty.call(variables, key)) {
        return String(variables[key]);
      }
      return match;
    });
  }
}
