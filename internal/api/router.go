package api

import (
	"net/http"

	"github.com/joelewaldo/go-micro-service/internal/api/config"
	"github.com/joelewaldo/go-micro-service/internal/api/handler"
	"github.com/joelewaldo/go-micro-service/pkg/middleware"
	"github.com/joelewaldo/go-micro-service/pkg/shared"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", shared.Chain(
		http.HandlerFunc(handler.HealthHandler),
		middleware.Logger,
	))

	mux.Handle("GET /subtract/{minuend}/{subtrahend}", shared.Chain(
		http.HandlerFunc(handler.SubtractHandler),
		middleware.Auth(cfg.RSA_KEY),
		middleware.RequireScope("read:api"),
		middleware.Logger,
	))

	return mux
}
