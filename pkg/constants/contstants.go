package constants

import "os"

const (
	// twilio sends audio data as 160 byte messages containing 20ms of audio each
	// we will buffer 20 twilio messages corresponding to 0.4 seconds of audio to improve throughput performance
	BUFFER_SIZE       = 20 * 160
	MQ_PUB_PORT       = "8899"
	DEFAULT_SAY_VOICE = "Polly.Joanna"

	//TODO move this to account specific config
	MAX_TOKEN_COUNT = 5000
)

func ConfigTableConStr() string {
	return os.Getenv("CONFIG_TABLE_CON_STR")
}

func ConfigTableName() string {
	return os.Getenv("CONFIG_TABLE_NAME")
}

func YamlConfigfileDir() string {
	return os.Getenv("YAML_CONFIG_DIR")
}
