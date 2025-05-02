package llmtools

import (
	"fmt"
	"time"

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

func CaptureMethodOfPayment(args map[string]any) string {
	fmt.Println("Capture method of payment args", args)
	return "successfull"
}

func CaptureCard(args map[string]any) string {
	fmt.Println("Capture card args", args)
	return "successfull"
}

func CaptureCVV(args map[string]any) string {
	fmt.Println("Capture CVV args", args)
	return "successfull"
}

func CaptureExpiry(args map[string]any) string {
	fmt.Println("Capture expiry args", args)
	expiryDate, ok := args["expiryDate"].(string)
	if !ok {
		return "invalid expiry date format"
	}

	// Parse the expiry date
	parsedDate, err := time.Parse("01/06", expiryDate)
	if err != nil {
		return "invalid expiry date format"
	}

	// Check if the expiry date is in the future
	if parsedDate.Before(time.Now()) {
		return "expiry date is invalid or expired"
	}
	return "successfull"
}

func ProcessPayments(args map[string]any) string {
	fmt.Println("Process payment args", args)
	return "successfull"
}

func PaymentConfirmation(args map[string]any) string {
	fmt.Println("Payment confirmation args", args)
	return "ABCD1234"
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
				Name:        "capture_method_of_payment",
				Description: "Capture the method of payment",
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
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "capture_card_number",
				Description: "Capture the credit or debit card number",
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "capture_cvv_number",
				Description: "Capture the credit or debit card's cvv number",
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "capture_expiry_date",
				Description: "Capture the credit or debit card's expiry date in MM/YY format",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"expiryDate": map[string]any{
							"type":        "string",
							"description": "Expiry date of the card in MM/YY format",
						},
					},
					"required": []string{"expiryDate"},
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
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "payment_confirmation",
				Description: "Confirms the payment and provides a confirmation code",
			},
		},
	}

	return availableTools
}
