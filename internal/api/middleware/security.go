package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/cbt-platform/internal/security"
)

func SecurityMiddleware(waf *security.WAF, rateLimiter *security.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// WAF check
		if !waf.CheckRequest(c.Request()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Request blocked by security policy",
			})
		}

		// Rate limiting
		clientIP := c.IP()
		if !rateLimiter.Allow(clientIP) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}

		// Security headers
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return c.Next()
	}
}
