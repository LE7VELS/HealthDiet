package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 只允许配置的前端来源，不使用通配符开放任意跨域访问。
// 同源请求通常没有 Origin，仍应继续处理；跨域预检则必须精确匹配允许来源。
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// 只有匹配的来源才写响应头，浏览器不会把其他来源的响应暴露给页面脚本。
		if origin != "" && origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			// 预检请求不会进入业务 Handler，来源不匹配时直接拒绝。
			if origin != allowedOrigin {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
