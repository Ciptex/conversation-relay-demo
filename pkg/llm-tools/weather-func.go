package llmtools

import "github.com/tmc/langchaingo/llms"

type WeatherToolParams struct {
	Location string `json:"location"`
}

func WeatherFunc(location string) string {
	// Here you would implement the logic to get the weather information
	// For now, we'll just return a dummy response
	return "The weather in " + location + " is sunny."
}

func WeatherFuncDefinition() []llms.Tool {
	var availableTools = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "getCurrentWeather",
				Description: "Get the current weather in a given location",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city and state, e.g. San Francisco, CA",
						},
						"unit": map[string]any{
							"type": "string",
							"enum": []string{"fahrenheit", "celsius"},
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}
	return availableTools
}
