# Go code

Refines `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` is generated. Never edit it by hand. Change `config.yaml` or the spec handling in `scripts/` and run `make generate`.
- **GO2** Spec fixes happen before oapi-codegen runs, as a step in `gen.go` that writes a patched copy under the gitignored `api/` directory. The downloaded spec is never committed or edited by hand.
- **GO3** `scripts/fetch-spec.sh` reads `openapi_spec_url` from the `.stats.yml` of meta-llama/llama-api-python and accepts only a URL under the Stainless spec bucket. Keep that check.
- **GO4** `number` maps to `float64` in `config.yaml`. Keep it.
- **GO5** `make lint` passes with zero issues. Lines are at most 80 columns in hand-written code; split calls one argument per line rather than adding `//nolint`.
