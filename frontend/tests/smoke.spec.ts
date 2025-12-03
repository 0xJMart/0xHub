import { expect, test } from '@playwright/test';
import { findProjectBySlug, listProjects, mockTags } from '../src/mocks/data';

test.beforeEach(async ({ page }) => {
  await page.route('http://localhost:8080/projects', async (route, request) => {
    const url = new URL(request.url());
    const filters = {
      tag: url.searchParams.get('tag') ?? undefined,
      category: url.searchParams.get('category') ?? undefined,
      search: url.searchParams.get('search') ?? undefined,
    };
    const response = listProjects(filters);
    await route.fulfill({
      status: 200,
      body: JSON.stringify(response),
      headers: { 'content-type': 'application/json' },
    });
  });

  await page.route('http://localhost:8080/projects/*', async (route, request) => {
    const slug = request.url().split('/').at(-1) ?? '';
    const project = findProjectBySlug(slug);
    if (!project) {
      await route.fulfill({
        status: 404,
        body: JSON.stringify({ error: 'Not found' }),
        headers: { 'content-type': 'application/json' },
      });
      return;
    }
    await route.fulfill({
      status: 200,
      body: JSON.stringify(project),
      headers: { 'content-type': 'application/json' },
    });
  });

  await page.route('http://localhost:8080/tags', async (route) => {
    await route.fulfill({
      status: 200,
      body: JSON.stringify(mockTags),
      headers: { 'content-type': 'application/json' },
    });
  });
});

test('lists projects and navigates to detail view', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Projects' })).toBeVisible();
  await expect(page.getByText('Pi-hole Gateway')).toBeVisible();

  await page.getByText('Pi-hole Gateway').click();

  await expect(page.getByRole('heading', { name: 'Pi-hole Gateway' })).toBeVisible();
  await expect(page.getByText('Network-wide DNS filtering')).toBeVisible();
});

test('shows an error toast when the API fails', async ({ page }) => {
  await page.route('http://localhost:8080/projects', async (route) => {
    await route.fulfill({
      status: 500,
      body: JSON.stringify({ error: 'upstream failure' }),
      headers: { 'content-type': 'application/json' },
    });
  });

  await page.goto('/');
  await expect(page.getByText(/Unable to load projects/i)).toBeVisible();
});


