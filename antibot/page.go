package antibot

import (
	"net/http"
	"os"
	"strings"
	"time"
	"waffe/utils"

	"gorm.io/gorm"
)

func ShowVerificationPage(db *gorm.DB, w http.ResponseWriter, r *http.Request, token string, finalURL string) {
	cfg := utils.LoadConfig("config.yml")
	cookieExpiry := time.Now().Add(time.Duration(cfg.AntiBot.TokenValidForSeconds) * time.Second)
	tokenCookie := &http.Cookie{
		Name:    "__t",
		Value:   token,
		Path:    "/",
		Expires: cookieExpiry,
	}
	http.SetCookie(w, tokenCookie)

	finalUrlCookie := &http.Cookie{
		Name:    "__u",
		Value:   finalURL,
		Path:    "/",
		Expires: cookieExpiry,
	}
	http.SetCookie(w, finalUrlCookie)

	html, err := os.ReadFile("assets/bot_protection/index.html")
	if err != nil {
		http.Error(w, "Internal error reading html file", http.StatusInternalServerError)
		return
	}

	js, err := os.ReadFile("assets/bot_protection/index.js")
	if err != nil {
		http.Error(w, "Internal error reading js file", http.StatusInternalServerError)
		return
	}

	ofuscatedJs := ObfuscateJS(string(js))

	htmlStr := strings.ReplaceAll(string(html), "//JS//", string(ofuscatedJs))
	w.Write([]byte(htmlStr))
}
