package utils

import (
	"net/http"
	"os"
)

func RenderPage(folder string, w http.ResponseWriter, r *http.Request) {
	cfg := LoadConfig("config.yml")
	pathBase := "assets/" + folder

	html, err := os.ReadFile(pathBase + "/index.html")
	if err != nil {
		http.Error(w, "Internal error reading html file", http.StatusInternalServerError)
		return
	}

	js, err := os.ReadFile(pathBase + "/index.js")
	if err != nil {
		http.Error(w, "Internal error reading js file", http.StatusInternalServerError)
		return
	}

	jsFinal := string(js)

	if cfg.Other.ObfuscateJavaScript {
		jsFinal = ObfuscateJS(jsFinal)
	}

	htmlFinal := string(html) + "\n<script>\n" + jsFinal + "\n</script>\n"
	w.Write([]byte(htmlFinal))
}
