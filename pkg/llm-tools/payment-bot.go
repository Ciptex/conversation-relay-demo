package llmtools

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func ValidateAccount(args map[string]any) string {
	if args["accNum"] == "123" {
		return "valid"
	}
	return "invalid"
}

func GetAccBalance(args map[string]any) string {
	return "1005"
}

func ProcessPayments(args map[string]any) string {
	fmt.Println("Process payment args", args)
	return "successfull"
}

func PaymentBotDef() []llms.Tool {
	var availableTools = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "validate_account",
				Description: "Validates if an account number exists in the system",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"accNum": map[string]any{
							"type":        "string",
							"description": "The account number provided by the caller",
						},
					},
					"required": []string{"accNum"},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "get_account_balance",
				Description: "Retrieves the current balance due for a valid account",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"accNum": map[string]any{
							"type":        "string",
							"description": "The account number provided by the caller",
						},
					},
					"required": []string{"accNum"},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "process_payment",
				Description: "Processes a payment for the specified account",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"accNum": map[string]any{
							"type":        "string",
							"description": "The account number provided by the caller",
						},
						"paymentMethod": map[string]any{
							"type":        "string",
							"enum":        []string{"credit_card", "bank_transfer", "debit_card"},
							"description": "The payment method selected by the caller",
						},
						"amount": map[string]any{
							"type":        "number",
							"description": "The payment amount",
						},
					},

					"required": []string{"accNum"},
				},
			},
		},
	}

	return availableTools
}
