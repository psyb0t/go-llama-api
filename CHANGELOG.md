# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `llama`, a REST client generated with oapi-codegen from the OpenAPI document Meta publishes for its Llama API SDKs (Llama API Specification, spec version v1, OpenAPI 3.1.0). It covers all 6 paths: chat completions, models, moderations and uploads.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` reads the current spec URL from the `.stats.yml` of meta-llama/llama-api-python, downloads the spec from Stainless and regenerates the client. The spec is not committed.
- Added the `go-llama-api` agent skill: `claude plugin install go-llama-api@psyb0t` (Claude Code), `codex plugin add go-llama-api@psyb0t` (Codex), `openclaw skills install @psyb0t/go-llama-api` (ClawHub).
