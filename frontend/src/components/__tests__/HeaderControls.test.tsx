import { screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom";
import HeaderControls from "@/components/HeaderControls";
import { renderWithProviders } from "@/testUtils";
import i18n from "@/i18n";
import { AppLanguage } from "@/constants/language";

describe("HeaderControls", () => {
  beforeEach(async () => {
    await i18n.changeLanguage(AppLanguage.En);
    localStorage.setItem("mosaic-theme", "light");
  });

  it("renders language and theme controls", () => {
    renderWithProviders(<HeaderControls />);

    expect(screen.getByRole("group", { name: /language/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /switch to dark mode/i })).toBeInTheDocument();
  });

  it("switches language when a locale button is clicked", async () => {
    renderWithProviders(<HeaderControls />);

    fireEvent.click(screen.getByRole("button", { name: "ES" }));
    expect(i18n.language).toMatch(/^es/);
  });

  it("toggles theme when the theme button is clicked", () => {
    renderWithProviders(<HeaderControls />);

    const themeButton = screen.getByRole("button", { name: /switch to dark mode/i });
    fireEvent.click(themeButton);

    expect(document.documentElement.dataset.theme).toBe("dark");
    expect(
      screen.getByRole("button", { name: /switch to light mode/i })
    ).toBeInTheDocument();
  });
});
