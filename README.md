# go-llama-api

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-llama-api.svg)](https://pkg.go.dev/github.com/psyb0t/go-llama-api)
[![CI](https://github.com/psyb0t/go-llama-api/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-llama-api/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-llama-api/badges/version.svg)](https://github.com/psyb0t/go-llama-api/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-llama-api/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-llama-api/badges/importers.svg)](https://github.com/psyb0t/go-llama-api/blob/badges/importers.md)

Unofficial Go client for the [Meta Llama API](https://llama.developer.meta.com), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Meta publishes for its SDKs. I am not affiliated with Meta.

The package is `llama`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-llama-api
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-llama-api"
)

func main() {
	auth := llama.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("LLAMA_API_KEY"))

		return nil
	})

	client, err := llama.NewClientWithResponses("https://api.llama.com/v1", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListModelsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	if res.JSON200 == nil {
		panic(res.Status())
	}

	for _, model := range res.JSON200.Data {
		fmt.Println(model.Id)
	}
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-llama-api) and in `go doc github.com/psyb0t/go-llama-api`.

Chat messages are unions. Build one with the `From...` helpers, for example `UserMessage_Content.FromUserMessageContent0("hello")` for plain text and then `Message.FromUserMessage(...)`, which also sets the `role`. Read the reply with `CompletionMessage.Content.AsAssistantMessageContent0()`.

The chat endpoint can also stream server-sent events when `stream` is true. The generated client does not parse event streams; use `CreateChatCompletion` (without `WithResponse`) and read the `*http.Response` body yourself.

## Auth

The Llama API takes an API key as a Bearer token in the `Authorization` header. Keys come from the API keys page of the Llama developer dashboard. The spec names a security requirement but does not define the scheme, so the client adds no auth by itself. Pass the request editor above to every client.

The base URL is `https://api.llama.com/v1`. The spec's `servers` entry lists `http://api.llama.com` without the version prefix, so always pass the base URL yourself.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` reads `openapi_spec_url` from the `.stats.yml` of [meta-llama/llama-api-python](https://github.com/meta-llama/llama-api-python) and downloads that URL to `api/openapi.yaml`. The spec is hosted by Stainless under a URL that changes with every update, so the script looks it up each time.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64`, because the spec gives temperatures and scores no format and oapi-codegen would otherwise use `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-llama-api) teaches a coding agent how to use go-llama-api when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-llama-api@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-llama-api@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-llama-api
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from the API specification Meta publishes for its SDKs. The specification itself is not redistributed here.
