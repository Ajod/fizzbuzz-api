package controllers

import (
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Error(msg string, args ...any) {}
func (m *mockLogger) Debug(msg string, args ...any) {}

func Test_GenerateFizzBuzz_Valid(t *testing.T) {
	assert := assert.New(t)

	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        100,
		MaxStringLength: 100,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "Fizz",
		Str2:  "Buzz",
	}

	expected := []string{
		"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz",
		"11", "Fizz", "13", "14", "FizzBuzz",
	}
	resp, err := ctrl.GenerateFizzBuzz(req)

	assert.NoError(err, "Error should be nil")
	assert.Equal(len(expected), len(resp.Result), "Result length should match expected length")
	for i := range expected {
		assert.Equal(expected[i], resp.Result[i], "Result element should match expected element at index %d", i)
	}
}

func Test_GenerateFizzBuzz_NegLimit(t *testing.T) {
	assert := assert.New(t)
	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        100,
		MaxStringLength: 100,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: -1,
		Str1:  "Fizz",
		Str2:  "Buzz",
	}

	resp, err := ctrl.GenerateFizzBuzz(req)
	assert.Error(err, "Error should not be nil for negative limit")
	assert.Equal(0, len(resp.Result), "Result should be empty for negative limit")
}

func Test_GenerateFizzBuzz_ZeroInt1(t *testing.T) {
	assert := assert.New(t)
	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        100,
		MaxStringLength: 100,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  0,
		Int2:  5,
		Limit: 15,
		Str1:  "Fizz",
		Str2:  "Buzz",
	}
	resp, err := ctrl.GenerateFizzBuzz(req)
	assert.Error(err, "Error should not be nil for zero Int1")
	assert.Equal(0, len(resp.Result), "Result should be empty for zero Int1")
}

func Test_GenerateFizzBuzz_StrOverLimit(t *testing.T) {
	assert := assert.New(t)
	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        100,
		MaxStringLength: 10,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "I am longer than ten characters",
		Str2:  "I'm not",
	}
	resp, err := ctrl.GenerateFizzBuzz(req)
	assert.Error(err, "Error should not be nil for string exceeding max length")
	assert.Equal(ErrStringLengthExceeded, err, "Error should be ErrStringLengthExceeded")
	assert.Equal(0, len(resp.Result), "Result should be empty for string exceeding max length")
}

func Test_GenerateFizzBuzz_LimitExceeded(t *testing.T) {
	assert := assert.New(t)
	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        50,
		MaxStringLength: 100,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "Fizz",
		Str2:  "Buzz",
	}
	resp, err := ctrl.GenerateFizzBuzz(req)
	assert.Error(err, "Error should not be nil for limit exceeding max limit")
	assert.Equal(ErrLimitExceeded, err, "Error should be ErrLimitExceeded")
	assert.Equal(0, len(resp.Result), "Result should be empty for limit exceeding max limit")
}

func Test_GenerateFizzBuzz_ZeroLimit(t *testing.T) {
	assert := assert.New(t)
	ctrl := NewFizzBuzzController(types.FizzBuzzLimits{
		MaxLimit:        100,
		MaxStringLength: 100,
	}, &mockLogger{})
	req := types.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 0,
		Str1:  "Fizz",
		Str2:  "Buzz",
	}
	resp, err := ctrl.GenerateFizzBuzz(req)
	assert.NoError(err, "Error should be nil for zero limit")
	assert.Equal(0, len(resp.Result), "Result should be empty for zero limit")
}
