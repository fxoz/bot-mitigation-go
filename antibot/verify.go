package antibot

import (
	"net/http"
	"net/url"
	"strings"
	"time"
	"waffe/utils"

	"gorm.io/gorm"
)

func VerifyHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("BOT_TOKEN")
	if err != nil {
		http.Error(w, "Bot check failed", http.StatusTooManyRequests)
		return
	}

	token := cookie.Value

	tokenMutex.Lock()
	expiry, exists := validTokens[token]
	if exists {
		delete(validTokens, token)
	}
	tokenMutex.Unlock()

	if !exists || time.Now().After(expiry) {
		http.Error(w, "Bot check failed (invalid)", http.StatusTooManyRequests)
		return
	}

	targetURL := r.URL.Query().Get("to")
	if targetURL == "" {
		http.Error(w, "Bot check failed (no target)", http.StatusTooManyRequests)
		return
	}

	AddClientToWhitelist(db, utils.GetIP(r))

	decodedURL, err := url.QueryUnescape(targetURL)
	if err != nil {
		http.Error(w, "Invalid encoded target URL", http.StatusBadRequest)
		return
	}

	decodedURL = strings.ReplaceAll(decodedURL, "\\", "/")
	target, err := url.Parse(decodedURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	if target.Hostname() == "" {
		http.Redirect(w, r, target.String(), http.StatusFound)
	} else {
		http.Error(w, "Redirect failed (invalid target)", http.StatusTooManyRequests)
	}
}
