---
name: go-llama-api
description: Unofficial Go client for the Meta Llama API, generated from Meta's OpenAPI spec. `llama.NewClientWithResponses("https://api.llama.com/v1", auth)` from `github.com/psyb0t/go-llama-api` gives typed methods for chat completions (`CreateChatCompletionWithResponse`), models (`ListModelsWithResponse`, `RetrieveModelWithResponse`), moderations (`ModerationsWithResponse`) and file uploads (`PostUploadsWithResponse`, `GetUploadsUploadIdWithResponse`, `PostUploadsUploadIdWithResponse`), with responses decoded into `JSON200` and error structs. Auth is a Bearer API key set through a `RequestEditorFn`. Use when Go code needs to call the Llama API or already imports go-llama-api.
homepage: https://github.com/psyb0t/go-llama-api
user-invocable: true
metadata:
  openclaw:
    emoji: "🦙"
    requires:
      bins:
        - go
      env:
        - LLAMA_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the client sends requests only to the base URL you pass, `https://api.llama.com/v1`, over HTTPS."
---

# go-llama-api

Typed Go client for the Meta Llama API. Package `llama`, generated with oapi-codegen; every operation in the spec is a method on `ClientWithResponses`.

For install, the full operation list, unions and streaming, see [references/setup.md](references/setup.md).

## Security & safety

- Read the API key from the `LLAMA_API_KEY` environment variable. Never hardcode it, commit it, or log it, and never print request headers.
- Requests go only to the base URL you pass. Use `https://api.llama.com/v1`.
- Prompts and uploaded files are sent to Meta. Only send data the user meant to send, and only upload files when the user asked for it.

## When to use

- Go code needs Llama API chat completions, model listing, moderation or uploads.
- The project already imports `github.com/psyb0t/go-llama-api`.

## When NOT to use

- The code is not Go. Use Meta's official SDKs instead.
- You need a parsed streaming client. The generated client returns the raw `*http.Response` for `stream: true`; it does not parse server-sent events.
- You are calling an OpenAI-compatible endpoint with an OpenAI SDK; this client follows Meta's native spec.

## Usage

```go
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-llama-api"
)

auth := llama.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LLAMA_API_KEY"))

	return nil
})

client, err := llama.NewClientWithResponses("https://api.llama.com/v1", auth)
if err != nil {
	return err
}

models, err := client.ListModelsWithResponse(ctx)
if err != nil {
	return err // transport error
}

if models.JSON200 == nil {
	return fmt.Errorf("list models: %s", models.Status()) // API error, see models.Body
}

for _, model := range models.JSON200.Data {
	fmt.Println(model.Id)
}
```

Chat completion:

```go
var content llama.UserMessage_Content
if err := content.FromUserMessageContent0("hello"); err != nil {
	return err
}

var msg llama.Message
if err := msg.FromUserMessage(llama.UserMessage{
	Content: content,
	Role:    llama.UserMessageRoleUser,
}); err != nil {
	return err
}

res, err := client.CreateChatCompletionWithResponse(ctx, llama.CreateChatCompletionRequest{
	Model:    "<model-id>", // an Id from ListModels
	Messages: []llama.Message{msg},
})
if err != nil {
	return err
}

if res.JSON200 == nil {
	return fmt.Errorf("chat: %s", res.Status())
}

text, err := res.JSON200.CompletionMessage.Content.AsAssistantMessageContent0()
```

Check `JSON200` for nil before using it. On errors the client fills `JSON400`, `JSON429`, `JSON500` or `JSONDefault`, and the raw body is always in `Body`.
