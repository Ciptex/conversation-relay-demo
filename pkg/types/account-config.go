package types

type AccountConfig struct {
	AccountSid                         string
	ConfigID                           string
	PromptConfigFile                   string
	AzureOpenAIKey                     string
	AzureOpenAIEndpoint                string
	AzureOpenAIRegion                  string
	AzureOpenAIDeploymentName          string
	AzureOpenAIEmbeddingDeploymentName string
	TwilioWorkFlowSid                  string

	// TwilioApiKey                       string
	// TwilioApiSecret                    string
	// TwilioFlowSid                      string
}
