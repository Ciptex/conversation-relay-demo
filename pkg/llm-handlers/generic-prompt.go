package llmhandlers

import "conversation-relay/pkg/types"

type GenericPromptHandler struct {
	accSid    string
	configSid string
	context   *HandlerContext
}

func newGenericPromptHandler(accSid, configSid string, context *HandlerContext) *GenericPromptHandler {
	return &GenericPromptHandler{
		accSid:    accSid,
		configSid: configSid,
		context:   context,
	}
}

func (h *GenericPromptHandler) Handle() (string, error) {
	// Fetch the prompt configuration from the repository
	promptConfig, err := h.context.Repo.GetPromptConfig(h.accSid, h.configSid)
	if err != nil {
		return "", err
	}
	config, _ := h.context.Repo.GetAccountConfig(h.accSid, h.configSid)
	prompt := newPrompt(h.context.CallSid, h.accSid, h.configSid, h.context.Repo, h.context.Span)
	parsedPrompt, err := prompt.getGenericPrompt(promptConfig.Config.OpenAI.GenericPrompt)
	h.context.Span.Debug("GenericPrompt::Handle transcipt len", "len", len(h.context.Transcript))
	// h.context.Span.Debug("GenericPrompt::Handle Parsed Generic Prompt: ", "parsedPrompt", parsedPrompt)
	model := h.context.LLM.New(types.LLMModelContext{Transcript: h.context.Transcript})
	response, err := model.CreateChatCompletion(config, h.context.CallSid, parsedPrompt, h.context.Span)
	if err != nil {
		return "", err
	}

	return response, nil
}
