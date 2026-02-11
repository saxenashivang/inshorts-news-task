package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shivangsaxena/inshorts-task/internal/core/port"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

type LangChainService struct {
	llm llms.Model
}

func NewLangChainService(llm llms.Model) *LangChainService {
	return &LangChainService{llm: llm}
}

func (s *LangChainService) ParseQuery(ctx context.Context, query string) (*port.LLMResponse, error) {
	// 1. Define Prompt Template
	template := `
Analyze the user query: "{{.query}}"
Return a JSON object with:
- "intent": "nearby", "search", or "category"
- "location": if a place is mentioned
- "category": if a news category is mentioned (e.g., tech, sports)
- "entities": list of keywords
Response must be valid JSON only.
`
	prompt := prompts.NewPromptTemplate(template, []string{"query"})

	// 2. Create Chain
	chain := chains.NewLLMChain(s.llm, prompt)

	// 3. Execute
	res, err := chains.Call(ctx, chain, map[string]any{"query": query})
	if err != nil {
		return nil, err
	}

	// 4. Parse JSON (Manual for now, can use OutputParsers)
	// LangChainGo chains return map[string]any usually with "text" key if using simple LLMChain
	output, ok := res["text"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected chain output format")
	}

	// Clean markdown code blocks if present
	output = strings.TrimPrefix(output, "```json")
	output = strings.TrimPrefix(output, "```")
	output = strings.TrimSuffix(output, "```")
	output = strings.TrimSpace(output)

	var parsed port.LLMResponse
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse llm json: %w. Output: %s", err, output)
	}
	return &parsed, nil
}

func (s *LangChainService) Summarize(ctx context.Context, text string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, s.llm, "Summarize this article in 1 sentence: "+text)
	if err != nil {
		return "", err
	}
	return completion, nil
}
