package antibot

import (
	"net/http"
	"net/url"
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

	http.Redirect(w, r, decodedURL, http.StatusFound)

}
