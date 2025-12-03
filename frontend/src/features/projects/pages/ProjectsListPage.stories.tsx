import type { Meta, StoryObj } from '@storybook/react';
import { HttpResponse, http } from 'msw';
import { ProjectsListPage } from './ProjectsListPage';
import { listProjects } from '@/mocks/data';

const meta = {
  title: 'Projects/ProjectsListPage',
  component: ProjectsListPage,
  parameters: {
    layout: 'fullscreen',
  },
} satisfies Meta<typeof ProjectsListPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AllProjects: Story = {};

export const EmptyProjects: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get('http://localhost:8080/projects', ({ request }) => {
          const url = new URL(request.url);
          const filters = {
            tag: url.searchParams.get('tag') ?? undefined,
            category: url.searchParams.get('category') ?? undefined,
            search: url.searchParams.get('search') ?? undefined,
          };
          const response = listProjects(filters);
          return HttpResponse.json({ ...response, items: [], total: 0 });
        }),
      ],
    },
  },
};

export const ErrorState: Story = {
  parameters: {
    msw: {
      handlers: [
        http.get('http://localhost:8080/projects', () =>
          HttpResponse.json({ error: 'Internal server error' }, { status: 500 }),
        ),
      ],
    },
  },
};


