import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// Build output goes straight into the Go embed directory.
export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: '../server/internal/webui/dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:8080',
      '/j': 'http://127.0.0.1:8080',
    },
  },
})
