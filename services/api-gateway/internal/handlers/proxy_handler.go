package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"Phantom_backend/services/api-gateway/internal/config"
)

func ProxyRequest(cfg *config.Config, serviceURL string) http.HandlerFunc {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		url := serviceURL + r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}

		var body io.Reader
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			body = bytes.NewBuffer(bodyBytes)
		}

		req, err := http.NewRequest(r.Method, url, body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
