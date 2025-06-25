package api

import (
	"net/http"

	"github.com/joelewaldo/go-micro-service/pkg/api/handler"
	"github.com/joelewaldo/go-micro-service/pkg/api/middleware"
	"github.com/joelewaldo/go-micro-service/pkg/shared"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/health", shared.Chain(
		http.HandlerFunc(handler.HealthHandler),
		middleware.Logger,
	))

	return mux
}
