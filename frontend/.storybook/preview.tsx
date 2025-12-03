import type { Preview } from '@storybook/react';
import { initialize, mswDecorator } from 'msw-storybook-addon';
import { handlers } from '../src/mocks/handlers';
import { TestProviders, type TestProvidersProps } from '../src/test-utils/TestProviders';

initialize({ onUnhandledRequest: 'error' });

const preview: Preview = {
  decorators: [
    mswDecorator,
    (Story, context) => {
      const providerProps = context.parameters.providerProps as TestProvidersProps | undefined;
      return (
        <TestProviders {...providerProps}>
          <Story />
        </TestProviders>
      );
    },
  ],
  parameters: {
    layout: 'fullscreen',
    controls: { expanded: true },
    actions: { argTypesRegex: '^on[A-Z].*' },
    msw: {
      handlers,
    },
    backgrounds: {
      default: 'surface',
      values: [
        { name: 'surface', value: '#020617' },
        { name: 'light', value: '#f8fafc' },
      ],
    },
  },
};

export default preview;


