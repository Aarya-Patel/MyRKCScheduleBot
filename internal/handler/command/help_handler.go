package command

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/template"
)

// Help Handler Command
const helpCommandName = "help"

// Start Handler Response
func helpCommandResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(ctx.Update.Message.Chat.Id, template.HelpTemplate, &gotgbot.SendMessageOpts{
		ParseMode: "MarkdownV2",
	})
	if err != nil {
		return err
	}

	return nil
}

var HelpHandler = handlers.NewCommand(helpCommandName, helpCommandResponse)
