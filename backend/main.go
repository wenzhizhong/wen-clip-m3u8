package backend

import (
	"clipM3u8Media/backend/routers"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Backend() {
	router := gin.Default()
	// router.Use(cors.Default()) // 使用默认配置
	router.Use(cors.New(corsConfig()))

	routers.Init(router) // 路由
	router.Run(":34116")
}

func corsConfig() cors.Config {
	// 完整的 CORS 配置
	return cors.Config{
		// 允许的域名，支持正则表达式
		AllowOrigins: []string{
			"*",
		},

		// 自定义源验证函数
		AllowOriginFunc: func(origin string) bool {
			// 允许所有本地请求
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// 允许所有本地请求
			if strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			// 允许特定 IP 地址
			if strings.HasPrefix(origin, "http://192.168.") {
				return true
			}

			return false
		},

		// 允许的 HTTP 方法
		AllowMethods: []string{
			"GET",     // 读取数据
			"POST",    // 创建数据
			"PUT",     // 更新数据
			"PATCH",   // 部分更新
			"DELETE",  // 删除数据
			"HEAD",    // 获取头部信息
			"OPTIONS", // 预检请求
		},

		// 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"X-Client-Version",
			"X-Platform",
			"Cache-Control",
			"Pragma",
		},

		// 允许客户端访问的响应头
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},

		// 是否允许携带凭证（cookies、HTTP认证等）
		AllowCredentials: true,

		// 预检请求缓存时间（秒）
		MaxAge: 12 * time.Hour,

		// 允许客户端访问的响应头
		AllowWildcard: true,

		// 允许浏览器暴露响应头给 JavaScript
		AllowBrowserExtensions: false,

		// 允许 WebSocket
		AllowWebSockets: true,

		// 允许 Files
		AllowFiles: false,
	}

}
