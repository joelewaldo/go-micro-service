package main

import (
	"net/http"

	"github.com/joelewaldo/go-micro-service/internal/api"
	"github.com/sirupsen/logrus"
)

func main() {
	router := api.NewRouter()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	addr := "127.0.0.1:3000"
	logrus.WithField("address: ", addr).Info("Starting Server")

	if err := http.ListenAndServe(addr, router); err != nil {
		logrus.WithField("Event", "Start Server").Fatal(err)
	}
}
