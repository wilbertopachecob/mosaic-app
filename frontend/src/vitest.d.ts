import "vitest";

// jest-axe ships matchers that augment Jest's expect. Re-declare them for
// Vitest's expect so `toHaveNoViolations()` is recognized under type-check.
declare module "vitest" {
  interface Assertion<T = unknown> {
    toHaveNoViolations(): T;
  }
  interface AsymmetricMatchersContaining {
    toHaveNoViolations(): void;
  }
}
