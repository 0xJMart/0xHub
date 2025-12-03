import { Route, Routes } from 'react-router-dom';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getProjectMock, listProjectsMock } from '@/test-utils/mockHubApi';
import { findProjectBySlug, listProjects } from '@/mocks/data';
import { ProjectDetailPage } from './ProjectDetailPage';
import { renderWithProviders } from '@/test-utils/renderWithProviders';

const renderDetailRoute = (route: string) =>
  renderWithProviders(
    <Routes>
      <Route path="/projects/:slug" element={<ProjectDetailPage />} />
    </Routes>,
    { providerProps: { routeEntries: [route] } },
  );

beforeEach(() => {
  listProjectsMock.mockReset().mockImplementation(async (params) => listProjects(params ?? {}));
  getProjectMock.mockReset().mockImplementation(async (slug) => {
    const project = findProjectBySlug(slug);
    if (!project) {
      throw new Error('Project not found');
    }
    return project;
  });
});

describe('ProjectDetailPage', () => {
  it('shows project metadata and content', async () => {
    renderDetailRoute('/projects/pi-hole-gateway');

    expect(await screen.findByRole('heading', { name: /Pi-hole Gateway/i })).toBeInTheDocument();
    expect(screen.getByText(/Network-wide DNS filtering/i)).toBeInTheDocument();
    expect(screen.getByRole('img', { name: /Pi-hole query dashboard/i })).toBeInTheDocument();
    expect(screen.getByText(/Kubernetes deployment with health probes/i)).toBeInTheDocument();
  });

  it('surfaces an error state when the project cannot be fetched', async () => {
    getProjectMock.mockRejectedValueOnce(new Error('Not found'));

    renderDetailRoute('/projects/missing');

    expect(await screen.findByText(/Something went wrong/i)).toBeInTheDocument();
    const retryButton = screen.getByRole('button', { name: /Retry/i });
    expect(retryButton).toBeInTheDocument();
  });

  it('copies the project URL to the clipboard', async () => {
    const user = userEvent.setup();
    const writeText = vi.fn().mockResolvedValue(undefined);
    const originalClipboard = navigator.clipboard;

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    renderDetailRoute('/projects/pi-hole-gateway');

    const copyButton = await screen.findByRole('button', { name: /Copy project URL/i });
    await user.click(copyButton);

    expect(writeText).toHaveBeenCalledWith(expect.stringContaining('http://localhost'));

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: originalClipboard,
    });
  });
});


