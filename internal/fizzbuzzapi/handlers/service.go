package handlers

import (
	"errors"
	"fizzbuzz-api/internal/fizzbuzzapi/config"
	"fizzbuzz-api/internal/fizzbuzzapi/controllers"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FizzBuzzHandler struct {
	cfg           *config.Config
	fbGenerator   FizzBuzzGenerator
	statsRecorder FizzBuzzStatsRecorder

	log logger.Logger
}

// Controller interfaces for dependency injection
type FizzBuzzGenerator interface {
	GenerateFizzBuzz(req types.FizzBuzzRequest) (types.FizzBuzzResponse, error)
}

type FizzBuzzStatsRecorder interface {
	GetStats() types.FizzBuzzStats
	SaveStat(req types.FizzBuzzRequest) error
}

func NewFizzBuzzHandler(cfg *config.Config, log logger.Logger, generator FizzBuzzGenerator, statsRecorder FizzBuzzStatsRecorder) *FizzBuzzHandler {
	return &FizzBuzzHandler{
		cfg:           cfg,
		fbGenerator:   generator,
		log:           log,
		statsRecorder: statsRecorder,
	}
}

func (h *FizzBuzzHandler) GenerateFizzBuzz(c *gin.Context) {
	var req types.FizzBuzzRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Info("received FizzBuzz request", "request", req)

	result, err := h.fbGenerator.GenerateFizzBuzz(req)
	if err != nil {
		h.log.Error("failed to generate FizzBuzz", "error", err)
		if errors.Is(controllers.ErrLimitExceeded, err) || errors.Is(controllers.ErrStringLengthExceeded, err) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	err = h.statsRecorder.SaveStat(req)
	if err != nil {
		h.log.Error("failed to save stats", "error", err)
		// Proceed without failing the request
	}

	c.JSON(http.StatusOK, gin.H{
		"result":      result.Result,
		"duration_ms": result.Duration,
	})
}

func (h *FizzBuzzHandler) GetFizzBuzzStats(c *gin.Context) {
	stats := h.statsRecorder.GetStats()
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
