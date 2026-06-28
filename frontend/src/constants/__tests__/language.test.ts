import { AppLanguage, resolveAppLanguage } from "../language";

describe("resolveAppLanguage", () => {
  it("maps Spanish locale codes to Es", () => {
    expect(resolveAppLanguage("es")).toBe(AppLanguage.Es);
    expect(resolveAppLanguage("es-MX")).toBe(AppLanguage.Es);
  });

  it("maps other locale codes to En", () => {
    expect(resolveAppLanguage("en")).toBe(AppLanguage.En);
    expect(resolveAppLanguage("en-US")).toBe(AppLanguage.En);
    expect(resolveAppLanguage(undefined)).toBe(AppLanguage.En);
  });
});
