export interface HealthStatus {
  status: string;
  timestamp: string;
}

export interface ErrorResponse {
  error: string;
  detail?: string;
  correlationId?: string;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  sortOrder?: number;
}

export interface Tag {
  id: string;
  name: string;
  displayName: string;
  color?: string;
}

export interface ProjectLink {
  id: string;
  linkType: string;
  url: string;
  label?: string;
}

export type MediaStatus = 'pending' | 'available' | 'errored';

export interface MediaAsset {
  id: string;
  projectId?: string | null;
  storagePath: string;
  originalFilename: string;
  mimeType?: string;
  sizeBytes?: number;
  description?: string;
  status: MediaStatus;
  createdAt?: string;
  updatedAt?: string;
  expiresAt?: string | null;
}

export type ProjectStatus = 'draft' | 'active' | 'archived' | 'deprecated';

export interface ProjectSummary {
  id: string;
  title: string;
  slug: string;
  summary?: string;
  status: ProjectStatus;
  category?: Category;
  tags: Tag[];
  heroMedia?: MediaAsset;
}

export interface Project extends ProjectSummary {
  description?: string;
  links: ProjectLink[];
  media: MediaAsset[];
}

export interface ProjectListResponse {
  items: ProjectSummary[];
  total: number;
}

export interface ListProjectsParams {
  tag?: string;
  category?: string;
  search?: string;
  limit?: number;
  offset?: number;
}

export interface ApiErrorResponse extends ErrorResponse {}

export class ApiError extends Error {
  readonly status: number;
  readonly payload?: ApiErrorResponse;

  constructor(status: number, payload?: ApiErrorResponse, message?: string) {
    super(message ?? payload?.error ?? `Request failed with status ${status}`);
    this.name = 'ApiError';
    this.status = status;
    this.payload = payload;
  }
}

export interface HubApiClientOptions {
  baseUrl?: string;
  fetch?: typeof fetch;
  defaultHeaders?: HeadersInit;
}

const DEFAULT_BASE_URL = 'http://localhost:8080';

export class HubApiClient {
  private readonly baseUrl: string;
  private readonly fetchImpl: typeof fetch;
  private readonly defaultHeaders: HeadersInit;

  constructor(options: HubApiClientOptions = {}) {
    const { baseUrl, fetch: fetchImpl, defaultHeaders } = options;
    this.baseUrl = (baseUrl ?? DEFAULT_BASE_URL).replace(/\/+$/, '');
    this.fetchImpl = fetchImpl ?? globalThis.fetch;
    if (!this.fetchImpl) {
      throw new Error(
        'A fetch implementation must be provided when running outside of environments with a global fetch.',
      );
    }
    this.defaultHeaders = defaultHeaders ?? { Accept: 'application/json' };
  }

  async health(): Promise<HealthStatus> {
    return this.get<HealthStatus>('/healthz');
  }

  async readiness(): Promise<HealthStatus> {
    return this.get<HealthStatus>('/readyz');
  }

  async listProjects(params?: ListProjectsParams): Promise<ProjectListResponse> {
    const query = this.buildQuery(params);
    return this.get<ProjectListResponse>(`/projects${query}`);
  }

  async getProject(slug: string): Promise<Project> {
    if (!slug) {
      throw new Error('slug is required');
    }
    const encoded = encodeURIComponent(slug);
    return this.get<Project>(`/projects/${encoded}`);
  }

  async listTags(): Promise<Tag[]> {
    return this.get<Tag[]>('/tags');
  }

  private buildUrl(path: string): string {
    if (!path.startsWith('/')) {
      path = `/${path}`;
    }
    return `${this.baseUrl}${path}`;
  }

  private buildQuery(params?: ListProjectsParams): string {
    if (!params) {
      return '';
    }
    const searchParams = new URLSearchParams();
    const entries: Array<[key: string, value: string | number | undefined]> = [
      ['tag', params.tag],
      ['category', params.category],
      ['search', params.search],
      ['limit', params.limit],
      ['offset', params.offset],
    ];
    for (const [key, value] of entries) {
      if (value === undefined || value === null || value === '') {
        continue;
      }
      searchParams.set(key, String(value));
    }
    const serialized = searchParams.toString();
    return serialized ? `?${serialized}` : '';
  }

  private mergeHeaders(headers?: HeadersInit): Headers {
    const merged = new Headers(this.defaultHeaders);
    if (headers) {
      const additional = new Headers(headers);
      additional.forEach((value, key) => merged.set(key, value));
    }
    return merged;
  }

  private async get<T>(path: string): Promise<T> {
    return this.request<T>(path, { method: 'GET' });
  }

  private async request<T>(path: string, init: RequestInit): Promise<T> {
    const response = await this.fetchImpl(this.buildUrl(path), {
      ...init,
      headers: this.mergeHeaders(init.headers),
    });

    const contentType = response.headers.get('content-type') ?? '';
    const isJson = contentType.includes('application/json');
    const rawBody = await response.text();
    const parsedBody = this.parseBody(rawBody, isJson);

    if (!response.ok) {
      const payload =
        isJson && parsedBody && typeof parsedBody === 'object'
          ? (parsedBody as ApiErrorResponse)
          : undefined;
      throw new ApiError(response.status, payload, payload?.error ?? response.statusText);
    }

    return (parsedBody ?? (undefined as unknown)) as T;
  }

  private parseBody(body: string, isJson: boolean): unknown {
    if (!body) {
      return undefined;
    }
    if (isJson) {
      try {
        return JSON.parse(body);
      } catch {
        return undefined;
      }
    }
    return body;
  }
}

export const createHubApiClient = (options?: HubApiClientOptions): HubApiClient =>
  new HubApiClient(options);

