import { createHubApiClient } from '@0xhub/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export const hubApiClient = createHubApiClient({
  baseUrl: API_BASE_URL,
});


