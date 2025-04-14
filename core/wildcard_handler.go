package core

import (
	"strings"
	"waffe/antibot"
	"waffe/captcha"
	"waffe/utils"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

var cfg = utils.LoadConfig("config.yml")

func isIgnoredPath(path string) bool {
	if strings.HasPrefix(path, "/.__core_") {
		return true
	}
	if strings.HasPrefix(path, "/debug/pprof") && cfg.Server.UseProfiler {
		return true
	}

	return false
}

func OnRequestHandler(c *fiber.Ctx) error {
	if isIgnoredPath(c.Path()) {
		return c.Next()
	}

	ip := c.IP()

	color.Blue("IsVerified: %s", captcha.IsVerified(ip))

	if captcha.IsVerified(ip) || antibot.IsVerified(ip) {
		if err := utils.RequestOrigin(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error processing request")
		}

		return nil
	}

	if !utils.IsHTMLRequest(c) {
		return c.Status(fiber.StatusForbidden).SendString("Access denied")
	}
	return utils.RenderPage("antibot", c)
}
