import type { ReactNode } from 'react';
import { cn } from '../lib/cn';

export interface TagPillProps {
  children: ReactNode;
  color?: string;
  className?: string;
}

export const TagPill = ({ children, color, className }: TagPillProps): JSX.Element => (
  <span
    className={cn(
      'inline-flex items-center gap-1 rounded-full border border-border/70 bg-surface px-3 py-1 text-xs font-medium uppercase tracking-wide text-text-muted',
      className,
    )}
    style={
      color
        ? {
            color,
            borderColor: color,
            backgroundColor: `${color}20`,
          }
        : undefined
    }
  >
    {children}
  </span>
);


