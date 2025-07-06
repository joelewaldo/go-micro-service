package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	if l.statusCode == 0 {
		l.statusCode = http.StatusOK
	}
	n, err := l.ResponseWriter.Write(b)
	l.size += n
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}

		defer func() {
			if rec := recover(); rec != nil {
				logrus.WithFields(logrus.Fields{
					"method":     r.Method,
					"path":       r.URL.Path,
					"remote":     r.RemoteAddr,
					"panic":      rec,
					"stack":      string(debug.Stack()),
					"user_agent": r.UserAgent(),
				}).Error("panic recovered in handler")

				if lrw.statusCode == 0 {
					http.Error(lrw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		entry := logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      lrw.statusCode,
			"duration_ms": float64(duration.Microseconds()) / 1000.0,
			"size":        lrw.size,
			"remote":      r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})

		switch {
		case lrw.statusCode >= 500:
			entry.Error("server error")
		case lrw.statusCode >= 400:
			entry.Warn("client error")
		default:
			entry.Info("request handled")
		}
	})
}
