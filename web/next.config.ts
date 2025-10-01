import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  // image: {
  //   remotePatterns: [new URL('http://localhost:3000/image/instructions/**')],
  // },
  experimental: {
    serverActions: {
      bodySizeLimit: '5mb',
    },
  },
};

export default nextConfig;
