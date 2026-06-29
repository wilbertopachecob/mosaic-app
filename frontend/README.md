# Mosaic App Frontend

This is the React + TypeScript frontend for Mosaic, built with [Vite](https://vite.dev/)
and tested with [Vitest](https://vitest.dev/).

## Available Scripts

In the project directory, you can run:

### `npm run dev` (alias: `npm start`)

Runs the app in development mode with Vite. Open
[http://localhost:3000](http://localhost:3000) to view it in the browser.

API requests to `/api` are proxied to the backend at `http://localhost:8080`
(configured in `vite.config.ts`). The page reloads automatically on edits via
Fast Refresh.

### `npm run build`

Type-checks the project (`tsc --noEmit`) and builds an optimized production
bundle to the `dist` folder.

### `npm run preview`

Serves the production build locally so you can verify it before deploying.

### `npm test`

Launches Vitest in watch mode. Use `npm run test:run` for a single
non-watching run, `npm run test:coverage` for coverage, and `npm run test:a11y`
for the accessibility suite.

### `npm run lint` / `npm run lint:fix`

Lints (and optionally auto-fixes) the `src` directory with ESLint.

### `npm run type-check`

Runs the TypeScript compiler in no-emit mode to check types.

## Path aliases

Imports use the `@/` alias, which resolves to `src/`. The alias is configured
in both `vite.config.ts` and `tsconfig.json`.

## Learn More

- [Vite documentation](https://vite.dev/guide/)
- [Vitest documentation](https://vitest.dev/guide/)
- [React documentation](https://react.dev/)
