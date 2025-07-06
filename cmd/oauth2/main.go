package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joelewaldo/go-micro-service/internal/oauth2"
	"github.com/joelewaldo/go-micro-service/internal/oauth2/config"
	"github.com/joelewaldo/go-micro-service/internal/oauth2/db"
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

	sqlDB, err := db.Connect(cfg)
	if err != nil {
		logrus.Fatalf("could not connect to db: %v", err)
	}
	defer sqlDB.Close()

	router := oauth2.NewRouter(cfg, sqlDB)

	addr := cfg.ISSUER
	logrus.WithField("address: ", addr).Info("Starting Server")

	if err := http.ListenAndServe(addr, router); err != nil {
		logrus.WithField("Event", "Start Server").Fatal(err)
	}
}
