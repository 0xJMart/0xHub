import type { PropsWithChildren } from 'react';
import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { ToastProvider } from '@/app/providers/ToastProvider';
import { GlobalErrorBoundary } from '@/app/providers/GlobalErrorBoundary';

export interface TestProvidersProps {
  routeEntries?: string[];
}

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

export const TestProviders = ({
  children,
  routeEntries = ['/'],
}: PropsWithChildren<TestProvidersProps>): JSX.Element => {
  const [client] = useState(() => createTestQueryClient());

  return (
    <GlobalErrorBoundary>
      <QueryClientProvider client={client}>
        <MemoryRouter initialEntries={routeEntries}>
          <ToastProvider>{children}</ToastProvider>
        </MemoryRouter>
      </QueryClientProvider>
    </GlobalErrorBoundary>
  );
};


