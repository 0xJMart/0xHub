import { Link } from 'react-router-dom';

export const NotFoundPage = (): JSX.Element => (
  <div className="flex min-h-screen flex-col items-center justify-center gap-6 bg-surface text-center">
    <div className="inline-flex rounded-full border border-border/70 bg-surface-raised px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
      404
    </div>
    <div className="space-y-3 px-6">
      <h1 className="text-3xl font-semibold text-text">Page not found</h1>
      <p className="text-sm text-text-muted">
        The view you&apos;re looking for doesn&apos;t exist. It may have been moved or archived.
      </p>
    </div>
    <Link
      to="/"
      className="rounded-full border border-brand/70 bg-brand/10 px-4 py-2 text-sm font-medium text-brand transition hover:bg-brand/20"
    >
      Back to projects
    </Link>
  </div>
);


