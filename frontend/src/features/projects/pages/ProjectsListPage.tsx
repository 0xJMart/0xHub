import { useEffect, useMemo } from 'react';
import { PageHeader, ProjectCard } from '@0xhub/ui';
import { useNavigate } from 'react-router-dom';
import { useProjectsQuery, useTagsQuery } from '@/app/api/hooks';
import { useToast } from '@/app/providers/ToastProvider';
import { ErrorState } from '@/components/ErrorState';
import { ProjectsFilters } from '@/features/projects/components/ProjectsFilters';
import {
  normalizeFilters,
  type ProjectsFiltersState,
  useProjectsSearchParams,
} from '@/features/projects/utils/useProjectsSearchParams';

export const ProjectsListPage = (): JSX.Element => {
  const { filters, setFilters, listParams } = useProjectsSearchParams();
  const { toast } = useToast();
  const navigate = useNavigate();
  const { data, isLoading, isError, refetch } = useProjectsQuery(listParams);
  const { data: tagsData } = useTagsQuery();

  useEffect(() => {
    if (isError) {
      toast({
        title: 'Unable to load projects',
        description: 'We could not reach the Hub API. Check connectivity and retry.',
        variant: 'error',
      });
    }
  }, [isError, toast]);

  const projects = useMemo(() => data?.items ?? [], [data]);
  const total = data?.total ?? 0;

  const categories = useMemo(() => {
    const unique = new Map<string, { id: string; name: string; slug: string }>();
    for (const project of projects) {
      if (project.category) {
        unique.set(project.category.slug, project.category);
      }
    }
    return Array.from(unique.values());
  }, [projects]);

  if (isLoading) {
    return (
      <div className="space-y-8">
        <PageHeader title="Projects" description="Exploring knowledge across the homelab." />
        <ProjectsFilters
          filters={filters}
          tags={tagsData ?? []}
          categories={categories}
          onFiltersChange={setFilters}
          isLoading
        />
        <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <ProjectCard.Skeleton key={index} />
          ))}
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <ErrorState
        onRetry={() => {
          toast({ title: 'Retrying project fetch' });
          refetch();
        }}
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="Projects"
        description="Explore documented builds, services, and experiments running in the homelab."
        leading={
          <span className="inline-flex items-center gap-2 rounded-full border border-border/70 bg-surface px-4 py-1 text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
            {total} {total === 1 ? 'project' : 'projects'}
          </span>
        }
      />

      <ProjectsFilters
        filters={filters}
        tags={tagsData ?? []}
        categories={categories}
        onFiltersChange={(nextFilters: ProjectsFiltersState) => {
          const currentNormalized = normalizeFilters(filters);
          const nextNormalized = normalizeFilters(nextFilters);
          const changed =
            JSON.stringify(currentNormalized) !== JSON.stringify(nextNormalized);
          setFilters(nextFilters);
          if (changed) {
            toast({
              title: 'Filters updated',
              description: 'The project list reflects the latest filters.',
              variant: 'success',
            });
          }
        }}
      />

      {projects.length ? (
        <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
          {projects.map((project) => (
            <ProjectCard
              key={project.id}
              project={{
                ...project,
                heroMedia: project.heroMedia ?? undefined,
                category: project.category ?? undefined,
              }}
              onSelect={() => navigate(`/projects/${project.slug}`)}
              footer="Click to open full project notes"
            />
          ))}
        </div>
      ) : (
        <div className="rounded-3xl border border-border/60 bg-surface-raised px-8 py-12 text-center shadow-surface-sm">
          <h3 className="text-xl font-semibold text-text">No projects match your filters</h3>
          <p className="mt-2 text-sm text-text-muted">
            Try adjusting the search term or removing selected tags to see more results.
          </p>
        </div>
      )}
    </div>
  );
};


