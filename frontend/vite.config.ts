import path from "node:path";
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  server: {
    port: 3000,
    // Mirrors the former CRA "proxy" field: forward API calls to the backend.
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
  build: {
    // Keep CRA's output folder so existing deploy scripts (scripts/build.sh,
    // main.go static serving) continue to work without changes.
    outDir: "build",
  },
  test: {
    globals: true,
    environment: "happy-dom",
    setupFiles: "./src/setupTests.ts",
    css: true,
    coverage: {
      provider: "v8",
      reporter: ["text", "lcov"],
    },
  },
});
