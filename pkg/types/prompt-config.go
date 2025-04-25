package types

type OpenAIConfig struct {
	GenericPrompt string `yaml:"genericPrompt"`
}

type Config struct {
	OpenAI OpenAIConfig `yaml:"openAI"`
}

type PromptConfig struct {
	Version string `yaml:"version"`
	Config  Config `yaml:"config"`
}
