package antibot

import (
	"encoding/json"
	"io"
	"net/http"
	"waffe/utils"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Checks struct {
	UserAgentFails              bool `json:"userAgentFails"`
	UsesWebDriver               bool `json:"usesWebDriver"`
	SusProperties               bool `json:"susProperties"`
	UsesHeadlessChrome          bool `json:"usesHeadlessChrome"`
	ChromeDiscrepancy           bool `json:"chromeDiscrepancy"`
	LackingCodecSupport         bool `json:"lackingCodecSupport"`
	PlaywrightStealthPixelRatio bool `json:"playwrightStealthPixelRatio"`
	ReportedUserAgent           string
}

type ChecksResponse struct {
	Verified bool `json:"verified"`
}

func JudgeClient(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cfg := utils.LoadConfig("config.yml")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var c Checks
	err = json.Unmarshal(body, &c)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	score := 0

	// score >=1000: it is 100% certain that the client is a bot

	if c.UserAgentFails {
		score += 100
	}
	if c.UsesWebDriver {
		score += 1000
	}
	if c.SusProperties {
		score += 1000
	}
	if c.UsesHeadlessChrome {
		score += 1000
	}
	if c.ChromeDiscrepancy {
		score += 400
	}
	if c.LackingCodecSupport {
		score += 300
	}
	if c.PlaywrightStealthPixelRatio {
		score += 1000
	}
	if !c.UserAgentFails && c.ReportedUserAgent != "" {
		if c.ReportedUserAgent != r.UserAgent() {
			score += 700
		}
	}

	AddClient(db, utils.GetIP(r))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if score >= cfg.AntiBot.Threshold {
		color.Red("Client score: %d above threshold (%d)", score, cfg.AntiBot.Threshold)
		color.Red("This would trigger a CAPTCHA challenge")

		json.NewEncoder(w).Encode(
			ChecksResponse{
				Verified: false,
			},
		)
		return
	}

	SetClientVerified(db, utils.GetIP(r))

	color.Green("Client score: %d below threshold (%d)", score, cfg.AntiBot.Threshold)
	json.NewEncoder(w).Encode(
		ChecksResponse{
			Verified: true,
		},
	)
}
