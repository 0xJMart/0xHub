import { useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import type { ListProjectsParams } from '@0xhub/api';

export interface ProjectsFiltersState {
  search?: string;
  tag?: string;
  category?: string;
}

const FILTER_KEYS: Array<keyof ProjectsFiltersState> = ['search', 'tag', 'category'];

export const normalizeFilters = (filters: ProjectsFiltersState): ProjectsFiltersState => {
  const cleaned: ProjectsFiltersState = {};
  FILTER_KEYS.forEach((key) => {
    const value = filters[key];
    if (value && value.trim() !== '') {
      cleaned[key] = value.trim();
    }
  });
  return cleaned;
};

export const useProjectsSearchParams = () => {
  const [searchParams, setSearchParams] = useSearchParams();

  const filters = useMemo<ProjectsFiltersState>(() => {
    const current: ProjectsFiltersState = {};
    FILTER_KEYS.forEach((key) => {
      const value = searchParams.get(key);
      if (value) {
        current[key] = value;
      }
    });
    return current;
  }, [searchParams]);

  const setFilters = useCallback(
    (nextFilters: ProjectsFiltersState) => {
      const cleaned = normalizeFilters(nextFilters);
      const nextSearchParams = new URLSearchParams();
      FILTER_KEYS.forEach((key) => {
        const value = cleaned[key];
        if (value) {
          nextSearchParams.set(key, value);
        }
      });
      setSearchParams(nextSearchParams, { replace: true });
    },
    [setSearchParams],
  );

  const listParams = useMemo<ListProjectsParams>(() => filters, [filters]);

  return {
    filters,
    setFilters,
    listParams,
  };
};


