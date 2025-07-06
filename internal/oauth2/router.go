package oauth2

import (
	"database/sql"
	"net/http"

	"github.com/joelewaldo/go-micro-service/internal/oauth2/config"
	"github.com/joelewaldo/go-micro-service/internal/oauth2/handler"
	"github.com/joelewaldo/go-micro-service/pkg/middleware"
	"github.com/joelewaldo/go-micro-service/pkg/shared"
)

func NewRouter(config *config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	handle := handler.NewTokenHandler(config, db)

	mux.Handle("POST /oauth2/token", shared.Chain(
		http.HandlerFunc(handle.Handle),
		middleware.Logger,
	))

	return mux
}
