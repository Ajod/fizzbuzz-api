package controllers

import (
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStats(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})

	req1 := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}
	req2 := types.FizzBuzzRequest{Int1: 2, Int2: 4, Limit: 10, Str1: "Foo", Str2: "Bar"}
	recorder.SaveStat(req1)
	recorder.SaveStat(req1)
	recorder.SaveStat(req2)

	stats := recorder.GetStats()
	assert.Equal(2, stats.Count)
	assert.Contains(stats.MostFrequentRequests, req1)
}

func Test_GetStats_Empty(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	stats := recorder.GetStats()
	assert.Equal(0, stats.Count)
	assert.Empty(stats.MostFrequentRequests)
}

func Test_SerializeDeserializeRequest(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	req := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}

	serialized, err := recorder.serializeRequest(req)
	assert.NoError(err)

	deserialized := recorder.deserializeRequests([]string{serialized})
	assert.Len(deserialized, 1)
	assert.Equal(req, deserialized[0])
}

func Test_StatsExAequo(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	req1 := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}
	req2 := types.FizzBuzzRequest{Int1: 2, Int2: 4, Limit: 10, Str1: "Foo", Str2: "Bar"}
	req3 := types.FizzBuzzRequest{Int1: 1, Int2: 2, Limit: 5, Str1: "A", Str2: "B"}

	recorder.SaveStat(req1)
	recorder.SaveStat(req1)
	recorder.SaveStat(req2)
	recorder.SaveStat(req2)
	recorder.SaveStat(req3)
	stats := recorder.GetStats()

	assert.Equal(2, stats.Count)
	assert.Contains(stats.MostFrequentRequests, req1)
	assert.Contains(stats.MostFrequentRequests, req2)
	assert.NotContains(stats.MostFrequentRequests, req3)
}

func Test_SaveStatsLargeConcurrency(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	req := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}
	concurrency := 1000
	done := make(chan bool)
	for i := 0; i < concurrency; i++ {
		go func() {
			err := recorder.SaveStat(req)
			assert.NoError(err)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}
	stats := recorder.GetStats()
	assert.Equal(concurrency, stats.Count)
	assert.Len(stats.MostFrequentRequests, 1)
	assert.Equal(req, stats.MostFrequentRequests[0])
}

func Test_SaveAndGetStatsConcurrency(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	req := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}
	concurrency := 1000
	done := make(chan bool)
	for i := 0; i < concurrency; i++ {
		if i%10 == 0 {
			// Every 10th goroutine retrieves stats
			go func() {
				stats := recorder.GetStats()
				assert.GreaterOrEqual(stats.Count, 0)
				done <- true
			}()
			continue
		}
		go func() {
			err := recorder.SaveStat(req)
			assert.NoError(err)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// Test that reversed requests are not considered equal, and their counts are not aggregated
func Test_ReversedRequestsAreNotEqual(t *testing.T) {
	assert := assert.New(t)
	recorder := NewFizzBuzzStatsController(&mockLogger{})
	req1 := types.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"}
	req2 := types.FizzBuzzRequest{Int1: 5, Int2: 3, Limit: 15, Str1: "Buzz", Str2: "Fizz"}
	// Save both requests
	err := recorder.SaveStat(req1)
	assert.NoError(err)
	err = recorder.SaveStat(req2)
	assert.NoError(err)
	stats := recorder.GetStats()
	assert.Equal(1, stats.Count)
	assert.Len(stats.MostFrequentRequests, 2)
}
