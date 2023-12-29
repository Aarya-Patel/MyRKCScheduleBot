package env

import (
	"os"
)

type EnvConfig struct {
	// Bot Token granted by @BotFather
	BotToken string
	// Webhook domain that is used by Telegram API to send updates
	WebhookDomain string
	// Webhook secret token that is used by Telegram API when sending updates
	WebhookSecret string
	// GCP Project ID
	ProjectId string
}

var envConfig *EnvConfig

func GetEnvConfig() *EnvConfig {
	if envConfig == nil {
		// Get token from the environment variable.
		token := os.Getenv("TOKEN")
		if token == "" {
			panic("TOKEN environment variable is empty")
		}

		// Get the webhook domain from the environment variable.
		webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
		if webhookDomain == "" {
			panic("WEBHOOK_DOMAIN environment variable is empty")
		}
		// Get the webhook secret from the environment variable.
		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		if webhookSecret == "" {
			panic("WEBHOOK_SECRET environment variable is empty")
		}
		// Get the GCP project id from the environment variable.
		projectId := os.Getenv("PROJECT_ID")
		if projectId == "" {
			panic("PROJECT_ID environment variable is empty")
		}

		envConfig = &EnvConfig{
			BotToken:      token,
			WebhookDomain: webhookDomain,
			WebhookSecret: webhookSecret,
			ProjectId:     projectId,
		}
	}
	return envConfig
}
