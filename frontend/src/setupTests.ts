// jest-dom adds custom matchers for asserting on DOM nodes, e.g.:
// expect(element).toHaveTextContent(/react/i)
// The /vitest subpath wires the matchers into Vitest's expect.
// learn more: https://github.com/testing-library/jest-dom
import "@testing-library/jest-dom/vitest";
import { expect } from "vitest";
import { toHaveNoViolations } from "jest-axe";

expect.extend(toHaveNoViolations);
