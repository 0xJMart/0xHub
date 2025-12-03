import { vi } from 'vitest';
import type { ListProjectsParams, Project } from '@0xhub/api';
import { findProjectBySlug, listProjects, mockTags } from '@/mocks/data';

export const listProjectsMock = vi.fn(
  async (params?: ListProjectsParams) => listProjects(params ?? {}),
);

export const getProjectMock = vi.fn(async (slug: string): Promise<Project> => {
  const project = findProjectBySlug(slug);
  if (!project) {
    throw new Error('Project not found');
  }
  return project;
});

export const listTagsMock = vi.fn(async () => mockTags);

vi.mock('@/app/api/client', () => ({
  hubApiClient: {
    listProjects: (...args: Parameters<typeof listProjectsMock>) => listProjectsMock(...args),
    getProject: (...args: Parameters<typeof getProjectMock>) => getProjectMock(...args),
    listTags: (...args: Parameters<typeof listTagsMock>) => listTagsMock(...args),
  },
}));


