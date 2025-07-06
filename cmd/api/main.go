package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joelewaldo/go-micro-service/internal/api"
	"github.com/joelewaldo/go-micro-service/internal/api/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.ExtractConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	lvl, err := logrus.ParseLevel(cfg.LOG_LEVEL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid LOG_LEVEL %q: %v\n", cfg.LOG_LEVEL, err)
		os.Exit(1)
	}
	logrus.SetLevel(lvl)

	if err != nil {
		logrus.Error(err)
		return
	}

	router := api.NewRouter(cfg)

	bindAddr := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)
	logrus.WithField("Event", "Start Server").
		Infof("listening on %s (issuer=%s)", bindAddr, cfg.ISSUER)

	if err := http.ListenAndServe(bindAddr, router); err != nil {
		logrus.WithField("Event", "Start Server").Fatal(err)
	}
}
