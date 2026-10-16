package llama_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-llama-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *llama.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := llama.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := llama.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBearerAuth(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
}

func TestListModels(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/models", r.URL.Path)
		assertBearerAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"object":"list","data":[{"id":"model-a","created":1700000000,` +
				`"object":"model","owned_by":"owner-x"}]}`,
		))
	})

	res, err := client.ListModelsWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Data, 1)

	model := res.JSON200.Data[0]
	assert.Equal(t, "model-a", model.Id)
	assert.Equal(t, 1700000000, model.Created)
	assert.Equal(t, llama.ModelObjectModel, model.Object)
	assert.Equal(t, "owner-x", model.OwnedBy)
}

func TestCreateChatCompletion(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearerAuth(t, r)

		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Temperature float64 `json:"temperature"`
		}

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "model-a", body.Model)
		assert.InDelta(t, 0.25, body.Temperature, 0)
		assert.Len(t, body.Messages, 1)

		if len(body.Messages) == 1 {
			assert.Equal(t, "user", body.Messages[0].Role)
			assert.Equal(t, "hello", body.Messages[0].Content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"id":"cmpl-1","completion_message":{"role":"assistant",` +
				`"content":"hi there","stop_reason":"stop"}}`,
		))
	})

	var content llama.UserMessage_Content
	require.NoError(t, content.FromUserMessageContent0("hello"))

	var msg llama.Message
	require.NoError(t, msg.FromUserMessage(llama.UserMessage{
		Content: content,
		Role:    llama.UserMessageRoleUser,
	}))

	temperature := 0.25

	res, err := client.CreateChatCompletionWithResponse(
		context.Background(),
		llama.CreateChatCompletionRequest{
			Model:       "model-a",
			Messages:    []llama.Message{msg},
			Temperature: &temperature,
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, "cmpl-1", *res.JSON200.Id)

	reply := res.JSON200.CompletionMessage
	assert.Equal(t, llama.AssistantMessageRoleAssistant, reply.Role)
	require.NotNil(t, reply.Content)

	text, err := reply.Content.AsAssistantMessageContent0()
	require.NoError(t, err)
	assert.Equal(t, "hi there", text)
}
