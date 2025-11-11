export class ApiError extends Error {
    status;
    payload;
    constructor(status, payload, message) {
        super(message ?? payload?.error ?? `Request failed with status ${status}`);
        this.name = 'ApiError';
        this.status = status;
        this.payload = payload;
    }
}
const DEFAULT_BASE_URL = 'http://localhost:8080';
export class HubApiClient {
    baseUrl;
    fetchImpl;
    defaultHeaders;
    constructor(options = {}) {
        const { baseUrl, fetch: fetchImpl, defaultHeaders } = options;
        this.baseUrl = (baseUrl ?? DEFAULT_BASE_URL).replace(/\/+$/, '');
        this.fetchImpl = fetchImpl ?? globalThis.fetch;
        if (!this.fetchImpl) {
            throw new Error('A fetch implementation must be provided when running outside of environments with a global fetch.');
        }
        this.defaultHeaders = defaultHeaders ?? { Accept: 'application/json' };
    }
    async health() {
        return this.get('/healthz');
    }
    async readiness() {
        return this.get('/readyz');
    }
    async listProjects(params) {
        const query = this.buildQuery(params);
        return this.get(`/projects${query}`);
    }
    async getProject(slug) {
        if (!slug) {
            throw new Error('slug is required');
        }
        const encoded = encodeURIComponent(slug);
        return this.get(`/projects/${encoded}`);
    }
    async listTags() {
        return this.get('/tags');
    }
    buildUrl(path) {
        if (!path.startsWith('/')) {
            path = `/${path}`;
        }
        return `${this.baseUrl}${path}`;
    }
    buildQuery(params) {
        if (!params) {
            return '';
        }
        const searchParams = new URLSearchParams();
        const entries = [
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
    mergeHeaders(headers) {
        const merged = new Headers(this.defaultHeaders);
        if (headers) {
            const additional = new Headers(headers);
            additional.forEach((value, key) => merged.set(key, value));
        }
        return merged;
    }
    async get(path) {
        return this.request(path, { method: 'GET' });
    }
    async request(path, init) {
        const response = await this.fetchImpl(this.buildUrl(path), {
            ...init,
            headers: this.mergeHeaders(init.headers),
        });
        const contentType = response.headers.get('content-type') ?? '';
        const isJson = contentType.includes('application/json');
        const rawBody = await response.text();
        const parsedBody = this.parseBody(rawBody, isJson);
        if (!response.ok) {
            const payload = isJson && parsedBody && typeof parsedBody === 'object'
                ? parsedBody
                : undefined;
            throw new ApiError(response.status, payload, payload?.error ?? response.statusText);
        }
        return (parsedBody ?? undefined);
    }
    parseBody(body, isJson) {
        if (!body) {
            return undefined;
        }
        if (isJson) {
            try {
                return JSON.parse(body);
            }
            catch {
                return undefined;
            }
        }
        return body;
    }
}
export const createHubApiClient = (options) => new HubApiClient(options);
//# sourceMappingURL=index.js.map