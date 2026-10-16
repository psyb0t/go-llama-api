# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `llama_test`, and use only the exported API.
- **TST2** Tests talk to an `httptest` server and check the method, path, `Authorization: Bearer` header and request body, then the decoded response. They never call the live API.
- **TST3** `example_test.go` holds the README usage example. When the README example changes, change it too; `make test` compiles it.
- **TST4** `make test` runs `go test -race` in the dev image.
