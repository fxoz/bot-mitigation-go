package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

func RequestOrigin(responseWriter http.ResponseWriter, originRequest *http.Request) {
	color.Yellow("Requesting origin server: %s %s", originRequest.Method, originRequest.URL)
	client := &http.Client{}
	resp, err := client.Do(originRequest)
	if err != nil {
		http.Error(responseWriter, "Failed to reach origin server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			responseWriter.Header().Add(key, value)
		}
	}
	responseWriter.WriteHeader(resp.StatusCode)

	if originRequest.Method != http.MethodHead && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		_, err = io.Copy(responseWriter, resp.Body)
		if err != nil {
			fmt.Printf("Error copying response body: %v\n", err)
		}
	}
}

func IsOriginAlive(origin string) bool {
	_, err := http.Get(origin)
	return err == nil
}
