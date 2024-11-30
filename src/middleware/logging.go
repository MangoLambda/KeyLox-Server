package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// LogRequestResponse is a middleware that logs the request and response bodies
func LogRequestResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request body
		var requestBody bytes.Buffer
		tee := io.TeeReader(r.Body, &requestBody)
		body, err := io.ReadAll(tee)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
		}
		r.Body = io.NopCloser(&requestBody)
		log.Printf("Request: %s %s\nBody: %s", r.Method, r.URL, body)

		// Capture the response body
		rec := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
		next.ServeHTTP(rec, r)

		// Log the response body
		log.Printf("Response: %d\nBody: %s", rec.status, rec.body.String())
	})
}

// responseRecorder is a custom ResponseWriter that captures the response body
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}
