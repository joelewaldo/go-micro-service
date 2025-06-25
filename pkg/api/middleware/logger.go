package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration,
			"remote":   r.RemoteAddr,
		}).Info("Handled Request")
	})
}
