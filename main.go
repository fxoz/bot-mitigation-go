package main

import (
	"log"
	"time"
	"waffe/antibot"
	"waffe/captcha"
	"waffe/core"
	"waffe/utils"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

var cfg = utils.LoadConfig("config.yml")

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
	app.Post("/.__core_/api/captcha/verify", captcha.VerifyCaptchaRoute)
	app.Get("/.__core_/api/captcha/generate", captcha.GenerateCaptchaRoute)
	app.Get("/.__core_/captcha", captcha.DisplayCaptchaRoute)
	app.All("/*", core.OnRequestHandler)

	if cfg.Server.UseProfiler {
		color.Green("Profiler enabled. Try: go tool pprof http://%s/debug/pprof/profile?seconds=10", cfg.Server.Proxy)
		app.Use(pprof.New())
	}

	log.Fatal(app.Listen(cfg.Server.Proxy))
}
