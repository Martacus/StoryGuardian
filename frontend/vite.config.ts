import path from "node:path"
import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import tailwindcss from "@tailwindcss/vite"
import wails from "@wailsio/runtime/plugins/vite"

export default defineConfig({
  plugins: [vue(), tailwindcss(), wails("./bindings")],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
