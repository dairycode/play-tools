import { defineConfig } from "vite";
import uni from "@dcloudio/vite-plugin-uni";
import tailwindcss from "@tailwindcss/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss({
      // 明确指定 CSS 文件路径
      base: '/src/style.css'
    }),
    uni.default ? uni.default() : uni()
  ],
  server: {
    host: '0.0.0.0', // 允许局域网访问
    port: 5173,
  }
});
