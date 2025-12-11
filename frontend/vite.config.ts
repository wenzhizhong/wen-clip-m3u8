import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // server: {
  //   proxy: {
  //     // 配置以 /proxy/video 开头的请求直接返回 404，使其被 Wails 的 AssetServer 处理
  //     '/proxy': {
  //       target: 'http://localhost:34115', // 这里可以随便填一个地址，主要是触发bypass规则
  //       bypass: function (req, res, options) {
  //         // 直接返回 null 或 false，Vite 会返回 404 给 Wails
  //         // Wails 收到 404 后，就会使用你自定义的 Go Handler
  //         return null; // 或者 return false
  //       },
  //     },
  //   },
  // },
  // server: {
  //   proxy: {
  //     // 关键修改：使用 bypass 代替 target
  //     '^/proxy/video': {
  //       target: 'http://localhost:34115', // 这个地址现在只是占位符，实际不会被连接
  //       changeOrigin: true,
  //       // bypass 函数决定如何处理请求。返回 null 表示不代理，让请求“落空”
  //       bypass: (req, res, options) => {
  //         // 当请求路径是 /proxy/video 时，Vite直接返回null，Wails会接管
  //         return null;
  //       },
  //       rewrite: (path) => path.replace(/^\/proxy\/video/, '/proxy/video'),
  //     },
  //     // 你可以用同样的方式处理其他需要Go后端处理的API
  //     // '^/api': { ... }
  //   },
  // },
  server: {
    proxy: {
      // // 将以下开头的请求，直接代理到Wails后端，避免Vite处理
      // '^/proxy': {
      //   target: 'http://localhost:34115', // Wails后端地址
      //   changeOrigin: true,
      //   // 关键：确保ws (WebSocket) 也被代理，如果不需要可省略
      //   ws: true,
      // },
      // 你可以添加更多需要代理的前端API路径
      '^/api': {
        target: "http://localhost:34116",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    }
  }
})
