# AGENTS.md

This repository welcomes contributions from AI agents.

## Project Overview

This is a Go API wrapper for the EODHD (End of Day Historical Data) financial API. The codebase provides HTTP client functionality with rate limiting, retry logic, and support for JSON/CSV response formats.

## Repository Structure

- `eodhd.go` — Core HTTP client implementation with rate limiting
- `eodhd_options.go` — Client configuration options
- `bulk_eod.go` — Bulk end-of-day data endpoint
- `ohlcv.go` — OHLCV historical data endpoint
- `tickers.go` — Ticker listing endpoint
- `exchanges.go` — Exchange listing endpoint
- `params.go` — Params interface
- `util.go` — Utility helpers
- `cmd/examples/` — Usage examples

## Running Tests

```bash
go test ./...
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use GoDoc comments on all exported types and functions
- Keep functions focused and small

## Git Commits

- One-line commit messages only (no body, no additional comments)
- No attribution lines (no Co-Authored-By, no Signed-off-by)
- Imperative mood (e.g., "add rate limiter", "fix validation bug")

## Areas That Need Work

- Test coverage is at 38% — many exported functions are untested
- Most exported functions lack GoDoc comments
- README is minimal and needs expansion
- Input validation is missing in several places
- There is a known bug in `ohlcv.go` where the `To` field has an incorrect URL tag
