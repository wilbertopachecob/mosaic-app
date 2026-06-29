/** Supported UI locales for the mosaic app. */
export enum AppLanguage {
  En = "en",
  Es = "es",
}

/** Normalizes an i18next language code to a supported app locale. */
export function resolveAppLanguage(languageCode: string | undefined): AppLanguage {
  return (languageCode ?? AppLanguage.En).startsWith(AppLanguage.Es)
    ? AppLanguage.Es
    : AppLanguage.En;
}
