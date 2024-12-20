import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
// import devtools from 'solid-devtools/vite';

export default defineConfig({
  plugins: [
    /* 
    Uncomment the following line to enable solid-devtools.
    For more info see https://github.com/thetarnav/solid-devtools/tree/main/packages/extension#readme
    */
    // devtools(),
    solidPlugin(),
  ],
  server: {
    port: 3923,
    host: "127.0.0.1",
    proxy: {
      "/api": "http://localhost:3924",
    },
  },
  build: {
    target: "esnext",
    outDir: "../cmd/server/public/web-dist",
    emptyOutDir: true,
  },
});
