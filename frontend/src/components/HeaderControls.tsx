import React from "react";
import { Moon, Sun } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useTheme } from "../context/ThemeContext";

const LanguageToggle: React.FC = () => {
  const { i18n, t } = useTranslation();
  const currentLang = (i18n.language ?? "en").startsWith("es") ? "es" : "en";

  const setLanguage = (lang: "en" | "es") => {
    if (lang !== currentLang) {
      i18n.changeLanguage(lang);
    }
  };

  return (
    <div
      className="lang-toggle"
      role="group"
      aria-label={t("controls.language")}
    >
      <span
        className="lang-toggle-slider"
        data-lang={currentLang}
        aria-hidden="true"
      />
      <button
        type="button"
        className={`lang-toggle-option ${currentLang === "es" ? "is-active" : ""}`}
        aria-pressed={currentLang === "es"}
        onClick={() => setLanguage("es")}
      >
        ES
      </button>
      <button
        type="button"
        className={`lang-toggle-option ${currentLang === "en" ? "is-active" : ""}`}
        aria-pressed={currentLang === "en"}
        onClick={() => setLanguage("en")}
      >
        EN
      </button>
    </div>
  );
};

const ThemeToggle: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const { t } = useTranslation();
  const isDark = theme === "dark";

  return (
    <button
      type="button"
      className="theme-toggle"
      onClick={toggleTheme}
      aria-label={isDark ? t("controls.themeDark") : t("controls.themeLight")}
      title={isDark ? t("controls.themeDark") : t("controls.themeLight")}
    >
      {isDark ? (
        <Moon size={18} aria-hidden="true" />
      ) : (
        <Sun size={18} aria-hidden="true" />
      )}
    </button>
  );
};

const HeaderControls: React.FC = () => {
  return (
    <div className="header-controls">
      <LanguageToggle />
      <ThemeToggle />
    </div>
  );
};

export default HeaderControls;
