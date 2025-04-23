import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
// import path from 'nodepath'
// import { fileURLToPath } from 'url'
import path from "path"
import { fileURLToPath } from "url"

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": path.resolve(path.dirname(fileURLToPath(import.meta.url)), "./src"),
    },
  },
})
