import path from "node:path"
import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import tailwindcss from "@tailwindcss/vite"
import wails from "@wailsio/runtime/plugins/vite"
import vueDevTools from "vite-plugin-vue-devtools"

export default defineConfig({
  plugins: [vue(), tailwindcss(), wails("./bindings"), vueDevTools()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
