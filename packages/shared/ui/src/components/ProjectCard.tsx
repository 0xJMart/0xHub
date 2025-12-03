import type { ReactNode } from 'react';
import { cn } from '../lib/cn';
import { TagPill } from './TagPill';

export interface ProjectLink {
  id: string;
  url: string;
  label?: string;
}

export interface ProjectTag {
  id: string;
  name: string;
  displayName: string;
  color?: string;
}

export interface ProjectMedia {
  id: string;
  originalFilename: string;
  storagePath: string;
  description?: string;
}

export interface ProjectCategory {
  id: string;
  name: string;
  slug: string;
}

export type ProjectStatus = 'draft' | 'active' | 'archived' | 'deprecated';

export interface ProjectSummary {
  id: string;
  title: string;
  slug: string;
  summary?: string;
  heroMedia?: ProjectMedia | null;
  tags: ProjectTag[];
  category?: ProjectCategory | null;
  status?: ProjectStatus;
}

export interface ProjectCardProps {
  project: ProjectSummary;
  onSelect?: (project: ProjectSummary) => void;
  footer?: ReactNode;
  className?: string;
}

const statusCopy: Record<ProjectStatus, { label: string; tone: string }> = {
  draft: { label: 'Draft', tone: 'border-warning/60 text-warning' },
  active: { label: 'Active', tone: 'border-success/60 text-success' },
  archived: { label: 'Archived', tone: 'border-border/60 text-text-muted' },
  deprecated: { label: 'Deprecated', tone: 'border-danger/60 text-danger' },
};

export const ProjectCard = ({ project, onSelect, footer, className }: ProjectCardProps) => {
  const status = project.status ? statusCopy[project.status] : null;

  const handleClick = (): void => {
    if (onSelect) {
      onSelect(project);
    }
  };

  return (
    <article
      className={cn(
        'group flex h-full flex-col overflow-hidden rounded-3xl border border-border/70 bg-surface-raised shadow-surface-sm transition hover:-translate-y-1 hover:shadow-surface-md',
        onSelect && 'cursor-pointer',
        className,
      )}
      onClick={onSelect ? handleClick : undefined}
    >
      {project.heroMedia ? (
        <div className="relative aspect-[16/9] overflow-hidden">
          <img
            src={project.heroMedia.storagePath}
            alt={project.heroMedia.description ?? project.heroMedia.originalFilename}
            className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
            loading="lazy"
            decoding="async"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-surface via-surface/10 to-transparent opacity-75" />
        </div>
      ) : (
        <div className="flex aspect-[16/9] items-center justify-center bg-surface-muted/60">
          <div className="rounded-xl border border-border/70 px-4 py-2 text-xs uppercase tracking-wide text-text-muted">
            No media yet
          </div>
        </div>
      )}

      <div className="flex flex-1 flex-col gap-6 px-6 py-6">
        <div className="flex items-center gap-3">
          {project.category ? (
            <TagPill className="text-[0.625rem] uppercase tracking-[0.26em] text-text-muted">
              {project.category.name}
            </TagPill>
          ) : null}
          {status ? (
            <TagPill className={cn('text-[0.625rem]', status.tone)}>{status.label}</TagPill>
          ) : null}
        </div>

        <div className="space-y-3">
          <h3 className="text-xl font-semibold tracking-tight text-text">{project.title}</h3>
          {project.summary ? (
            <p className="line-clamp-3 text-sm leading-relaxed text-text-muted">{project.summary}</p>
          ) : null}
        </div>

        {project.tags.length ? (
          <div className="flex flex-wrap gap-2">
            {project.tags.map((tag) => (
              <TagPill key={tag.id} color={tag.color}>
                {tag.displayName}
              </TagPill>
            ))}
          </div>
        ) : null}

        {footer ? <div className="mt-auto pt-3 text-sm text-text-muted">{footer}</div> : null}
      </div>
    </article>
  );
};

ProjectCard.Skeleton = function ProjectCardSkeleton(): JSX.Element {
  return (
    <article className="flex h-full animate-pulse flex-col overflow-hidden rounded-3xl border border-border/60 bg-surface-raised/70">
      <div className="aspect-[16/9] bg-surface-muted/40" />
      <div className="flex flex-1 flex-col gap-4 px-6 py-6">
        <div className="h-6 w-1/3 rounded-full bg-surface-muted/60" />
        <div className="space-y-2">
          <div className="h-6 w-3/4 rounded bg-surface-muted/70" />
          <div className="h-4 w-full rounded bg-surface-muted/60" />
          <div className="h-4 w-5/6 rounded bg-surface-muted/60" />
        </div>
        <div className="mt-auto h-5 w-2/5 rounded-full bg-surface-muted/50" />
      </div>
    </article>
  );
};


