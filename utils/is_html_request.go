package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func IsHTMLRequest(c *fiber.Ctx) bool {
	method := c.Method()
	if method != fiber.MethodGet && method != fiber.MethodPost {
		return false
	}
	if strings.Contains(c.Get("Accept"), "text/html") {
		return true
	}
	return strings.HasSuffix(c.Path(), ".html")
}
