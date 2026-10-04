export const locales = ["tk", "ru", "en"] as const;
export type Locale = (typeof locales)[number];
export const defaultLocale: Locale = "tk";
export const LOCALE_COOKIE = "NEXT_LOCALE";
export const PROJECT_COOKIE = "habarchy_project";
export const ACCESS_COOKIE = "habarchy_access";
export const REFRESH_COOKIE = "habarchy_refresh";

export function isLocale(v: string | undefined): v is Locale {
  return !!v && (locales as readonly string[]).includes(v);
}
