package callback

import (
	"context"
	"math"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/google/uuid"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/constant"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/env"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/storage"
	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/template"
)

// Breakout Callback Query Filter
var breakoutCallbackQueryFilter filters.CallbackQuery = func(cq *gotgbot.CallbackQuery) bool {
	hasPrefix := callbackquery.Prefix("BREAKOUT:")(cq)
	if !hasPrefix {
		return false
	}

	// Ensure that the unmarshalled callback query data is properly defined
	return uuid.MustParse(strings.Split(cq.Data, ":")[1]) != uuid.Nil
}

// Fetchs the appropriate entity from Datastore specified by the `BreakoutInstanceId` in the callback data
// Updates that entity and the `/breakout` message with the next step of the breakout workflow
func breakoutCallbackQueryResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	dbCtx := context.Background()
	envConfig := env.GetEnvConfig()
	dbClient, err := storage.GetDataStoreClient(envConfig.ProjectId)
	if err != nil {
		return err
	}

	// Fetch the current breakout progress entity from datastore
	var breakoutProgressData constant.BreakoutProgressData
	key := datastore.IDKey(constant.BreakoutProgressDataKind, ctx.CallbackQuery.Message.MessageId, nil)
	err = dbClient.Get(dbCtx, key, &breakoutProgressData)
	if err != nil {
		return err
	}

	// Extract `BreakoutInstanceId`, current `TransitionState` and the option that was selected
	split := strings.Split(ctx.CallbackQuery.Data, ":")
	breakoutInstanceId := split[1]

	temp, err := strconv.Atoi(split[2])
	if err != nil {
		return err
	}

	currentState := constant.TransitionState(temp)

	temp, err = strconv.Atoi(split[3])
	if err != nil {
		return err
	}

	// If the selected option isn't `Back` then move onto the next state, otherwise, move back a state
	delta := 1
	if temp < 0 {
		delta = -1
	}
	// Ensure that the next state is within the bounds of the TransitionState
	nextStateVal := int(math.Min(math.Max(0, float64(int(currentState)+delta)), float64(constant.ResultStep)))
	nextState := constant.TransitionState(nextStateVal)

	// Update the breakout progess data to the next state and populate the data accordingly
	// breakoutProgressData.State = nextState
	if temp != -1 {
		switch currentState {
		case constant.EISideStep:
			breakoutProgressData.Side = constant.EISide(temp)
		case constant.CentreStep:
			breakoutProgressData.Centre = constant.Centre(temp)
		case constant.WingStep:
			breakoutProgressData.Wing = constant.SevaWing(temp)
		case constant.SevaStep:
			breakoutProgressData.Role = constant.SevaRole(temp)
		}
	} else {
		switch currentState {
		case constant.CentreStep:
			breakoutProgressData.Side = -1
			breakoutProgressData.Centre = -1
		case constant.WingStep:
			breakoutProgressData.Centre = -1
			breakoutProgressData.Wing = -1
		case constant.SevaStep:
			breakoutProgressData.Wing = -1
			breakoutProgressData.Role = -1
		case constant.ResultStep:
			breakoutProgressData.Role = -1
		}
	}

	// Upsert
	key = datastore.IDKey(constant.BreakoutProgressDataKind, ctx.CallbackQuery.Message.MessageId, nil)
	key, err = dbClient.Put(dbCtx, key, &breakoutProgressData)
	if err != nil {
		return err
	}

	// Fetch the message that generated this callback query and edit it's prompt
	// and buttons according to next state
	_, _, err = b.EditMessageText(template.GetPromptForTransitionState(nextState, breakoutProgressData), &gotgbot.EditMessageTextOpts{
		ChatId:      ctx.CallbackQuery.Message.Chat.Id,
		MessageId:   ctx.CallbackQuery.Message.MessageId,
		ParseMode:   "MarkdownV2",
		ReplyMarkup: template.GenerateInlineKeyboardReplyMarkup(nextState, breakoutInstanceId),
	})

	if err != nil {
		return err
	}

	return nil
}

var BreakoutCallbackQueryHandler = handlers.NewCallback(breakoutCallbackQueryFilter, breakoutCallbackQueryResponse)
