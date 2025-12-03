import type { Config } from 'tailwindcss'
import { tailwindExtend } from '@0xhub/ui'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  darkMode: ['class'],
  theme: {
    extend: tailwindExtend,
  },
  plugins: [],
} satisfies Config
