package main

import (
	"log"
	"strings"
	"time"
	"waffe/antibot"
	"waffe/captcha"
	"waffe/utils"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

var cfg = utils.LoadConfig("config.yml")

func onRequestHandler(c *fiber.Ctx) error {
	if strings.HasPrefix(c.Path(), "/debug/pprof") && cfg.Server.UseProfiler {
		return c.Next()
	}

	if c.Path() == "/.__core_/captcha" || c.Path() == "/.__core_/api/judge" {
		return c.Next()
	}

	ip := c.IP()
	if antibot.RequiresReVerification(ip) {
		if !utils.IsHTMLRequest(c) {
			return c.Status(fiber.StatusForbidden).SendString("Access denied")
		}
		return utils.RenderPage("bot_protection", c)
	}

	if !antibot.IsClientCurrentlyVerified(ip) {
		return c.Status(fiber.StatusForbidden).SendString("Checks failed")
	}

	if err := utils.RequestOrigin(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error processing request")
	}

	return nil
}

func captchaHandler(c *fiber.Ctx) error {
	return utils.RenderPage("captcha", c)
}

func generateCaptchaHandler(c *fiber.Ctx) error {
	captchaImage := captcha.GenerateImageCaptcha()
	if captchaImage.DataUri == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate captcha")
	}

	response := map[string]string{
		"dataUri":       captchaImage.DataUri,
		"correctRegion": "",
	}
	return c.JSON(response)
}

func main() {
	if !utils.IsOriginAlive(cfg.Server.Origin) {
		if !fiber.IsChild() {
			color.Red("Origin server is not reachable! Exiting...")
		}
		return
	}

	if !fiber.IsChild() {
		color.Green("Server running at http://%s", cfg.Server.Proxy)
		color.Green("Private IP: %s", utils.GetPrivateIP())
	}

	app := fiber.New(fiber.Config{
		Prefork:          cfg.Server.Prefork,
		StrictRouting:    cfg.Server.StrictRouting,
		CaseSensitive:    cfg.Server.CaseSensitive,
		BodyLimit:        cfg.Server.BodyLimitBytes,
		ReadTimeout:      time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:     time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:      time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
		ProxyHeader:      cfg.Server.GetIpFromHeader,
		DisableKeepalive: !cfg.Server.EnableKeepAlive,
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
	})

	app.Use(logger.New())
	app.Post("/.__core_/api/judge", antibot.JudgeClient())
	app.Post("/.__core_/api/captcha/generate", generateCaptchaHandler)
	app.Get("/.__core_/captcha", captchaHandler)
	app.All("/*", onRequestHandler)

	if cfg.Server.UseProfiler {
		color.Green("Profiler enabled. Try: go tool pprof http://%s/debug/pprof/profile?seconds=10", cfg.Server.Proxy)
		app.Use(pprof.New())
	}

	log.Fatal(app.Listen(cfg.Server.Proxy))
}
