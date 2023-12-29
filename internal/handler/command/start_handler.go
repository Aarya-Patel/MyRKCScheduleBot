package command

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/template"
)

// Start Handler Command
const startCommandName = "start"

// Start Handler Response
func startCommandResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(ctx.Update.Message.Chat.Id, template.WelcomeTemplate, &gotgbot.SendMessageOpts{
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		return err
	}
	return nil
}

var StartHandler = handlers.NewCommand(startCommandName, startCommandResponse)
