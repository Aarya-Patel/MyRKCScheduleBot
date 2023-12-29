package main

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/env"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/handler"
)

// This bot repeats everything you say - but it uses webhooks instead of long polling.
// Webhooks are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
// For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
// Simply install ngrok, make an account on the website, and run:
// `ngrok http 8080`
// Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
// from the samples/echoWebhookBot directory:
// `TOKEN="<your_token_here>" WEBHOOK_DOMAIN="<your_domain_here>"  WEBHOOK_SECRET="<random_string_here>" go run .`
// Then, simply send /start to your bot; if it replies, you've successfully set up webhooks!

const WEBHOOK_PATH = "custom-path/"

func main() {
	envConfig := env.GetEnvConfig()

	// Create bot from environment value
	b, err := gotgbot.NewBot(envConfig.BotToken, nil)
	if err != nil {
		panic("Failed to create bot: " + err.Error())
	}

	_, err = b.SetMyCommands(handler.BotCommands, nil)
	if err != nil {
		panic("Failed to set bot commands: " + err.Error())
	}

	// Create dispatcher
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("An error occurred while handling dispatch:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	// Create updater
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		UnhandledErrFunc: func(err error) {
			// If an error is returned by a handler, log it and continue going
			log.Println("An error occurred while handling update: ", err.Error())
		},
	})

	// Add echo handler to reply to all text messages
	handler.ConfigureDispatcher(dispatcher)

	// Start the webhook server. We start the server before we set the webhook itself, so that when telegram starts
	// sending updates, the server is already ready
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  "0.0.0.0:8080",
		SecretToken: envConfig.WebhookSecret,
	}
	err = updater.StartWebhook(b, WEBHOOK_PATH+envConfig.BotToken, webhookOpts)
	if err != nil {
		panic("Failed to start webhook: " + err.Error())
	}

	// Set the webhook on Telegram
	err = updater.SetAllBotWebhooks(envConfig.WebhookDomain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
	})
	if err != nil {
		panic("Failed to set webhook: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
