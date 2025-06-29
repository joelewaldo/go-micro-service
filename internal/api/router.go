package api

import (
	"net/http"

	"github.com/joelewaldo/go-micro-service/internal/api/handler"
	"github.com/joelewaldo/go-micro-service/internal/api/middleware"
	"github.com/joelewaldo/go-micro-service/pkg/shared"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", shared.Chain(
		http.HandlerFunc(handler.HealthHandler),
		middleware.Logger,
	))

	mux.Handle("GET /subtract/{minuend}/{subtrahend}", shared.Chain(
		http.HandlerFunc(handler.SubtractHandler),
		middleware.Logger,
	))

	return mux
}
