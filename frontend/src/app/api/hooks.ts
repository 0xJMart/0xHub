import type { ListProjectsParams, Project, ProjectListResponse, Tag } from '@0xhub/api';
import { useQuery, useSuspenseQuery } from '@tanstack/react-query';
import { hubApiClient } from './client';

const queryNamespace = 'hub' as const;

export const projectKeys = {
  all: [queryNamespace, 'projects'] as const,
  list: (params: ListProjectsParams | undefined) =>
    [queryNamespace, 'projects', 'list', params ?? {}] as const,
  detail: (slug: string) => [queryNamespace, 'projects', 'detail', slug] as const,
  tags: () => [queryNamespace, 'tags'] as const,
};

export const useProjectsQuery = (params: ListProjectsParams | undefined) =>
  useQuery({
    queryKey: projectKeys.list(params),
    queryFn: async (): Promise<ProjectListResponse> => hubApiClient.listProjects(params),
    keepPreviousData: true,
  });

export const useSuspenseProjectsQuery = (params: ListProjectsParams | undefined) =>
  useSuspenseQuery({
    queryKey: projectKeys.list(params),
    queryFn: async (): Promise<ProjectListResponse> => hubApiClient.listProjects(params),
  });

export const useProjectQuery = (slug: string) =>
  useQuery({
    queryKey: projectKeys.detail(slug),
    queryFn: async (): Promise<Project> => hubApiClient.getProject(slug),
    enabled: Boolean(slug),
  });

export const useTagsQuery = () =>
  useQuery({
    queryKey: projectKeys.tags(),
    queryFn: async (): Promise<Tag[]> => hubApiClient.listTags(),
    staleTime: 5 * 60_000,
  });


