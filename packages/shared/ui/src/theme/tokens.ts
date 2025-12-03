/**
 * Shared design tokens for the Hub UI.
 * These tokens power the Tailwind theme extension and can be consumed
 * directly by component libraries or other styling systems.
 */

export const colors = {
  brand: {
    primary: '#38bdf8',
    primaryMuted: '#0ea5e9',
    secondary: '#c084fc',
    accent: '#f97316',
  },
  surface: {
    base: '#020617',
    raised: '#0f172a',
    muted: '#1e293b',
    overlay: '#111827cc',
  },
  border: {
    subtle: '#1f2937',
    default: '#334155',
    strong: '#475569',
  },
  text: {
    primary: '#e2e8f0',
    secondary: '#cbd5f5',
    muted: '#94a3b8',
    inverted: '#020617',
  },
  state: {
    success: '#34d399',
    info: '#38bdf8',
    warning: '#fbbf24',
    danger: '#f87171',
  },
} as const;

export const radii = {
  xs: '0.375rem',
  sm: '0.5rem',
  md: '0.75rem',
  lg: '1rem',
  xl: '1.5rem',
  pill: '9999px',
} as const;

export const shadows = {
  sm: '0 8px 24px -12px rgba(15, 23, 42, 0.45)',
  md: '0 16px 40px -20px rgba(15, 23, 42, 0.55)',
  lg: '0 24px 70px -32px rgba(8, 47, 73, 0.55)',
} as const;

export const spacing = {
  gutter: '1.5rem',
  container: '74rem',
  stackGap: '1.75rem',
} as const;

export const typography = {
  fonts: {
    sans: [
      'InterVariable',
      'Inter',
      'system-ui',
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'sans-serif',
    ],
    mono: ['"JetBrains Mono"', '"Fira Code"', 'ui-monospace', 'SFMono-Regular', 'monospace'],
  },
  tracking: {
    tight: '-0.01em',
    normal: '0',
    wide: '0.08em',
  },
} as const;

export const transitions = {
  subtle: '150ms ease-out',
  emphasized: '220ms cubic-bezier(0.16, 1, 0.3, 1)',
} as const;

export const tokens = {
  colors,
  radii,
  shadows,
  spacing,
  typography,
  transitions,
} as const;

export type ThemeTokens = typeof tokens;


