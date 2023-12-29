package command

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/google/uuid"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/constant"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/env"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/storage"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/template"
)

// Breakout Handler Command
const breakoutCommandName = "breakout"

// Schedule Handler Response
func breakoutCommandResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	dbCtx := context.Background()
	envConfig := env.GetEnvConfig()
	dbClient, dbErr := storage.GetDataStoreClient(envConfig.ProjectId)
	if dbErr != nil {
		return dbErr
	}

	breakoutProgressData := constant.BreakoutProgressData{
		State:              constant.EISideStep,
		BreakoutInstanceId: uuid.New().String(),
		UserId:             ctx.Update.Message.From.Id,
		Side:               -1,
		Centre:             -1,
		Wing:               -1,
	}

	msg, err := b.SendMessage(ctx.Update.Message.Chat.Id, template.GetPromptForTransitionState(breakoutProgressData.State, breakoutProgressData), &gotgbot.SendMessageOpts{
		ParseMode:   "MarkdownV2",
		ReplyMarkup: template.GenerateInlineKeyboardReplyMarkup(breakoutProgressData.State, breakoutProgressData.BreakoutInstanceId),
	})
	if err != nil {
		return err
	}

	// Add an BreakoutProgressData entity into the database regarding this message
	key := datastore.IDKey(constant.BreakoutProgressDataKind, msg.MessageId, nil)
	key, err = dbClient.Put(dbCtx, key, &breakoutProgressData)
	if err != nil {
		return err
	}

	return nil
}

var BreakoutHandler = handlers.NewCommand(breakoutCommandName, breakoutCommandResponse)
