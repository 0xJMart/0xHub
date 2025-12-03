import type { Config } from 'tailwindcss';
import { tokens } from './tokens';

const { colors, radii, shadows, spacing, typography, transitions } = tokens;

export const tailwindExtend = {
  colors: {
    brand: {
      DEFAULT: colors.brand.primary,
      muted: colors.brand.primaryMuted,
      secondary: colors.brand.secondary,
      accent: colors.brand.accent,
    },
    surface: {
      DEFAULT: colors.surface.base,
      raised: colors.surface.raised,
      muted: colors.surface.muted,
      overlay: colors.surface.overlay,
    },
    border: {
      subtle: colors.border.subtle,
      DEFAULT: colors.border.default,
      strong: colors.border.strong,
    },
    text: {
      DEFAULT: colors.text.primary,
      secondary: colors.text.secondary,
      muted: colors.text.muted,
      inverted: colors.text.inverted,
    },
    state: {
      success: colors.state.success,
      info: colors.state.info,
      warning: colors.state.warning,
      danger: colors.state.danger,
    },
  },
  borderRadius: {
    xs: radii.xs,
    sm: radii.sm,
    md: radii.md,
    lg: radii.lg,
    xl: radii.xl,
    pill: radii.pill,
  },
  fontFamily: {
    sans: typography.fonts.sans,
    mono: typography.fonts.mono,
  },
  letterSpacing: {
    tight: typography.tracking.tight,
    wide: typography.tracking.wide,
  },
  boxShadow: {
    'surface-sm': shadows.sm,
    'surface-md': shadows.md,
    'surface-lg': shadows.lg,
  },
  maxWidth: {
    content: spacing.container,
  },
  spacing: {
    gutter: spacing.gutter,
    stack: spacing.stackGap,
  },
  transitionTimingFunction: {
    subtle: transitions.subtle,
    emphasized: transitions.emphasized,
  },
};

export const tailwindPreset = {
  darkMode: ['class'],
  theme: {
    extend: tailwindExtend,
  },
  plugins: [],
} satisfies Config;

export type TailwindPreset = typeof tailwindPreset;


