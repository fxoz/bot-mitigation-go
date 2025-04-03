package antibot

import (
	"waffe/utils"

	"github.com/gofiber/fiber/v2"
)

type Checks struct {
	UserAgentFails              bool   `json:"userAgentFails"`
	UsesWebDriver               bool   `json:"usesWebDriver"`
	SusProperties               bool   `json:"susProperties"`
	UsesHeadlessChrome          bool   `json:"usesHeadlessChrome"`
	ChromeDiscrepancy           bool   `json:"chromeDiscrepancy"`
	LackingCodecSupport         bool   `json:"lackingCodecSupport"`
	PlaywrightStealthPixelRatio bool   `json:"playwrightStealthPixelRatio"`
	ReportedUserAgent           string `json:"reportedUserAgent"`
}

type ChecksResponse struct {
	Verified bool `json:"verified"`
}

func JudgeClient() fiber.Handler {
	cfg := utils.LoadConfig("config.yml")
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			return fiber.NewError(fiber.StatusMethodNotAllowed, "Invalid request method")
		}

		var checks Checks
		if err := c.BodyParser(&checks); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
		}

		score := 0
		if checks.UserAgentFails {
			score += 100
		}
		if checks.UsesWebDriver {
			score += 1000
		}
		if checks.SusProperties {
			score += 1000
		}
		if checks.UsesHeadlessChrome {
			score += 1000
		}
		if checks.ChromeDiscrepancy {
			score += 400
		}
		if checks.LackingCodecSupport {
			score += 300
		}
		if checks.PlaywrightStealthPixelRatio {
			score += 1000
		}
		if !checks.UserAgentFails && checks.ReportedUserAgent != "" && checks.ReportedUserAgent != c.Get("User-Agent") {
			score += 700
		}

		RegisterClient(c.IP())
		if score >= cfg.AntiBot.Threshold {
			return c.Status(fiber.StatusOK).JSON(ChecksResponse{Verified: false})
		}

		MarkClientVerified(c.IP())
		return c.Status(fiber.StatusOK).JSON(ChecksResponse{Verified: true})
	}
}
