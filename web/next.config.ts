import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

const api = process.env.HABARCHY_API_URL ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  experimental: { serverActions: { bodySizeLimit: "2mb" } },
  async rewrites() {
    // Swagger UI and asynqmon live on the API; expose them under the panel's
    // origin so the admin cookie (habarchy_access) authenticates them.
    return [
      { source: "/api-docs", destination: `${api}/api/docs` },
      { source: "/api/docs/:path*", destination: `${api}/api/docs/:path*` },
      { source: "/admin/queues", destination: `${api}/admin/queues` },
      { source: "/admin/queues/:path*", destination: `${api}/admin/queues/:path*` },
    ];
  },
};

export default withNextIntl(nextConfig);
