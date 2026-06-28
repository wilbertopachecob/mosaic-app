import React from "react";
import { render, RenderResult } from "@testing-library/react";
import { AllProviders } from "./AppProviders";

export { AllProviders };

/**
 * Renders a component wrapped in the same providers used by the application.
 */
export function renderWithProviders(ui: React.ReactElement): RenderResult {
  return render(<AllProviders>{ui}</AllProviders>);
}
