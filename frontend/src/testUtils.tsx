import React from "react";
import { render, RenderResult } from "@testing-library/react";
import { I18nextProvider } from "react-i18next";
import i18n from "./i18n";
import { ThemeProvider } from "./context/ThemeContext";

export function AllProviders({ children }: { children: React.ReactNode }) {
  return (
    <I18nextProvider i18n={i18n}>
      <ThemeProvider>{children}</ThemeProvider>
    </I18nextProvider>
  );
}

export function renderWithProviders(ui: React.ReactElement): RenderResult {
  return render(<AllProviders>{ui}</AllProviders>);
}
