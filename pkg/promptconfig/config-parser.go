package promptconfig

import (
	"conversation-relay/pkg/constants"
	"conversation-relay/pkg/types"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

func parseConfig(appConfig types.AccountConfig) (types.PromptConfig, error) {
	configFilePath := constants.YamlConfigfileDir()
	configFile := "default-config.yml"

	var config types.PromptConfig
	configuredFile := appConfig.PromptConfigFile
	if _, err := os.Stat(configFilePath + configuredFile); err == nil {
		// account specific config file exist use that config
		configFile = configuredFile
	}
	slog.Debug("promptConfig::parseConfig loading yaml config from", "configFilePath", configFilePath+configFile)
	yamlFile, err := os.ReadFile(configFilePath + configFile)
	if err != nil {
		slog.Error("promptConfig::parseConfig error parsing config file", "error", err)
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	return config, err
}

func LoadPromptConfig(appConfig types.AccountConfig) (types.PromptConfig, error) {
	config, err := parseConfig(appConfig)
	if err != nil {
		slog.Error("promptConfig::LoadConfig error parsing config file", "error", err)
		return config, err
	}
	// slog.Debug("promptConfig::LoadConfig loaded yaml config", "config", config)
	return config, nil
}
