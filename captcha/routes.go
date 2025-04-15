package captcha

import (
	"waffe/utils"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

func DisplayCaptchaRoute(c *fiber.Ctx) error {
	return utils.RenderPage("captcha", c)
}

func GenerateCaptchaRoute(c *fiber.Ctx) error {
	captchaImage := GenerateImageCaptcha()
	if captchaImage.DataUri == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate captcha")
	}

	response := map[string]string{
		"image": captchaImage.DataUri,
	}

	RegisterCaptcha(c.IP(), captchaImage.CorrectRegion)
	return c.JSON(response)
}

func VerifyCaptchaRoute(c *fiber.Ctx) error {
	clientIP := c.IP()

	if ExceededMaxFailedAttempts(clientIP) {
		color.Red("Captcha failed attempts exceeded, IP %s", clientIP)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"verified": false, "exceeded": true})
	}

	if IsVerified(clientIP) {
		color.Blue("Captcha verification not required, IP %s", clientIP)
		return c.Status(fiber.StatusForbidden).SendString("Captcha verification not required")
	}

	var request struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON format")
	}

	clickedX := int(request.X * float32(ScaledWidth))
	clickedY := int(request.Y * float32(ScaledHeight))

	if IsCaptchaCorrect(clientIP, clickedX, clickedY) {
		MarkCaptchaSolved(clientIP)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"verified": true, "exceeded": false})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"verified": false, "exceeded": false})
}
