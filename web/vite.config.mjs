import { defineConfig } from "vite";
import uni from "@dcloudio/vite-plugin-uni";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    uni.default ? uni.default() : uni()
  ],
  server: {
    host: '0.0.0.0', // 允许局域网访问
    port: 5173,
    hmr: {
      overlay: true,
      protocol: 'ws',
      host: '0.0.0.0'
    },
    fs: {
      strict: false
    },
    cors: true
  },
  optimizeDeps: {
    exclude: []
  }
});
