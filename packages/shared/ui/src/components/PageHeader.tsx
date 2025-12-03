import type { ReactNode } from 'react';
import { cn } from '../lib/cn';

export interface PageHeaderProps {
  title: string;
  description?: ReactNode;
  leading?: ReactNode;
  actions?: ReactNode;
  className?: string;
}

export const PageHeader = ({
  title,
  description,
  leading,
  actions,
  className,
}: PageHeaderProps): JSX.Element => (
  <div
    className={cn(
      'relative overflow-hidden rounded-3xl border border-border/80 bg-surface-raised/90 px-6 py-8 shadow-surface-md sm:px-8 md:rounded-[2.5rem]',
      className,
    )}
  >
    <div className="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-col gap-4">
        {leading}
        <div className="space-y-2">
          <h1 className="text-3xl font-semibold tracking-tight text-text sm:text-4xl">{title}</h1>
          {description ? (
            <div className="text-base leading-7 text-text-muted">{description}</div>
          ) : null}
        </div>
      </div>
      {actions ? <div className="flex flex-wrap items-center gap-3">{actions}</div> : null}
    </div>
  </div>
);


