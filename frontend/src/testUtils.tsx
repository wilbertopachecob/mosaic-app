import React from "react";
import { render, RenderResult } from "@testing-library/react";
import { AllProviders } from "./AppProviders";

export { AllProviders };

export function renderWithProviders(ui: React.ReactElement): RenderResult {
  return render(<AllProviders>{ui}</AllProviders>);
}
