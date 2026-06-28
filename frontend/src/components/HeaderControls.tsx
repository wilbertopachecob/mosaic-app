import React from "react";
import { Moon, Sun } from "lucide-react";
import { useTranslation } from "react-i18next";
import { AppLanguage, resolveAppLanguage } from "../constants/language";
import { useTheme } from "../context/ThemeContext";

/** Segmented control for switching between English and Spanish. */
const LanguageToggle: React.FC = () => {
  const { i18n, t } = useTranslation();
  const currentLang = resolveAppLanguage(i18n.language);

  const setLanguage = (lang: AppLanguage) => {
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
        className={`lang-toggle-option ${currentLang === AppLanguage.Es ? "is-active" : ""}`}
        aria-pressed={currentLang === AppLanguage.Es}
        onClick={() => setLanguage(AppLanguage.Es)}
      >
        ES
      </button>
      <button
        type="button"
        className={`lang-toggle-option ${currentLang === AppLanguage.En ? "is-active" : ""}`}
        aria-pressed={currentLang === AppLanguage.En}
        onClick={() => setLanguage(AppLanguage.En)}
      >
        EN
      </button>
    </div>
  );
};

/** Button that toggles between light and dark themes. */
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

/** Header toolbar with language and theme controls. */
const HeaderControls: React.FC = () => {
  return (
    <div className="header-controls">
      <LanguageToggle />
      <ThemeToggle />
    </div>
  );
};

export default HeaderControls;
