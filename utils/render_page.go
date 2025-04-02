package utils

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func RenderPage(folder string, c *fiber.Ctx) error {
	cfg := LoadConfig("config.yml")
	pathBase := "assets/" + folder

	html, err := os.ReadFile(pathBase + "/index.html")
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Page not found")
	}

	js, err := os.ReadFile(pathBase + "/index.js")
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Page not found")
	}

	jsFinal := string(js)

	if cfg.Other.ObfuscateJavaScript {
		jsFinal = ObfuscateJS(jsFinal)
	}

	htmlFinal := string(html) + "\n<script>\n" + jsFinal + "\n</script>\n"

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(htmlFinal)
}
