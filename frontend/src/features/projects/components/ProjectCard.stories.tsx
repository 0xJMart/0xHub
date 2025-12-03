import type { Meta, StoryObj } from '@storybook/react';
import { ProjectCard, type ProjectSummary } from '@0xhub/ui';
import { mockProjects, toProjectSummary } from '@/mocks/data';

const summaries: ProjectSummary[] = mockProjects.map(toProjectSummary);

const meta = {
  title: 'Projects/ProjectCard',
  component: ProjectCard,
  parameters: {
    layout: 'centered',
  },
  args: {
    project: summaries[0],
  },
} satisfies Meta<typeof ProjectCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const ActiveProject: Story = {
  args: {
    project: summaries[0],
    footer: 'Click to open full project notes',
  },
};

export const DraftProject: Story = {
  args: {
    project: summaries[2],
  },
};

export const WithoutMedia: Story = {
  args: {
    project: {
      ...summaries[1],
      heroMedia: undefined,
    },
  },
};


