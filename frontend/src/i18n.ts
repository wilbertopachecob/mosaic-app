import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { AppLanguage } from "@/constants/language";
import en from "@/locales/en.json";
import es from "@/locales/es.json";

const STORAGE_KEY = "mosaic-language";

/**
 * Resolves the initial UI language from localStorage or the browser locale.
 * Falls back to English when no supported preference is found.
 */
export function getInitialLanguage(): string {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === AppLanguage.En || stored === AppLanguage.Es) {
    return stored;
  }

  const browserLang = navigator.language.toLowerCase();
  return browserLang.startsWith(AppLanguage.Es) ? AppLanguage.Es : AppLanguage.En;
}

const initialLanguage = getInitialLanguage();
document.documentElement.lang = initialLanguage;

i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    es: { translation: es },
  },
  lng: initialLanguage,
  fallbackLng: AppLanguage.En,
  interpolation: {
    escapeValue: false,
  },
});

i18n.on("languageChanged", (lng) => {
  localStorage.setItem(STORAGE_KEY, lng);
  document.documentElement.lang = lng;
});

export default i18n;
