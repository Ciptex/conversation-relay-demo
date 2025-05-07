# Conversation Relay Demo

A twilio conversation relay voice bot demo that guides callers through a workflow to collect account numbers, select payment methods, and complete payments. This application does not collect sensitive information such as credit card or bank account numbers.

For simplicity, the OpenAI tools validate and return static values. For example, the valid account number is `123` and the due amount is `1005 pounds`.

## Environment Configuration
Create a `.env` file and configure the environment variables listed below, or rename the `env-sample` file and fill in the relevant details.

- `PUBLIC_ENDPOINT`: Application endpoint (e.g., `12b4-171-76-82-37.ngrok-free.app`)
- `YAML_CONFIG_DIR`: Full directory path where prompt configuration files are stored (e.g., `/prompt/config/`)
- `PROMPT_CONFIG_FILE`: Configuration file name (e.g., `payment-bot.yml`)
- `AZURE_OPENAI_ENDPOINT`: Azure OpenAI endpoint (e.g., `https://ciptex.openai.azure.com/`)
- `AZURE_OPENAI_KEY`: Key generated in Azure portal
- `AZURE_OPENAI_MODEL`: Model deployed in Azure portal (e.g., `gpt-4o`)
- `AZURE_OPENAI_EMBEDDING_MODEL`: Embedding model deployed in Azure portal (e.g., `text-embedding-ada-002`)
- `AZURE_OPENAI_REGION`: Azure region (e.g., `eastus`)
- `TWILIO_ACCOUNT_SID`: Twilio account Sid
- `TWILIO_WORKFLOW_SID`: Twilio workflow Sid to create the task to transfer to the agant.
- `CARD_EASY_UNAME`: CardEasy username
- `CARD_EASY_PWD`: CardEasy password

> System use CardEasy gateway to collect card details securely

## Running the Application
Download and install golang from https://go.dev/dl/

1. Clone the repository from GitHub:

```bash
make build
make run
```

2. Open a different terminal tab and run:

```bash
make ngrok
```

> Note: Edit the Makefile to change the ngrok location if needed.

## Twilio Console Configuration

1. Open the Twilio console
2. Navigate to Phone Numbers → Manage → Active numbers
3. Click on an active phone number and go to the Configure tab
4. Modify "A call comes in" setting and choose "Webhook"
5. Set the URL to `https://[your-ngrok-url]/v1.0/configsid/twiml`
   - Note: The `configsid` parameter is not currently used in this demo project
6. Click "Save Configuration"

## Testing

Dial the configured phone number and follow the voice prompts.