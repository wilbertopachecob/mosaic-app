import { getInitialLanguage } from "@/i18n";
import { AppLanguage } from "@/constants/language";

describe("getInitialLanguage", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("returns a stored language preference", () => {
    localStorage.setItem("mosaic-language", AppLanguage.Es);
    expect(getInitialLanguage()).toBe(AppLanguage.Es);
  });

  it("maps Spanish browser locales to es", () => {
    Object.defineProperty(navigator, "language", {
      value: "es-MX",
      configurable: true,
    });
    expect(getInitialLanguage()).toBe(AppLanguage.Es);
  });

  it("defaults to en for other browser locales", () => {
    Object.defineProperty(navigator, "language", {
      value: "en-US",
      configurable: true,
    });
    expect(getInitialLanguage()).toBe(AppLanguage.En);
  });
});
