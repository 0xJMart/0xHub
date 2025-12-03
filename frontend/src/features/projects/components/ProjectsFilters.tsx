import type { Tag } from '@0xhub/api';
import { useEffect, useMemo, useState } from 'react';
import { TagPill } from '@0xhub/ui';
import { cn } from '@0xhub/ui';
import type { ProjectsFiltersState } from '@/features/projects/utils/useProjectsSearchParams';

interface ProjectsFiltersProps {
  filters: ProjectsFiltersState;
  tags: Tag[];
  categories: Array<{ id: string; name: string; slug: string }>;
  onFiltersChange: (filters: ProjectsFiltersState) => void;
  isLoading?: boolean;
}

export const ProjectsFilters = ({
  filters,
  tags,
  categories,
  onFiltersChange,
  isLoading = false,
}: ProjectsFiltersProps): JSX.Element => {
  const [searchValue, setSearchValue] = useState(filters.search ?? '');

  useEffect(() => {
    // Synchronize local search input state when the URL query params change.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setSearchValue(filters.search ?? '');
  }, [filters.search]);

  const hasActiveFilters = useMemo(
    () => Boolean(filters.search || filters.tag || filters.category),
    [filters.category, filters.search, filters.tag],
  );

  const applyFilter = (next: ProjectsFiltersState): void => {
    onFiltersChange(next);
    if (next.search !== searchValue) {
      setSearchValue(next.search ?? '');
    }
  };

  return (
    <div className="space-y-6 rounded-3xl border border-border/60 bg-surface-raised px-6 py-6 shadow-surface-sm sm:px-8">
      <form
        className="flex flex-col gap-4 md:flex-row md:items-center md:gap-6"
        onSubmit={(event) => {
          event.preventDefault();
          applyFilter({ ...filters, search: searchValue.trim() || undefined });
        }}
      >
        <div className="flex-1">
          <label htmlFor="project-search" className="sr-only">
            Search projects
          </label>
          <div className="flex items-center gap-3 rounded-full border border-border/70 bg-surface px-4 py-2 focus-within:border-brand focus-within:shadow-surface-sm">
            <span className="text-sm font-semibold uppercase tracking-[0.28em] text-text-muted">
              Search
            </span>
            <input
              id="project-search"
              type="search"
              value={searchValue}
              onChange={(event) => setSearchValue(event.target.value)}
              placeholder="Search by title, tag, or summary"
              className="flex-1 bg-transparent text-sm text-text outline-none placeholder:text-text-muted/70"
              disabled={isLoading}
            />
          </div>
        </div>
        <div className="flex items-center gap-3">
          <button
            type="submit"
            className="rounded-full border border-brand/70 bg-brand/10 px-4 py-2 text-sm font-semibold uppercase tracking-[0.28em] text-brand transition hover:bg-brand/20 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={isLoading}
          >
            Apply
          </button>
          {hasActiveFilters ? (
            <button
              type="button"
              className="rounded-full border border-border/70 px-4 py-2 text-sm font-semibold uppercase tracking-[0.28em] text-text-muted transition hover:text-text"
              onClick={() => applyFilter({})}
              disabled={isLoading}
            >
              Clear
            </button>
          ) : null}
        </div>
      </form>

      {tags.length ? (
        <fieldset className="space-y-3">
          <legend className="text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
            Tags
          </legend>
          <div className="flex flex-wrap gap-3">
            {tags.map((tag) => {
              const isActive = filters.tag === tag.name;
              return (
                <button
                  key={tag.id}
                  type="button"
                  className={cn(
                    'transition',
                    isActive ? 'scale-105' : 'hover:scale-105',
                  )}
                  aria-pressed={isActive}
                  onClick={() =>
                    applyFilter({
                      ...filters,
                      tag: isActive ? undefined : tag.name,
                    })
                  }
                  disabled={isLoading}
                >
                  <TagPill
                    color={tag.color}
                    className={cn(
                      'border-border/60 px-4 py-1.5 text-xs font-semibold uppercase tracking-[0.28em]',
                      isActive
                        ? 'border-brand bg-brand/10 text-brand'
                        : 'text-text-muted hover:border-brand/50',
                    )}
                  >
                    {tag.displayName}
                  </TagPill>
                </button>
              );
            })}
          </div>
        </fieldset>
      ) : null}

      {categories.length ? (
        <fieldset className="space-y-3">
          <legend className="text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
            Categories
          </legend>
          <div className="flex flex-wrap gap-3">
            {categories.map((category) => {
              const isActive = filters.category === category.slug;
              return (
                <button
                  key={category.id}
                  type="button"
                  className="rounded-full border border-border/60 bg-surface px-4 py-1.5 text-xs font-semibold uppercase tracking-[0.28em] text-text-muted transition hover:border-brand/50 hover:text-brand disabled:cursor-not-allowed disabled:opacity-60"
                  aria-pressed={isActive}
                  onClick={() =>
                    applyFilter({
                      ...filters,
                      category: isActive ? undefined : category.slug,
                    })
                  }
                  disabled={isLoading}
                >
                  {category.name}
                </button>
              );
            })}
          </div>
        </fieldset>
      ) : null}
    </div>
  );
};


