package repo

import "conversation-relay/pkg/types"

type IRepo interface {
	SetAccountConfig(accConfig types.AccountConfig)
	GetAccountConfig(accountSid, configId string) (types.AccountConfig, error)
	GetPromptConfig(accountSid, configId string) (types.PromptConfig, error)
	AddCallContext(callSid, role, message string)
	GetCallContext(callSid string) []types.MessageTranscript
	GetLastMessage(callSid string) string
	ResetCallContext(callSid string)
}

func NewDB() IRepo {
	return newLocalRepo()
}
