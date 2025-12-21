# fizzbuzz-api

A small HTTP API that generates FizzBuzz sequences and exposes basic request statistics.

---

## Overview

- **Language:** Go
- **Framework:** Gin
- **Purpose:** Generate FizzBuzz-like sequences given input parameters and expose basic request statistics.
- **Entrypoint:** `cmd/fizzbuzz-api/main.go` which boots the server using `internal/fizzbuzzapi/http.NewServer()`.
- **Key packages:** `internal/fizzbuzzapi/handlers` (HTTP handlers), `controllers` (application services), `types` (shared types), `http` (server bootstrap).

This repository exposes:
- POST `/fizzbuzz/generate` — generate a FizzBuzz sequence and return the sequence with generation duration.
- GET `/fizzbuzz/stats` — return the most frequent request(s) recorded by the service and their counts.
- GET `/fizzbuzz/health` — basic health check, returns 200 OK with a simple body ("healthy")

---

## Quick start

Build and run locally:

```powershell
# set env as needed (optional)
$env:FBAPI_PORT="4255"; $env:FBAPI_HOST="localhost"
# run
go run ./cmd/fizzbuzz-api
```

Run tests:

```powershell
go test ./...
```

Configuration is loaded from environment variables with prefix `FBAPI_` (see `internal/fizzbuzzapi/config/config.go`):
- `FBAPI_PORT` (default `4255`)
- `FBAPI_HOST` (default `localhost`)
- `FBAPI_MAX_FIZZBUZZ_LIMIT` (default `100000`) — max allowed `limit` value
- `FBAPI_MAX_STRING_LENGTH` (default `30`) — max allowed length for `str1` / `str2`
- `FBAPI_STATS_STORAGE` (default `inmemory`) — storage type for stats (currently only `inmemory` is implemented)

---

## API Routes & Behavior

### POST /fizzbuzz/generate

- **Request JSON (all fields required):**

```json
{
  "int1": 3,
  "int2": 5,
  "limit": 15,
  "str1": "fizz",
  "str2": "buzz"
}
```

- **Success Response (200):**

```json
{
    "result": ["1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz"],
    "duration_ms": 1
}
```

- **Validation & Errors:**
  - `int1` and `int2` must be strictly positive integers (> 0). If not, the API returns `400 Bad Request`.
  - `limit` must be non-negative (>= 0); negative values result in `400 Bad Request`. A `limit` of `0` returns an empty sequence.
  - If `limit` exceeds the configured maximum (`FBAPI_MAX_FIZZBUZZ_LIMIT`), the API returns `422 Unprocessable Entity`.
  - If `str1` or `str2` exceeds `FBAPI_MAX_STRING_LENGTH`, the API returns `422 Unprocessable Entity`.
  - If the JSON cannot be bound, the API returns `400 Bad Request`.

- **Notes on behavior & performance:**
  - Generation is O(limit) in time and O(limit) in memory (returns the whole list). Large `limit` values can be CPU- and memory- intensive and produce large responses.
  - The controller logs generation duration (`duration_ms`) and returns it in the response.

### GET /fizzbuzz/stats

- **Success Response (200):**

```json
{
  "stats": {
    "most_frequent_request": [ { /* request object(s) */ } ],
    "count": 42
  }
}
```

- **Implementation details:**
  - Stats are recorded via a storage abstraction. The current implementation is an in-memory map (`map[string]int`). Other storage implementations are planned but not implemented.
  - You can configure `FBAPI_STATS_STORAGE` (e.g., `inmemory` or `file`) but currently only `inmemory` is implemented.
  - The API exposes the most frequent request(s) and the highest frequency count.
  - The API normalizes the parameters of the request such that int1 <= int2. This means that `"int1": 3, "int2": 5, "str1": "fizz", "str2": "buzz"` and `"int1": 5, "int2": 3, "str1": "buzz", "str2": "fizz"` are **equal**.

---

## Scope & Performance Trade-offs

- Memory:
  - Generating a large `limit` creates a large slice of strings and increases memory pressure. The configured `FBAPI_MAX_FIZZBUZZ_LIMIT` is a safety guard, but returning huge payloads still affects latency and network transfer costs.
  - The stats map is stored in memory and can grow with unique request payloads (unbounded unless controlled).

- Reliability & Scalability:
  - Stats are stored in-process: multiple instances will have inconsistent stats and restarting the process loses the data.
  - Saving stats is best-effort: if `SaveStat` errors, the request still succeeds (errors are logged but do not fail generation requests).

---

## Examples

Generate a fizzbuzz sequence (example via curl):

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"int1":3,"int2":5,"limit":15,"str1":"fizz","str2":"buzz"}' \
  http://localhost:4255/fizzbuzz/generate
```

Get stats:

```bash
curl http://localhost:4255/fizzbuzz/stats
```

---

## Future Improvements

- **Add a caching layer**
  - Cache recent request results keyed by the serialized request input so identical requests return the cached sequence and skip regeneration work.
  - Use TTLs and LRU strategies to bound cache memory.

- **Persist stats to a file or a database (Postgres, SQLite, etc.)**
  - Use a DB to store aggregated counts, time-series metrics, or raw events for long-term analytics.
  - Add background flushes or batch writes to reduce DB pressure.

- **Rate limiting & throttling**
  - Protect the service against excessive usage and DoS by implementing rate limits per client IP / API key.

- **Streaming / pagination**
  - For very large `limit` values, stream the sequence as a chunked response or offer pagination to reduce peak memory and transfer sizes.

- **Observability & metrics**
  - Add metrics such as request durations, error rates, request sizes, and structured logging to help monitor production behavior.

- **API refinement & docs**
  - Add OpenAPI/Swagger docs.

---
