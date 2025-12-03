import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { PageHeader } from '@0xhub/ui';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useProjectQuery } from '@/app/api/hooks';
import { useToast } from '@/app/providers/ToastProvider';
import { ErrorState } from '@/components/ErrorState';
import { LoadingState } from '@/components/LoadingState';
import { ProjectMediaGallery } from '@/features/projects/components/ProjectMediaGallery';
import { ProjectMetaPanel } from '@/features/projects/components/ProjectMetaPanel';

export const ProjectDetailPage = (): JSX.Element => {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const {
    data: project,
    isLoading,
    isError,
    refetch,
  } = useProjectQuery(slug ?? '');

  useEffect(() => {
    if (isError) {
      toast({
        title: 'Unable to load project',
        description: 'We could not fetch the project details. Retry in a moment.',
        variant: 'error',
      });
    }
  }, [isError, toast]);

  useEffect(() => {
    if (!slug) {
      navigate('/', { replace: true });
    }
  }, [navigate, slug]);

  useEffect(() => {
    if (project) {
      document.title = `${project.title} · 0xHub`;
    }
  }, [project]);

  if (isLoading) {
    return <LoadingState message="Loading project details…" />;
  }

  if (isError || !project) {
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
    <div className="space-y-10">
      <PageHeader
        title={project.title}
        description={
          project.summary ?? 'Detailed breakdown of this homelab build, including media and links.'
        }
        leading={
          <span className="inline-flex items-center gap-2 rounded-full border border-border/70 bg-surface px-4 py-1 text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
            {project.category?.name ?? 'Uncategorized'}
          </span>
        }
        actions={
          <button
            type="button"
            className="rounded-full border border-border/60 px-4 py-2 text-sm font-semibold uppercase tracking-[0.28em] text-text-muted transition hover:border-brand/60 hover:text-brand"
            onClick={() => navigate(-1)}
          >
            Back
          </button>
        }
      />

      {project.heroMedia ? (
        <figure className="overflow-hidden rounded-[2.75rem] border border-border/60 bg-surface-raised shadow-surface-md">
          <img
            src={project.heroMedia.storagePath}
            alt={project.heroMedia.description ?? project.heroMedia.originalFilename}
            className="h-[420px] w-full object-cover"
            loading="lazy"
            decoding="async"
          />
          {(project.heroMedia.description || project.heroMedia.originalFilename) && (
            <figcaption className="px-6 py-4 text-sm text-text-muted">
              {project.heroMedia.description ?? project.heroMedia.originalFilename}
            </figcaption>
          )}
        </figure>
      ) : null}

      <div className="grid gap-8 lg:grid-cols-[2fr_1fr]">
        <article className="prose prose-invert max-w-none space-y-6 prose-hr:border-border/40 prose-hr:border-t prose-hr:pt-4 prose-ul:list-disc prose-li:marker:text-brand">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
              h2: ({ children }) => (
                <h2 className="text-2xl font-semibold tracking-tight text-text">{children}</h2>
              ),
              h3: ({ children }) => (
                <h3 className="pt-4 text-xl font-semibold tracking-tight text-text">{children}</h3>
              ),
              a: ({ children, href }) => (
                <a
                  href={href}
                  target="_blank"
                  rel="noreferrer"
                  className="font-medium text-brand underline decoration-dotted underline-offset-4 hover:text-brand-muted"
                >
                  {children}
                </a>
              ),
              code: ({ children }) => (
                <code className="rounded bg-surface-muted/70 px-2 py-1 text-sm text-brand">
                  {children}
                </code>
              ),
            }}
          >
            {project.description ??
              'This project does not include a detailed description yet. Update it to document setup notes, hardware, and lessons learned.'}
          </ReactMarkdown>

          <ProjectMediaGallery media={project.media} />
        </article>

        <ProjectMetaPanel project={project} />
      </div>
    </div>
  );
};


