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
export interface ApiErrorResponse extends ErrorResponse {
}
export declare class ApiError extends Error {
    readonly status: number;
    readonly payload?: ApiErrorResponse;
    constructor(status: number, payload?: ApiErrorResponse, message?: string);
}
export interface HubApiClientOptions {
    baseUrl?: string;
    fetch?: typeof fetch;
    defaultHeaders?: HeadersInit;
}
export declare class HubApiClient {
    private readonly baseUrl;
    private readonly fetchImpl;
    private readonly defaultHeaders;
    constructor(options?: HubApiClientOptions);
    health(): Promise<HealthStatus>;
    readiness(): Promise<HealthStatus>;
    listProjects(params?: ListProjectsParams): Promise<ProjectListResponse>;
    getProject(slug: string): Promise<Project>;
    listTags(): Promise<Tag[]>;
    private buildUrl;
    private buildQuery;
    private mergeHeaders;
    private get;
    private request;
    private parseBody;
}
export declare const createHubApiClient: (options?: HubApiClientOptions) => HubApiClient;
//# sourceMappingURL=index.d.ts.map