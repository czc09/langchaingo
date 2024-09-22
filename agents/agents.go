package agents

import (
	"context"
	"github.com/czc09/langchaingo/chains"

	"github.com/czc09/langchaingo/schema"
	"github.com/czc09/langchaingo/tools"
)

// Agent is the interface all agents must implement.
type Agent interface {
	// Plan Given an input and previous steps decide what to do next. Returns
	// either actions or a finish.
	Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string, opts ...chains.ChainCallOption) ([]schema.AgentAction, *schema.AgentFinish, error) //nolint:lll
	GetInputKeys() []string
	GetOutputKeys() []string
	GetTools() []tools.Tool
}
