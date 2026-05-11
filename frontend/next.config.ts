import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,
  async redirects() {
    return [
      { source: "/en/:path*", destination: "/:path*", permanent: false },
    ];
  },
  async rewrites() {
    const englishPages = ["dashboard", "inventory", "alerts", "import", "reports", "settings", "login"];
    return englishPages.flatMap((page) => [
      { source: `/${page}`, destination: `/en/${page}` },
      { source: `/${page}/:path*`, destination: `/en/${page}/:path*` },
    ]);
  },
};

export default withNextIntl(nextConfig);
