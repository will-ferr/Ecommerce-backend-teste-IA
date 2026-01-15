package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent referrer leakage
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		c.Writer.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"font-src 'self'; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none';")

		// HSTS (only in HTTPS)
		if c.Request.TLS != nil {
			c.Writer.Header().Set("Strict-Transport-Security",
				"max-age=31536000; includeSubDomains; preload")
		}

		// Prevent caching of sensitive data
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")

		// Limit request size
		c.Request.ParseMultipartForm(32 << 20) // 32MB max

		c.Next()
	}
}
