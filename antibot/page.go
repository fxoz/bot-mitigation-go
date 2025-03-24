package antibot

import (
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
)

func StartJavaScriptVerification(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// cfg := utils.LoadConfig("config.yml")

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
