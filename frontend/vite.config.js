/**
 * Vite构建配置文件
 * 配置开发服务器、代理、路径别名和生产环境构建选项
 */
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()], // 注册Vue插件
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)) // @别名指向src目录
    }
  },
  server: {
    host: '0.0.0.0', // 允许外部网络访问开发服务器
    port: 5173,
    proxy: {
      // 开发环境代理：将/api请求转发到后端服务，解决跨域问题
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: 'dist',              // 生产构建输出目录
    assetsDir: 'assets',         // 静态资源存放子目录
    sourcemap: false,            // 生产环境不生成sourcemap
    minify: 'terser',            // 使用terser进行代码压缩
    chunkSizeWarningLimit: 1500, // chunk体积警告阈值（KB）
    rollupOptions: {
      output: {
        // 代码分割策略：将第三方库拆分为独立chunk，优化缓存和加载速度
        manualChunks: {
          vendor: ['vue', 'vue-router'],
          utils: ['axios']
        }
      }
    }
  },
  preview: {
    host: '0.0.0.0',
    port: 4173 // 预览生产构建的端口
  }
})
