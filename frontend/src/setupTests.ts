// jest-dom adds custom matchers for asserting on DOM nodes, e.g.:
// expect(element).toHaveTextContent(/react/i)
// The /vitest subpath wires the matchers into Vitest's expect.
// learn more: https://github.com/testing-library/jest-dom
import "@testing-library/jest-dom/vitest";
import { toHaveNoViolations } from "jest-axe";

// `expect` is available as a global because `globals: true` is set in the
// Vitest config (see vite.config.ts).
expect.extend(toHaveNoViolations);

// Node 22+ exposes an experimental, getter-only `globalThis.localStorage`
// accessor that returns undefined (and warns) unless `--localstorage-file` is
// provided. Because it has no setter, the DOM test environment cannot replace
// it, which breaks any code that reads `localStorage`. Define a working
// in-memory implementation for tests instead.
function createMemoryStorage(): Storage {
  const store = new Map<string, string>();
  return {
    get length() {
      return store.size;
    },
    clear: () => store.clear(),
    getItem: (key: string) => (store.has(key) ? store.get(key)! : null),
    key: (index: number) => Array.from(store.keys())[index] ?? null,
    removeItem: (key: string) => {
      store.delete(key);
    },
    setItem: (key: string, value: string) => {
      store.set(key, String(value));
    },
  } as Storage;
}

for (const storageName of ["localStorage", "sessionStorage"] as const) {
  Object.defineProperty(globalThis, storageName, {
    value: createMemoryStorage(),
    configurable: true,
    writable: true,
  });
}
