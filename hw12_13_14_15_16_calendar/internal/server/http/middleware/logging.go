package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server"
)

// NewLoggingMiddleware создает middleware для логирования сетевых запросов.
func NewLoggingMiddleware(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()
			ip := strings.Split(r.RemoteAddr, ":")[0]
			ts := start.Format(server.LogTimestampFormat)
			method := r.Method
			path := r.URL.RequestURI()
			userAgent := r.UserAgent()

			rl := &responseWriterProxy{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rl, r)

			latency := time.Since(start).Milliseconds()
			code := rl.statusCode

			logLine := fmt.Sprintf("%s [%s] %s %s %s %d %d %q", ip, ts, method, path, r.Proto, code, latency, userAgent)
			logger.Println(logLine)
		})
	}
}

// responseWriterProxy оборачивает http.ResponseWriter для перехвата кода ответа.
type responseWriterProxy struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает код ответа и сохраняет его в структуре для дальнейшего использования в middleware.
func (rw *responseWriterProxy) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
