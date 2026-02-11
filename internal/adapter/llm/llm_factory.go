package llm

import (
	"context"
	"fmt"

	"github.com/shivangsaxena/inshorts-task/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

func NewLLM(ctx context.Context, cfg *config.Config) (llms.Model, error) {
	switch cfg.LLMProvider {
	case "openai":
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("openai api key is missing")
		}
		return openai.New(openai.WithToken(cfg.OpenAIKey))
	case "gemini":
		return googleai.New(ctx, googleai.WithAPIKey(cfg.GeminiKey))
	case "claude":
		return anthropic.New(anthropic.WithToken(cfg.ClaudeKey))
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", cfg.LLMProvider)
	}
}
