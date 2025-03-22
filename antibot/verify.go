package antibot

import (
	"net/http"
	"net/url"
	"waffe/utils"

	"gorm.io/gorm"
)

func VerifyHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("__t")

	if err != nil || !IsValidTokenForIP(db, utils.GetIP(r), tokenCookie.Value) {
		http.Error(w, "Check failed (token)", http.StatusForbidden)
		return
	}

	finalURLCookie, err := r.Cookie("__u")
	if err != nil {
		http.Error(w, "Invalid final URL", http.StatusBadRequest)
		return
	}

	target, err := url.Parse(finalURLCookie.Value)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	if target.Hostname() != "" {
		http.Error(w, "Redirect failed (invalid target)", http.StatusBadRequest)
		return
	}

	SetClientVerified(db, utils.GetIP(r))
	http.Redirect(w, r, target.String(), http.StatusFound)
}
