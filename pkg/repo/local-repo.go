package repo

import (
	"conversation-relay/pkg/promptconfig"
	"conversation-relay/pkg/types"
	"log/slog"
	"os"
)

type LocalRepo struct {
	callContext       map[string][]types.MessageTranscript
	paymentMeta       map[string]types.PaymentMeta
	lastMessage       map[string]string
	localAccConfig    types.AccountConfig
	localPromptConfig types.PromptConfig
}

func newLocalRepo() *LocalRepo {
	return &LocalRepo{
		callContext: make(map[string][]types.MessageTranscript),
		lastMessage: make(map[string]string),
		paymentMeta: make(map[string]types.PaymentMeta),
	}
}

func (l *LocalRepo) SetPaymentMeta(callSid string, paymentMeta types.PaymentMeta) {
	l.paymentMeta[callSid] = paymentMeta
}
func (l *LocalRepo) GetPaymentMeta(callSid string) types.PaymentMeta {
	if paymentMeta, ok := l.paymentMeta[callSid]; ok {
		return paymentMeta
	}
	return types.PaymentMeta{}
}

func (l *LocalRepo) SetAccountConfig(accConfig types.AccountConfig) {
	// No-op for local repo
}
func (l *LocalRepo) GetAccountConfig(accountSid, configId string) (types.AccountConfig, error) {
	if l.localAccConfig.AccountSid != "" {
		return l.localAccConfig, nil
	}
	// localAccConfig = azure.GetConfig(accountSid, configId)
	l.localAccConfig = types.AccountConfig{
		AccountSid:                         os.Getenv("TWILIO_ACCOUNT_SID"),
		AzureOpenAIKey:                     os.Getenv("AZURE_OPENAI_KEY"),
		AzureOpenAIEndpoint:                os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIDeploymentName:          os.Getenv("AZURE_OPENAI_MODEL"),
		AzureOpenAIEmbeddingDeploymentName: os.Getenv("AZURE_OPENAI_EMBEDDING_MODEL"),
		AzureOpenAIRegion:                  os.Getenv("AZURE_OPENAI_REGION"),
		PromptConfigFile:                   os.Getenv("PROMPT_CONFIG_FILE"),
		// TwilioApiKey:                       os.Getenv("TWILIO_API_KEY"),
		// TwilioApiSecret:                    os.Getenv("TWILIO_API_SECRET"),
		// TwilioFlowSid:                      os.Getenv("TWILIO_STUDIO_FLOW_SID"),
		TwilioWorkFlowSid: os.Getenv("TWILIO_WORKFLOW_SID"),
	}
	slog.Info("Loaded account config", "accountSid", l.localAccConfig.AccountSid, "deploymentName", l.localAccConfig.AzureOpenAIDeploymentName)
	return l.localAccConfig, nil
}

func (l *LocalRepo) GetPromptConfig(accountSid, configId string) (types.PromptConfig, error) {
	if l.localPromptConfig.Version != "" {
		return l.localPromptConfig, nil
	}
	localPromptConfig, err := promptconfig.LoadPromptConfig(l.localAccConfig)
	l.localPromptConfig = localPromptConfig
	return localPromptConfig, err
}

func (l *LocalRepo) AddCallContext(callSid, role, message string) {
	currentMsgs := l.callContext[callSid]
	newMsg := types.MessageTranscript{
		Role:    role,
		Message: message,
	}
	currentMsgs = append(currentMsgs, newMsg)
	l.callContext[callSid] = currentMsgs
	l.lastMessage[callSid] = role + ": " + message
}

func (l *LocalRepo) ResetCallContext(callSid string) {
	delete(l.callContext, callSid)
}

func (l *LocalRepo) GetCallContext(callSid string) []types.MessageTranscript {
	return l.callContext[callSid]
}

func (l *LocalRepo) GetLastMessage(callSid string) string {
	return l.lastMessage[callSid]
}
