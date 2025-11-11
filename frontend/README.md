# Hub Frontend

React + TypeScript + Vite workspace for the Hub UI. Tooling is preconfigured for Tailwind CSS, ESLint (flat config), Prettier, Vitest, and Storybook so feature teams can focus on building modules instead of plumbing.

## Requirements

- Node.js ≥ 20.19.0 (local binaries are vendored under `/home/admin/node-v20.19.0-linux-x64/bin`)
- npm ≥ 10.8.0

Install dependencies with:

```bash
npm install
```

## Commands

- `npm run dev` – start the Vite dev server with HMR
- `npm run lint` – lint the project with ESLint + Storybook rules
- `npm run typecheck` – static type analysis with TypeScript
- `npm run format` – check formatting using Prettier (Tailwind plugin enabled)
- `npm run test` – run Vitest in CI mode (headless)
- `npm run test:watch` – run Vitest in watch/interactive mode
- `npm run storybook` – launch Storybook (Vite builder)
- `npm run storybook:build` – generate the static Storybook site for CI publishing
- `npm run build` – run TypeScript build references and bundle via Vite

> Storybook’s Vitest addon attempts to install Playwright browsers automatically. If the install fails due to missing system dependencies, run `npx playwright install chromium` once after prerequisites are met.

## Project Layout

- `src/` – application entry point and feature shells
- `src/stories/` – Storybook examples and component documentation
- `tailwind.config.js` – Tailwind theme + content scanning configuration
- `prettier.config.js` – shared formatting rules (Tailwind-aware)
- `.storybook/` – Storybook configuration, including Vitest test integration

Shared UI primitives live in `packages/shared/ui` and will be expanded in later phases alongside design tokens and component libraries.
