package controllers

import (
	"encoding/json"
	"fizzbuzz-api/internal/fizzbuzzapi/logger"
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"sync"
)

type StatsRecord map[string]int

type FizzBuzzStatsController struct {
	record StatsRecord
	log    logger.Logger

	sync.Mutex
}

func NewFizzBuzzStatsController(log logger.Logger) *FizzBuzzStatsController {
	return &FizzBuzzStatsController{
		record: make(StatsRecord),
		log:    log,
	}
}

func (ctrl *FizzBuzzStatsController) GetStats() types.FizzBuzzStats {
	ctrl.Lock()
	defer ctrl.Unlock()

	highestCount := 0
	var mostFrequentRequests []string
	for _, count := range ctrl.record {
		if count > highestCount {
			highestCount = count
		}
	}
	for req, count := range ctrl.record {
		if count == highestCount {
			mostFrequentRequests = append(mostFrequentRequests, req)
		}
	}

	return types.FizzBuzzStats{
		MostFrequentRequests: ctrl.deserializeRequests(mostFrequentRequests),
		Count:                highestCount,
	}
}

func (ctrl *FizzBuzzStatsController) SaveStat(req types.FizzBuzzRequest) error {
	ctrl.Lock()
	defer ctrl.Unlock()

	str, err := ctrl.serializeRequest(req)
	if err != nil {
		return err
	}
	ctrl.record[str]++
	ctrl.log.Info("stat recorded", "request", str, "new_count", ctrl.record[str])
	return nil
}

func (ctrl *FizzBuzzStatsController) serializeRequest(req types.FizzBuzzRequest) (string, error) {
	// Normalize the request parameters first to avoid different serializations for same logical requests
	// We organize Int1/Str1 and Int2/Str2 so that Int1 is always <= Int2
	if req.Int1 > req.Int2 {
		req.Int1, req.Int2 = req.Int2, req.Int1
		req.Str1, req.Str2 = req.Str2, req.Str1
	}

	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// str is a slice of serialized FizzBuzzRequests
func (ctrl *FizzBuzzStatsController) deserializeRequests(str []string) []types.FizzBuzzRequest {
	reqs := make([]types.FizzBuzzRequest, len(str))
	for i, s := range str {
		json.Unmarshal([]byte(s), &reqs[i])
	}
	return reqs
}
