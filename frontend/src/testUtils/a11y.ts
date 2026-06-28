import { axe, type AxeResults } from "jest-axe";
import type { RenderResult } from "@testing-library/react";

/** Runs axe against a rendered container and returns the audit results. */
export async function runAxeAudit(
  container: HTMLElement
): Promise<AxeResults> {
  return axe(container);
}

/** Asserts the rendered UI has no axe accessibility violations. */
export async function expectNoA11yViolations(
  renderResult: Pick<RenderResult, "container">
): Promise<void> {
  const results = await runAxeAudit(renderResult.container);
  expect(results).toHaveNoViolations();
}
