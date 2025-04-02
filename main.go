package main

import (
	"log"
	"waffe/antibot"
	"waffe/utils"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func onRequestHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/.__captcha" || c.Path() == "/.__/api/__judge" {
			return c.Next()
		}

		ip := c.IP()
		if antibot.NeedsVerification(db, ip) {
			if !utils.IsHTMLRequest(c) {
				return c.Status(fiber.StatusForbidden).SendString("Access denied")
			}
			return utils.RenderPage("bot_protection", c)
		}

		if !antibot.IsClientVerified(db, ip) {
			return c.Status(fiber.StatusForbidden).SendString("Checks failed")
		}

		if err := utils.RequestOrigin(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error processing request")
		}

		return nil
	}
}

func captchaHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.RenderPage("captcha", c)
	}
}

func main() {
	cfg := utils.LoadConfig("config.yml")
	db := antibot.InitDB()

	if !utils.IsOriginAlive(cfg.Server.Origin) {
		if !fiber.IsChild() {
			color.Red("Origin server is not reachable! Exiting...")
		}
		return
	}

	if !fiber.IsChild() {
		color.Green("Loaded config file")
		color.Green("Initialized database")
		color.Green("Origin server is reachable")
		color.Green("Server running at http://%s", cfg.Server.Proxy)
		color.Blue("Private IP: %s", utils.GetPrivateIP())
	}

	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	// app.Use(logger.New())
	app.Post("/.__/api/__judge", antibot.JudgeClient(db))
	app.Get("/.__captcha", captchaHandler(db))
	app.All("/*", onRequestHandler(db))

	log.Fatal(app.Listen(cfg.Server.Proxy))
}
