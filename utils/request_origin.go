package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

var cfg = LoadConfig("config.yml")

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    time.Duration(cfg.Transfer.TimeoutIdleSeconds) * time.Second,
		DisableCompression: !cfg.Transfer.Compress,
	},
	Timeout: time.Duration(cfg.Transfer.TimeoutSeconds) * time.Second,
}

func RequestOrigin(c *fiber.Ctx) error {
	originURL := cfg.Server.Origin + c.OriginalURL()
	// color.Yellow("Requesting origin server: %s %s", c.Method(), originURL)

	req, err := http.NewRequestWithContext(c.UserContext(), c.Method(), originURL, bytes.NewReader(c.Body()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to create origin request")
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})

	if cfg.Server.RealIpHeader != "" {
		req.Header.Add(cfg.Server.RealIpHeader, c.IP())
	}

	resp, err := client.Do(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to reach origin server")
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	c.Status(resp.StatusCode)

	if c.Method() != http.MethodHead && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		if _, err := io.Copy(c, resp.Body); err != nil {
			fmt.Printf("Error streaming response: %v\n", err)
		}
	}

	return nil
}

func IsOriginAlive(origin string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, origin, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
