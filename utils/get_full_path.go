package utils

import "net/http"

func GetFullPath(r *http.Request) string {
	fullPath := r.URL.Path
	if r.URL.RawQuery != "" {
		fullPath += "?" + r.URL.RawQuery
	}
	return fullPath
}
