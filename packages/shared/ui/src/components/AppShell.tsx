import type { PropsWithChildren, ReactNode } from 'react';
import { cn } from '../lib/cn';

export interface AppShellNavLink {
  label: string;
  href?: string;
  icon?: ReactNode;
  isActive?: boolean;
  isExternal?: boolean;
  onClick?: () => void;
}

export interface AppShellProps extends PropsWithChildren {
  brand: {
    name: string;
    description?: string;
    logo?: ReactNode;
  };
  navigation?: AppShellNavLink[];
  actions?: ReactNode;
  footer?: ReactNode;
  className?: string;
}

export const AppShell = ({
  brand,
  navigation,
  actions,
  footer,
  className,
  children,
}: AppShellProps): JSX.Element => (
  <div
    className={cn(
      'min-h-screen bg-surface/95 text-text antialiased backdrop-blur-sm',
      'selection:bg-brand/40 selection:text-text',
      className,
    )}
  >
    <div className="relative flex min-h-screen flex-col">
      <header className="border-b border-border/80 bg-surface/80 backdrop-blur">
        <div className="mx-auto flex w-full max-w-content items-center justify-between gap-6 px-6 py-4 lg:px-10">
          <div className="flex items-center gap-3">
            <span className="inline-flex h-10 w-10 items-center justify-center rounded-lg border border-border/80 bg-surface-raised shadow-surface-sm">
              {brand.logo ?? (
                <span className="text-lg font-semibold tracking-wide text-brand">0x</span>
              )}
            </span>
            <div>
              <p className="text-sm font-semibold uppercase tracking-wide text-brand">
                {brand.name}
              </p>
              {brand.description ? (
                <p className="text-sm text-text-muted">{brand.description}</p>
              ) : null}
            </div>
          </div>
          {navigation?.length ? (
            <nav className="hidden items-center gap-1 rounded-full border border-border/80 bg-surface-raised px-1 py-1 shadow-surface-sm md:flex">
              {navigation.map((item, index) => {
                const LinkTag = item.onClick ? 'button' : 'a';
                const props = item.onClick
                  ? {
                      type: 'button' as const,
                      onClick: item.onClick,
                    }
                  : {
                      href: item.href ?? '#',
                      target: item.isExternal ? '_blank' : undefined,
                      rel: item.isExternal ? 'noreferrer' : undefined,
                    };
                return (
                  <LinkTag
                    key={item.href ?? `${item.label}-${index}`}
                    className={cn(
                      'flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition',
                      item.isActive
                        ? 'bg-brand/10 text-brand shadow-surface-sm'
                        : 'text-text-muted hover:bg-surface-muted/60 hover:text-text-secondary',
                    )}
                    {...props}
                  >
                    {item.icon}
                    <span>{item.label}</span>
                  </LinkTag>
                );
              })}
            </nav>
          ) : null}
          <div className="hidden items-center gap-3 md:flex">{actions}</div>
        </div>
      </header>

      <main className="flex-1">
        <div className="mx-auto w-full max-w-content px-6 py-10 lg:px-10 lg:py-12">
          {children}
        </div>
      </main>

      {footer ? (
        <footer className="border-t border-border/80 bg-surface/70">
          <div className="mx-auto flex w-full max-w-content items-center justify-between px-6 py-5 text-sm text-text-muted lg:px-10">
            {footer}
          </div>
        </footer>
      ) : null}
    </div>
  </div>
);


