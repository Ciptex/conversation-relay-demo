package repo

import (
	"conversation-relay/pkg/promptconfig"
	"conversation-relay/pkg/types"
	"log/slog"
	"os"
)

type LocalRepo struct {
	callContext map[string][]types.MessageTranscript
	lastMessage map[string]string
}

func newLocalRepo() *LocalRepo {
	return &LocalRepo{
		callContext: make(map[string][]types.MessageTranscript),
		lastMessage: make(map[string]string),
	}
}

var localAccConfig types.AccountConfig
var localPromptConfig types.PromptConfig

func (l *LocalRepo) SetAccountConfig(accConfig types.AccountConfig) {
	// No-op for local repo
}
func (l *LocalRepo) GetAccountConfig(accountSid, configId string) (types.AccountConfig, error) {
	if localAccConfig.AccountSid != "" {
		return localAccConfig, nil
	}
	// localAccConfig = azure.GetConfig(accountSid, configId)
	localAccConfig = types.AccountConfig{
		AccountSid:                         os.Getenv("TWILIO_ACCOUNT_SID"),
		AzureOpenAIKey:                     os.Getenv("AZURE_OPENAI_KEY"),
		AzureOpenAIEndpoint:                os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIDeploymentName:          os.Getenv("AZURE_OPENAI_MODEL"),
		AzureOpenAIEmbeddingDeploymentName: os.Getenv("AZURE_OPENAI_EMBEDDING_MODEL"),
		AzureOpenAIRegion:                  os.Getenv("AZURE_OPENAI_REGION"),
		PromptConfigFile:                   os.Getenv("PROMPT_CONFIG_FILE"),
	}
	slog.Info("Loaded account config", "accountSid", localAccConfig.AccountSid, "deploymentName", localAccConfig.AzureOpenAIDeploymentName)
	return localAccConfig, nil
}

func (l *LocalRepo) GetPromptConfig(accountSid, configId string) (types.PromptConfig, error) {
	if localPromptConfig.Version != "" {
		return localPromptConfig, nil
	}
	localPromptConfig, err := promptconfig.LoadPromptConfig(localAccConfig)
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
