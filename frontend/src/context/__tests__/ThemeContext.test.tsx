import { renderHook, act } from "@testing-library/react";
import {
  ThemeProvider,
  useTheme,
  getInitialTheme,
  applyTheme,
} from "@/context/ThemeContext";

describe("getInitialTheme", () => {
  beforeEach(() => {
    localStorage.clear();
    delete document.documentElement.dataset.theme;
  });

  it("returns a stored theme preference", () => {
    localStorage.setItem("mosaic-theme", "dark");
    expect(getInitialTheme()).toBe("dark");
  });

  it("falls back to system preference when nothing is stored", () => {
      window.matchMedia = vi.fn().mockImplementation((query: string) => ({
        matches: query === "(prefers-color-scheme: dark)",
        media: query,
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
      }));

    expect(getInitialTheme()).toBe("dark");
  });

  it("defaults to light when storage and matchMedia are unavailable", () => {
    const originalMatchMedia = window.matchMedia;
    // @ts-expect-error simulate environments without matchMedia
    window.matchMedia = undefined;

    expect(getInitialTheme()).toBe("light");

    window.matchMedia = originalMatchMedia;
  });
});

describe("applyTheme", () => {
  it("sets the data-theme attribute on the document root", () => {
    applyTheme("dark");
    expect(document.documentElement.dataset.theme).toBe("dark");
  });
});

describe("useTheme", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("throws when used outside ThemeProvider", () => {
    expect(() => renderHook(() => useTheme())).toThrow(
      "useTheme must be used within ThemeProvider"
    );
  });

  it("toggles between light and dark", () => {
    localStorage.setItem("mosaic-theme", "light");

    const { result } = renderHook(() => useTheme(), {
      wrapper: ThemeProvider,
    });

    expect(result.current.theme).toBe("light");

    act(() => {
      result.current.toggleTheme();
    });

    expect(result.current.theme).toBe("dark");
    expect(localStorage.getItem("mosaic-theme")).toBe("dark");
    expect(document.documentElement.dataset.theme).toBe("dark");
  });
});
