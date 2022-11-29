package analytics

import (
	"net/http"
	"time"

	"github.com/tom-draper/api-analytics/analytics/go/core"
)

// http.ResponseWriter wrapper to store status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func createResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func Analytics(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			ctx := r.Context()
			rw := createResponseWriter(w) // Wrap to store status code

			start := time.Now()
			next.ServeHTTP(rw, r.WithContext(ctx))
			elapsed := time.Since(start).Milliseconds()

			data := core.Data{
				APIKey:       apiKey,
				Hostname:     r.Host,
				Path:         r.URL.Path,
				UserAgent:    r.UserAgent(),
				Method:       core.MethodMap[r.Method],
				ResponseTime: elapsed,
				Status:       rw.status,
				Framework:    7,
			}

			go core.LogRequest(data)
		})
	}
}
