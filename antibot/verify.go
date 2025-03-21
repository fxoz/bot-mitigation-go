package antibot

import (
	"net/http"
	"net/url"
	"strings"
	"waffe/utils"

	"gorm.io/gorm"
)

func VerifyHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if !IsValidTokenForIP(db, utils.GetIP(r), token) {
		http.Error(w, "Check failed (invalid token)", http.StatusTooManyRequests)
		return
	}

	SetClientVerified(db, utils.GetIP(r))

	targetURL := r.URL.Query().Get("to")
	if targetURL == "" {
		http.Error(w, "Check failed (no target)", http.StatusTooManyRequests)
		return
	}

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
