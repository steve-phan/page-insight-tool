import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'www.home24.de',
        pathname: '/**',
      },
    ],
  },
};

export default nextConfig;
