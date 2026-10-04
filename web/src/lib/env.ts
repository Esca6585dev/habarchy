/** Server-side configuration (never exposed to the browser). */
export const serverEnv = {
  /** Base URL of the Habarchy API, e.g. http://api:8080 inside docker. */
  apiUrl: process.env.HABARCHY_API_URL ?? "http://localhost:8080",
  /** Set "true" when the panel is served over HTTPS (secure cookies). */
  secureCookies: process.env.HABARCHY_SECURE_COOKIES === "true",
};
