package antibot

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"waffe/utils"
)

var (
	validTokens = make(map[string]time.Time)
	tokenMutex  = sync.Mutex{}
)

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func BotProtectionHandler(w http.ResponseWriter, r *http.Request) {
	token, err := generateRandomToken(16)
	if err != nil {
		http.Error(w, "Internal error generating token", http.StatusInternalServerError)
		return
	}

	tokenMutex.Lock()
	validTokens[token] = time.Now().Add(5 * time.Minute)
	tokenMutex.Unlock()

	cookie := &http.Cookie{
		Name:     "BOT_TOKEN",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	originURL := utils.GetFullPath(r)
	encodedURL := url.QueryEscape(originURL)
	finalURL := fmt.Sprintf("/__verify?to=%s", encodedURL)

	html, err := os.ReadFile("assets/html/bot_protection.html")
	if err != nil {
		http.Error(w, "Internal error reading html file", http.StatusInternalServerError)
		return
	}

	js, err := os.ReadFile("assets/js/bot_protection.js")
	if err != nil {
		http.Error(w, "Internal error reading js file", http.StatusInternalServerError)
		return
	}

	jsWithToken := strings.ReplaceAll(string(js), "{-URL-}", finalURL)
	ofuscatedJs := ObfuscateJS(jsWithToken)

	os.WriteFile("obfuscated.js", []byte(ofuscatedJs), 0644)

	htmlStr := strings.ReplaceAll(string(html), "//JS//", string(ofuscatedJs))
	w.Write([]byte(htmlStr))
}
