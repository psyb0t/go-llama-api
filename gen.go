// Package llama is a Go client for the Meta Llama API, generated with
// oapi-codegen from the OpenAPI document Meta publishes.
package llama

//go:generate bash scripts/fetch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.yaml
