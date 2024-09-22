package agents

import (
	"context"
	_ "embed"
	"regexp"
	"strings"

	"github.com/czc09/langchaingo/callbacks"
	"github.com/czc09/langchaingo/chains"
	"github.com/czc09/langchaingo/llms"
	"github.com/czc09/langchaingo/prompts"
	"github.com/czc09/langchaingo/schema"
	"github.com/czc09/langchaingo/tools"
)

const (
	_conversationalFinalAnswerAction = "AI:"
)

// ConversationalAgent is a struct that represents an agent responsible for deciding
// what to do or give the final output if the task is finished given a set of inputs
// and previous steps taken.
//
// Other agents are often optimized for using tools to figure out the best response,
// which is not ideal in a conversational setting where you may want the agent to be
// able to chat with the user as well.
type ConversationalAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain chains.Chain
	// Tools is a list of the tools the agent can use.
	Tools []tools.Tool
	// Output key is the key where the final output is placed.
	OutputKey string
	// CallbacksHandler is the handler for callbacks.
	CallbacksHandler callbacks.Handler
}

var _ Agent = (*ConversationalAgent)(nil)

func NewConversationalAgent(llm llms.Model, tools []tools.Tool, opts ...Option) *ConversationalAgent {
	options := conversationalDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &ConversationalAgent{
		Chain: chains.NewLLMChain(
			llm,
			options.getConversationalPrompt(tools),
			chains.WithCallback(options.callbacksHandler),
		),
		Tools:            tools,
		OutputKey:        options.outputKey,
		CallbacksHandler: options.callbacksHandler,
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *ConversationalAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
	opts ...chains.ChainCallOption,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if a.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			a.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}
	options := make([]chains.ChainCallOption, 0, len(opts)+1)
	options = append(options, opts...)
	options = append(options, chains.WithStreamingFunc(stream))
	options = append(options, chains.WithStopWords([]string{"\nObservation:", "\n\tObservation:"}))
	output, err := chains.Predict(
		ctx,
		a.Chain,
		fullInputs,
		options...,
	)
	if err != nil {
		return nil, nil, err
	}

	return a.parseOutput(output)
}

func (a *ConversationalAgent) GetInputKeys() []string {
	chainInputs := a.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *ConversationalAgent) GetOutputKeys() []string {
	return []string{a.OutputKey}
}

func (a *ConversationalAgent) GetTools() []tools.Tool {
	return a.Tools
}

func constructScratchPad(steps []schema.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}

func (a *ConversationalAgent) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, _conversationalFinalAnswerAction) {
		splits := strings.Split(output, _conversationalFinalAnswerAction)

		finishAction := &schema.AgentFinish{
			ReturnValues: map[string]any{
				a.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}

		return nil, finishAction, nil
	}

	r := regexp.MustCompile(`Action: (.*?)[\n]*Action Input: (.*)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) != 0 {
		return []schema.AgentAction{
			{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
		}, nil, nil
	}

	// Compatible output information without any prefix
	finishAction := &schema.AgentFinish{
		ReturnValues: map[string]any{
			a.OutputKey: output,
		},
		Log: output,
	}
	return nil, finishAction, nil
}

//go:embed prompts/conversational_prefix.txt
var _defaultConversationalPrefix string //nolint:gochecknoglobals

//go:embed prompts/conversational_format_instructions.txt
var _defaultConversationalFormatInstructions string //nolint:gochecknoglobals

//go:embed prompts/conversational_suffix.txt
var _defaultConversationalSuffix string //nolint:gochecknoglobals

func createConversationalPrompt(tools []tools.Tool, prefix, instructions, suffix string) prompts.PromptTemplate {
	template := strings.Join([]string{prefix, instructions, suffix}, "\n\n")

	return prompts.PromptTemplate{
		Template:       template,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad"},
		PartialVariables: map[string]any{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
			"history":           "",
		},
	}
}
