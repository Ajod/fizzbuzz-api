package controllers

import (
	"errors"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"strconv"
	"time"
)

type FizzBuzzConfig struct {
	MaxLimit        int
	MaxStringLength int
}

type FizzBuzzController struct {
	types.FizzBuzzLimits
	log logger.Logger
}

var (
	ErrLimitExceeded        = errors.New("limit exceeds maximum allowed")
	ErrStringLengthExceeded = errors.New("string length exceeds maximum allowed")
	ErrNegativeParameter    = errors.New("limit, int1, and int2 must be strictly positive integers")
)

func NewFizzBuzzController(limits types.FizzBuzzLimits, log logger.Logger) *FizzBuzzController {
	return &FizzBuzzController{
		FizzBuzzLimits: limits,
		log:            log,
	}
}

func (ctrl *FizzBuzzController) GenerateFizzBuzz(req types.FizzBuzzRequest) (types.FizzBuzzResponse, error) {
	if req.Limit < 0 || req.Int1 <= 0 || req.Int2 <= 0 {
		return types.FizzBuzzResponse{}, ErrNegativeParameter
	}

	if req.Limit > ctrl.MaxLimit {
		return types.FizzBuzzResponse{}, ErrLimitExceeded
	}

	if len(req.Str1) > ctrl.MaxStringLength || len(req.Str2) > ctrl.MaxStringLength {
		return types.FizzBuzzResponse{}, ErrStringLengthExceeded
	}

	start := time.Now()
	var result = make([]string, 0, req.Limit)
	for i := 1; i <= req.Limit; i++ {
		switch {
		case i%req.Int1 == 0 && i%req.Int2 == 0:
			result = append(result, req.Str1+req.Str2)
		case i%req.Int1 == 0:
			result = append(result, req.Str1)
		case i%req.Int2 == 0:
			result = append(result, req.Str2)
		default:
			result = append(result, strconv.Itoa(i))
		}
	}
	duration := time.Since(start)
	ctrl.log.Info("fizzBuzz generated", "limit", req.Limit, "duration_ms", duration.Milliseconds())

	return types.FizzBuzzResponse{
		Result:   result,
		Duration: duration.Milliseconds(),
	}, nil
}
