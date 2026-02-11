package port

import "context"

type LLMResponse struct {
	Intent   string   `json:"intent"`
	Entities []string `json:"entities"`
	Location string   `json:"location,omitempty"`
	Category string   `json:"category,omitempty"`
}

type LLMService interface {
	ParseQuery(ctx context.Context, query string) (*LLMResponse, error)
	Summarize(ctx context.Context, text string) (string, error)
}
