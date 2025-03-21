package utils

import (
	"net/http"
)

func ProcessRequest(r *http.Request) *http.Request {
	cfg := LoadConfig("config.yml")
	targetURL := cfg.Server.Origin + GetFullPath(r)
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return nil
	}

	req.Header = r.Header.Clone()

	clientIP := r.RemoteAddr
	req.Header.Add(cfg.Server.RealIpHeader, clientIP)

	return req
}
