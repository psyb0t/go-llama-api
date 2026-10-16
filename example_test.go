package llama_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-llama-api"
)

// This is the README usage example. It talks to the live API, so it has no
// Output comment and only has to compile.
//
//nolint:testableexamples // needs a real API key, compile check only
func Example() {
	auth := llama.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			key := os.Getenv("LLAMA_API_KEY")
			req.Header.Set("Authorization", "Bearer "+key)

			return nil
		},
	)

	client, err := llama.NewClientWithResponses(
		"https://api.llama.com/v1",
		auth,
	)
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
