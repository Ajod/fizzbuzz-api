package main

import (
	"context"
	"fizzbuzz-api/internal/fizzbuzzapi/http"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logger.NewSlogLogger()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		server, err := http.NewServer(logger)
		if err != nil {
			panic(err)
		}
		server.Run()
	}()
	<-ctx.Done()
	logger.Info("shutting down fizzbuzz-api")
}
