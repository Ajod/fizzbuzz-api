package types

/*
	Defines the types used in the FizzBuzz API that are shared across different packages.
	Types that are only relevant to a specific package should be defined within that package.
*/

type FizzBuzzRequest struct {
	Int1  int    `json:"int1" binding:"required"`
	Int2  int    `json:"int2" binding:"required"`
	Limit int    `json:"limit" binding:"required"`
	Str1  string `json:"str1" binding:"required"`
	Str2  string `json:"str2" binding:"required"`
}

type FizzBuzzLimits struct {
	MaxLimit        int
	MaxStringLength int
}

type FizzBuzzResponse struct {
	Result   []string `json:"result"`
	Duration int64    `json:"duration_ms"`
}

type FizzBuzzStats struct {
	MostFrequentRequests []FizzBuzzRequest `json:"most_frequent_request"`
	Count                int               `json:"count"`
}
