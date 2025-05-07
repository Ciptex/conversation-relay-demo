package llms

import (
	"context"
	llmtools "conversation-relay/pkg/llm-tools"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"encoding/json"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type GPTLLM struct{}

func newGTPLLM() ILLM {
	return GPTLLM{}
}

func (l GPTLLM) New(context types.LLMModelContext) LLMModel {
	return &GPTModel{
		Transcript: context.Transcript,
	}
}

type GPTModel struct {
	Transcript []types.MessageTranscript
}

func (g *GPTModel) CreateChatCompletion(config types.AccountConfig, sid, prompt string, span trace.ISpan) (string, error) {
	llm, err := g.createClient(config)
	if err != nil {
		span.Error("langchain::CreateChatCompletion::error creating openai client", err)
		return "", err
	}
	ctx := context.Background()

	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, prompt),
	}

	for i := 0; i < len(g.Transcript); i++ {
		messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageType(g.Transcript[i].Role), g.Transcript[i].Message))
	}

	resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(llmtools.PaymentBotDef()), llms.WithTemperature(0.0))

	messageHistory, currResp := g.executeToolCalls(ctx, sid, llm, messageHistory, resp, span)
	return currResp, err
}

func (g *GPTModel) CreateEmbedding(config types.AccountConfig, text string, span trace.ISpan) ([]float32, error) {
	if text == "" {
		return []float32{0}, errors.New("Unable to create embedding. Text is empty")
	}
	llm, err := g.createClient(config)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	embedings, err := llm.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return embedings[0], nil
}

func (g *GPTModel) createClient(config types.AccountConfig) (*openai.LLM, error) {
	var opts []openai.Option
	opts = append(opts, openai.WithBaseURL(config.AzureOpenAIEndpoint))
	opts = append(opts, openai.WithToken(config.AzureOpenAIKey))
	opts = append(opts, openai.WithModel(config.AzureOpenAIDeploymentName))
	opts = append(opts, openai.WithAPIType(openai.APITypeAzure))
	opts = append(opts, openai.WithEmbeddingModel(config.AzureOpenAIEmbeddingDeploymentName))
	llm, err := openai.New(opts...)
	return llm, err
}

// TODO remove hardcoded tool names and use plugin approach
func (g *GPTModel) executeToolCalls(ctx context.Context, sid string, llm *openai.LLM, messageHistory []llms.MessageContent, resp *llms.ContentResponse, span trace.ISpan) ([]llms.MessageContent, string) {
	span.Debug("GPTModel::Executing tools", "toolsLen", len(resp.Choices[0].ToolCalls))
	if len(resp.Choices[0].ToolCalls) <= 0 {
		return messageHistory, resp.Choices[0].Content
	}
	var currResp string = ""
	for _, toolCall := range resp.Choices[0].ToolCalls {
		span.Debug("GPTModel::executeToolCalls", "tool", toolCall.FunctionCall.Name, "args", toolCall.FunctionCall.Arguments)
		switch toolCall.FunctionCall.Name {
		case "validate_account":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls validate_account", err)
			}
			args["callSid"] = sid
			resp := llmtools.ValidateAccount(args)
			currResp = resp

		case "get_account_balance":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls get_account_balance", err)
			}
			args["callSid"] = sid
			resp := llmtools.GetAccBalance(args)
			currResp = resp
		case "capture_method_of_payment":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls capture_method_of_payment", err)
			}
			args["callSid"] = sid
			resp := llmtools.CaptureMethodOfPayment(args)
			currResp = resp
		case "capture_card_number":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls capture_card_number", err)
			}
			args["callSid"] = sid
			resp := llmtools.CaptureCard(args)
			currResp = resp
		case "capture_cvv_number":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls capture_cvv_number", err)
			}
			args["callSid"] = sid
			resp := llmtools.CaptureCVV(args)
			currResp = resp
		case "capture_expiry_date":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls capture_expiry_date", err)
			}
			args["callSid"] = sid
			resp := llmtools.CaptureExpiry(args)
			currResp = resp
		case "process_payment":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls process_payment", err)
			}
			args["callSid"] = sid
			resp := llmtools.ProcessPayments(args)
			currResp = resp
		case "payment_confirmation":
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls payment_confirmation", err)
			}
			args["callSid"] = sid
			resp := llmtools.PaymentConfirmation(args)
			currResp = resp
		case "getCurrentWeather":
			var args struct {
				Location string `json:"location"`
				Unit     string `json:"unit"`
			}
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				span.Error("GPTModel::executeToolCalls getCurrentWeather", err)
			}

			response := llmtools.WeatherFunc(args.Location)
			currResp = response

		default:
			span.Warn("GPTModel::executeToolCalls Unsupported tool: %s", toolCall.FunctionCall.Name)
		}
		toolCallIn := llms.MessageContent{
			Role: llms.ChatMessageTypeAI,
			Parts: []llms.ContentPart{
				llms.ToolCall{
					ID:   toolCall.ID,
					Type: "function",
					FunctionCall: &llms.FunctionCall{
						Name:      toolCall.FunctionCall.Name,
						Arguments: toolCall.FunctionCall.Arguments,
					},
				},
			},
		}
		messageHistory = append(messageHistory, toolCallIn)

		toolsResp := llms.MessageContent{
			Role: llms.ChatMessageTypeTool,
			Parts: []llms.ContentPart{
				llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    currResp,
				},
			},
		}
		messageHistory = append(messageHistory, toolsResp)
	}
	finalResp := currResp
	_ = finalResp
	ct, err := llm.GenerateContent(ctx, messageHistory)
	span.Debug("GPTModel::executeToolCalls final resp", "choicesLen", len(ct.Choices), "choiceContent", ct.Choices[0].Content, "err", err)
	if err != nil {
		span.Error("GPTModel::executeToolCalls error in final response", err)
		return messageHistory, finalResp
	}

	return g.executeToolCalls(ctx, sid, llm, messageHistory, ct, span)
}
