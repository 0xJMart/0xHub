import type { Project } from '@0xhub/api';
import { TagPill } from '@0xhub/ui';
import { useToast } from '@/app/providers/ToastProvider';

const statusCopy: Record<Project['status'], { label: string; tone: string }> = {
  draft: { label: 'Draft', tone: 'border-state-warning/60 text-state-warning' },
  active: { label: 'Active', tone: 'border-state-success/60 text-state-success' },
  archived: { label: 'Archived', tone: 'border-border/60 text-text-muted' },
  deprecated: { label: 'Deprecated', tone: 'border-state-danger/60 text-state-danger' },
};

export interface ProjectMetaPanelProps {
  project: Project;
}

export const ProjectMetaPanel = ({ project }: ProjectMetaPanelProps): JSX.Element => {
  const { toast } = useToast();

  return (
    <aside className="space-y-8 rounded-3xl border border-border/60 bg-surface-raised px-6 py-6 shadow-surface-sm">
      <section className="space-y-2">
        <h3 className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">Status</h3>
        <TagPill className={statusCopy[project.status].tone}>{statusCopy[project.status].label}</TagPill>
      </section>

      {project.category ? (
        <section className="space-y-2">
          <h3 className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">
            Category
          </h3>
          <TagPill className="border-brand/50 bg-brand/5 text-brand">{project.category.name}</TagPill>
        </section>
      ) : null}

      {project.tags.length ? (
        <section className="space-y-3">
          <h3 className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">Tags</h3>
          <div className="flex flex-wrap gap-2">
            {project.tags.map((tag) => (
              <TagPill key={tag.id} color={tag.color}>
                {tag.displayName}
              </TagPill>
            ))}
          </div>
        </section>
      ) : null}

      {project.links.length ? (
        <section className="space-y-3">
          <h3 className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">Links</h3>
          <ul className="space-y-2">
            {project.links.map((link) => (
              <li key={link.id}>
                <a
                  href={link.url}
                  target="_blank"
                  rel="noreferrer"
                  className="group flex items-center justify-between gap-3 rounded-2xl border border-border/60 bg-surface px-4 py-3 text-sm text-text transition hover:border-brand/60 hover:text-brand"
                >
                  <span className="truncate">{link.label ?? link.url}</span>
                  <span className="text-xs uppercase tracking-[0.28em] text-text-muted group-hover:text-brand">
                    Open
                  </span>
                </a>
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      <section className="space-y-3">
        <h3 className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">Share</h3>
        <button
          type="button"
          className="w-full rounded-full border border-border/60 px-4 py-2 text-sm font-semibold uppercase tracking-[0.28em] text-text-muted transition hover:border-brand/60 hover:text-brand"
          onClick={async () => {
            try {
              await navigator.clipboard.writeText(window.location.href);
              toast({
                title: 'Link copied',
                description: 'You can now share this project with the homelab crew.',
                variant: 'success',
              });
            } catch (error) {
              console.error(error);
              toast({
                title: 'Unable to copy link',
                description: 'Use your browser menu to copy the URL manually.',
                variant: 'error',
              });
            }
          }}
        >
          Copy project URL
        </button>
      </section>
    </aside>
  );
};


