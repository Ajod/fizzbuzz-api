package main

import (
	"fizzbuzz-api/internal/fizzbuzzapi/http"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
)

func main() {
	logger := logger.NewSlogLogger()

	server, err := http.NewServer(logger)
	if err != nil {
		panic(err)
	}
	server.Run()
	logger.Info("shutting down fizzbuzz-api")
}
