package utils

import (
	"net/http"
	"path"
	"strings"
)

func IsHTMLRequest(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		return false
	}

	if accept := r.Header.Get("Accept"); strings.Contains(accept, "text/html") {
		return true
	}

	ext := path.Ext(r.URL.Path)
	return ext == "" || ext == ".html"
}
