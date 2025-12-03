import type { Meta, StoryObj } from '@storybook/react';
import { HttpResponse, http } from 'msw';
import { Route, Routes } from 'react-router-dom';
import { ProjectDetailPage } from './ProjectDetailPage';

const meta = {
  title: 'Projects/ProjectDetailPage',
  component: ProjectDetailPage,
  parameters: {
    layout: 'fullscreen',
    providerProps: {
      routeEntries: ['/projects/pi-hole-gateway'],
    },
  },
  decorators: [
    (Story) => (
      <Routes>
        <Route path="/projects/:slug" element={<Story />} />
      </Routes>
    ),
  ],
} satisfies Meta<typeof ProjectDetailPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const MissingProject: Story = {
  name: 'Error – Missing project',
  parameters: {
    providerProps: {
      routeEntries: ['/projects/missing-project'],
    },
    msw: {
      handlers: [
        http.get('http://localhost:8080/projects/:slug', () =>
          HttpResponse.json({ error: 'Not found' }, { status: 404 }),
        ),
      ],
    },
  },
};


