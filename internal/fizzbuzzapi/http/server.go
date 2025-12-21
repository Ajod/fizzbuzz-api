package http

import (
	"context"
	"fizzbuzz-api/internal/fizzbuzzapi/config"
	"fizzbuzz-api/internal/fizzbuzzapi/controllers"
	"fizzbuzz-api/internal/fizzbuzzapi/handlers"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	HttpServer      *http.Server
	cfg             *config.Config
	log             logger.Logger
	fizzbuzzHandler *handlers.FizzBuzzHandler
}

func NewServer(log logger.Logger) (*Server, error) {
	cfg, err := config.LoadConfig(log)
	if err != nil {
		return nil, err
	}

	// Define and initialize controllers
	fizzbuzzLimits := types.FizzBuzzLimits{
		MaxLimit:        cfg.MaxFizzBuzzLimit,
		MaxStringLength: cfg.MaxStringLength,
	}
	fizzbuzzController := controllers.NewFizzBuzzController(fizzbuzzLimits, log)

	var fizzBuzzStatsController *controllers.FizzBuzzStatsController
	// Note: Currently only in-memory stats recorder is implemented, placeholder for future extensions
	switch cfg.StatsStorage {
	case "inmemory":
		log.Info("using in-memory stats recorder")
		fizzBuzzStatsController = controllers.NewFizzBuzzStatsController(log)
	default:
		log.Info("using in-memory stats recorder (default)")
		fizzBuzzStatsController = controllers.NewFizzBuzzStatsController(log)
	}

	// Define and initialize handlers
	fizzbuzzHandler := handlers.NewFizzBuzzHandler(cfg, log, fizzbuzzController, fizzBuzzStatsController)

	router := gin.Default()
	return &Server{
		HttpServer: &http.Server{Addr: cfg.Host + ":" + cfg.Port, Handler: router},
		cfg:        cfg,
		log:        log,

		fizzbuzzHandler: fizzbuzzHandler,
	}, nil
}

func (s *Server) Run() {
	s.Routes(s.HttpServer.Handler.(*gin.Engine))

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("failed to start HTTP server", "error", err)
		}
	}()
	s.log.Info("fizzbuzz-api server running", "host", s.cfg.Host, "port", s.cfg.Port)

	// Graceful shutdown on interrupt signal (finish requests before shutting down)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	s.log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		s.log.Error("server forced to shutdown", "error", err)
	}
}

func (s *Server) Routes(router *gin.Engine) {
	// Define API routes here
	router.GET("/fizzbuzz/health", handlers.HealthCheck)

	router.POST("/fizzbuzz/generate", s.fizzbuzzHandler.GenerateFizzBuzz)
	router.GET("/fizzbuzz/stats", s.fizzbuzzHandler.GetFizzBuzzStats)
}
