package handler

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/handler/callback"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/handler/command"
)

var str string = "asdf"
var BotCommands []gotgbot.BotCommand = []gotgbot.BotCommand{
	{
		Command:     "start",
		Description: "Displays the welcome message",
	},
	{
		Command:     "help",
		Description: "Displays helpful information regarding the bot",
	},
	{
		Command:     "breakout",
		Description: "Initiates breakout schedule workflow",
	},
}

func ConfigureDispatcher(dispatcher *ext.Dispatcher) {
	dispatcher.AddHandler(command.StartHandler)
	dispatcher.AddHandler(command.HelpHandler)
	dispatcher.AddHandler(command.BreakoutHandler)
	dispatcher.AddHandler(callback.BreakoutCallbackQueryHandler)
	dispatcher.AddHandler(handlers.NewMessage(message.Text, incrompehensibleMessageResponse)) // Catch All Handler
}

// echo replies to a messages with its own contents.
func incrompehensibleMessageResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Sorry, I didn't get that\\. Use `/help` to learn more about the tool\\!", &gotgbot.SendMessageOpts{
		ParseMode: "MarkdownV2",
	})
	if err != nil {
		return fmt.Errorf("Failed to respond to an incrompehensible message: %w", err)
	}
	return nil
}
