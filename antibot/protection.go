package antibot

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"waffe/utils"

	"gorm.io/gorm"
)

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func BotProtectionHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	token, err := generateRandomToken(16)
	if err != nil {
		http.Error(w, "Internal error generating token", http.StatusInternalServerError)
		return
	}

	originURL := utils.GetFullPath(r)
	encodedURL := url.QueryEscape(originURL)
	finalURL := fmt.Sprintf("/__verify?to=%s&token=%s", encodedURL, token)
	AddClientToWhitelist(db, utils.GetIP(r), token)

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

	jsWithToken := strings.ReplaceAll(string(js), "{-URL-}", finalURL)
	ofuscatedJs := ObfuscateJS(jsWithToken)

	htmlStr := strings.ReplaceAll(string(html), "//JS//", string(ofuscatedJs))
	w.Write([]byte(htmlStr))
}
