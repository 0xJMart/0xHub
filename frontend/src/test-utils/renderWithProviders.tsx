import type { ReactElement } from 'react';
import { render, type RenderOptions } from '@testing-library/react';
import { TestProviders, type TestProvidersProps } from './TestProviders';

export interface RenderWithProvidersOptions extends RenderOptions {
  providerProps?: TestProvidersProps;
}

export const renderWithProviders = (
  ui: ReactElement,
  options: RenderWithProvidersOptions = {},
) => {
  const { providerProps, ...renderOptions } = options;
  const Wrapper = ({ children }: { children: React.ReactNode }) => (
    <TestProviders {...providerProps}>{children}</TestProviders>
  );
  return render(ui, { wrapper: Wrapper, ...renderOptions });
};


