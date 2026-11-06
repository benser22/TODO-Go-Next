import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Imagen de Docker mínima: .next/standalone contiene solo lo necesario.
  output: "standalone",
  cacheComponents: true,
  partialPrefetching: true,
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
};

export default nextConfig;
