import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it } from 'vitest';
import { listProjectsMock, listTagsMock } from '@/test-utils/mockHubApi';
import { listProjects, mockTags } from '@/mocks/data';
import { ProjectsListPage } from './ProjectsListPage';
import { renderWithProviders } from '@/test-utils/renderWithProviders';

beforeEach(() => {
  listProjectsMock.mockReset().mockImplementation(async (params) => listProjects(params ?? {}));
  listTagsMock.mockReset().mockResolvedValue(mockTags);
});

describe('ProjectsListPage', () => {
  it('renders projects returned from the Hub API', async () => {
    renderWithProviders(<ProjectsListPage />);

    expect(await screen.findByRole('heading', { name: /Projects/i })).toBeInTheDocument();
    expect(await screen.findByText(/Pi-hole Gateway/i)).toBeInTheDocument();
    expect(screen.getByText(/GitOps Control Plane/i)).toBeInTheDocument();
    expect(screen.getByText(/Plex Media Stack/i)).toBeInTheDocument();
  });

  it('filters projects when a tag is toggled', async () => {
    const user = userEvent.setup();
    renderWithProviders(<ProjectsListPage />);

    const tagButton = await screen.findByRole('button', { name: /Automation/i });
    await user.click(tagButton);

    expect(tagButton).toHaveAttribute('aria-pressed', 'true');

    const visibleProjects = await screen.findAllByRole('heading', { level: 3 });
    const titles = visibleProjects.map((heading) => heading.textContent);
    expect(titles).toContain('Pi-hole Gateway');
    expect(titles).not.toContain('Plex Media Stack');
  });
});


