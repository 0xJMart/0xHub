import { http, HttpResponse } from 'msw';
import { findProjectBySlug, listProjects, mockTags } from './data';

const API_BASE = 'http://localhost:8080';

export const handlers = [
  http.get(`${API_BASE}/projects`, ({ request }) => {
    const url = new URL(request.url);
    const filters = {
      tag: url.searchParams.get('tag') ?? undefined,
      category: url.searchParams.get('category') ?? undefined,
      search: url.searchParams.get('search') ?? undefined,
    };

    const response = listProjects(filters);
    return HttpResponse.json(response);
  }),

  http.get(`${API_BASE}/projects/:slug`, ({ params }) => {
    const project = findProjectBySlug(String(params.slug));
    if (!project) {
      return HttpResponse.json({ error: 'Project not found' }, { status: 404 });
    }
    return HttpResponse.json(project);
  }),

  http.get(`${API_BASE}/tags`, () => HttpResponse.json(mockTags)),
];


