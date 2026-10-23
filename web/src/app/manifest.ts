import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "Habarçy SMS",
    short_name: "Habarçy SMS",
    description:
      "Operator SMPP (0796) we telefon şlýuzy arkaly SMS, OTP we habar ugratmak.",
    start_url: "/",
    display: "standalone",
    background_color: "#eef1f5",
    theme_color: "#1366c4",
    icons: [{ src: "/icon.svg", sizes: "any", type: "image/svg+xml" }],
  };
}
