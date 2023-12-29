package template

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"regexp"
	"strconv"

	"github.com/Aarya-Patel/MyRKCScheduleBot/internal/constant"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

const EISidePrompt = `
To find your breakout room, please help in answering the following sequence of questions by selecting one of the options below\.

*\[1/4\]* Which side does your seva reside within?
`

const CentrePrompt = `
To find your breakout room, please help in answering the following sequence of questions by selecting one of the options below\.

*\[2/4\]* What is your centre?
`

const SevaWingPrompt = `
To find your breakout room, please help in answering the following sequence of questions by selecting one of the options below\.

*\[3/4\]* What is your seva wing?
`

const SevaPrompt = `
To find your breakout room, please help in answering the following sequence of questions by selecting one of the options below\.

*\[4/4\]* What is your primary seva?
`

const BreakoutProgressDataTemplate = `
__This is your primary seva breakdown:__
*Side*: {{.Side.String}}
*Centre*: {{.Centre.String}}
*Wing*: {{.Wing.String}}
*Seva*: {{.Role.String}}
`

const BreakoutRoomTemplate = `
🗓 *Date*: {{.Time.DisplayDate}}
⏰ *Time*: {{.Time.DisplayTime}}
{{if .Room}}📍 *Room*: {{.Room}}{{else}}💻 *Link*: [{{.Link}}]({{.Link}})
{{if .Password}}🔐 *Password*: {{.Password}}{{end}}{{end}}
`

var TransitionStateToPrompt map[constant.TransitionState]string = map[constant.TransitionState]string{
	constant.EISideStep: EISidePrompt,
	constant.CentreStep: CentrePrompt,
	constant.WingStep:   SevaWingPrompt,
	constant.SevaStep:   SevaPrompt,
}

func escapeStringFromMarkdown(src []byte) []byte {
	// Wrap the special character inside a group
	specialChars := []string{"(-)", "(\\.)", "(\\!)"}
	for _, ch := range specialChars {
		re, err := regexp.Compile(ch)
		if err != nil {
			log.Println(err.Error())
			return src // return unedited version in case of error
		}
		src = re.ReplaceAll(src, []byte("\\$1"))
	}
	return src
}

func formatBreakoutRooms(breakoutProgressData constant.BreakoutProgressData, breakoutRooms []constant.BreakoutRoomData) string {
	var fmtBreakoutRooms bytes.Buffer
	progressTmpl, err := template.New("BreakoutProgressTemplate").Parse(BreakoutProgressDataTemplate)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	roomTmpl, err := template.New("BreakoutRoomTemplate").Parse(BreakoutRoomTemplate)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	fmtBreakoutRooms.Write([]byte("_Below is a catered list of your breakout sessions._\n"))
	err = progressTmpl.Execute(&fmtBreakoutRooms, breakoutProgressData)
	fmtBreakoutRooms.Write([]byte("\n"))
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	for idx, breakoutRoom := range breakoutRooms {
		fmtBreakoutRooms.Write([]byte("*__Breakout " + strconv.Itoa(idx+1) + "__*"))
		err = roomTmpl.Execute(&fmtBreakoutRooms, breakoutRoom)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		if idx != len(breakoutRooms)-1 {
			fmtBreakoutRooms.Write([]byte("\n"))
		}
	}

	if len(breakoutRooms) == 0 {
		fmtBreakoutRooms.Write([]byte("🏖 *No breakout sessions scheduled!*"))
	}

	// Finally escape all special characters deemed by MarkdownV2 parse mode
	escFmtBreakoutRooms := escapeStringFromMarkdown(fmtBreakoutRooms.Bytes())
	return string(escFmtBreakoutRooms)
}

func GetPromptForTransitionState(state constant.TransitionState, breakoutProgressData constant.BreakoutProgressData) string {
	// If we are at the last step of the breakout workflow, fetch breakout data
	if state == constant.ResultStep {
		data, err := constant.GetBreakoutData()
		if err != nil {
			log.Println(err.Error())
			return "Oops something went wrong!"
		}
		// log.Println(data[""])
		breakoutRooms := data[breakoutProgressData.Side.String()][breakoutProgressData.Centre.String()][breakoutProgressData.Wing.String()][breakoutProgressData.Role.String()]
		if err != nil {
			log.Println(err.Error())
			return "Oops something went wrong!"
		}

		return formatBreakoutRooms(breakoutProgressData, breakoutRooms)
	}
	return TransitionStateToPrompt[state]
}

func GenerateInlineKeyboardReplyMarkup(state constant.TransitionState, breakoutInstanceId string) gotgbot.InlineKeyboardMarkup {
	markup := gotgbot.InlineKeyboardMarkup{}

	switch state {
	case constant.EISideStep:
		markup.InlineKeyboard = generateEISideStepButtons(state, breakoutInstanceId)
	case constant.CentreStep:
		markup.InlineKeyboard = generateCentreStepButtons(state, breakoutInstanceId)
	case constant.WingStep:
		markup.InlineKeyboard = generateWingStepButtons(state, breakoutInstanceId)
	case constant.SevaStep:
		markup.InlineKeyboard = generateSevaStepButtons(state, breakoutInstanceId)
	case constant.ResultStep:
		markup.InlineKeyboard = generateResultStepButtons(state, breakoutInstanceId)
	default:
		markup.InlineKeyboard = [][]gotgbot.InlineKeyboardButton{}
	}

	return markup
}

func generateEISideStepButtons(state constant.TransitionState, breakoutInstanceId string) [][]gotgbot.InlineKeyboardButton {
	var buttons [][]gotgbot.InlineKeyboardButton = make([][]gotgbot.InlineKeyboardButton, 0)
	options := []constant.EISide{
		constant.ESide,
		constant.ISide,
	}

	const numOptionsPerRow int = 2
	numGroups := int(math.Ceil(float64(len(options)) / float64(numOptionsPerRow)))

	// Group the options in pairs. So each row has 2
	for i := 0; i < numGroups; i++ {
		curOptions := options[i*numOptionsPerRow : int(math.Min(float64(len(options)), float64(i*numOptionsPerRow+numOptionsPerRow)))]
		curRow := make([]gotgbot.InlineKeyboardButton, 0)

		for _, option := range curOptions {
			curRow = append(curRow, gotgbot.InlineKeyboardButton{
				Text:         option.String(),
				CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":" + strconv.Itoa(int(option)),
			})
		}
		buttons = append(buttons, curRow)
	}

	return buttons
}

func generateCentreStepButtons(state constant.TransitionState, breakoutInstanceId string) [][]gotgbot.InlineKeyboardButton {
	var buttons [][]gotgbot.InlineKeyboardButton = make([][]gotgbot.InlineKeyboardButton, 0)
	options := []constant.Centre{
		constant.Brandon,
		constant.Calgary,
		constant.Cambridge,
		constant.Edmonton,
		constant.FortMcMurray,
		constant.Montreal,
		constant.Ottawa,
		constant.Regina,
		constant.Saskatoon,
		constant.Scarborough,
		constant.Toronto,
		constant.Vancouver,
		constant.Windsor,
		constant.Winnipeg,
	}

	const numOptionsPerRow int = 2
	numGroups := int(math.Ceil(float64(len(options)) / float64(numOptionsPerRow)))

	// Group the options in pairs. So each row has 2
	for i := 0; i < numGroups; i++ {
		curOptions := options[i*numOptionsPerRow : int(math.Min(float64(len(options)), float64(i*numOptionsPerRow+numOptionsPerRow)))]
		curRow := make([]gotgbot.InlineKeyboardButton, 0)

		for _, option := range curOptions {
			curRow = append(curRow, gotgbot.InlineKeyboardButton{
				Text:         option.String(),
				CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":" + strconv.Itoa(int(option)),
			})
		}
		buttons = append(buttons, curRow)
	}

	// Add back button
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{
			Text:         "⏪ Back",
			CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":-1",
		},
	})
	return buttons
}

func generateWingStepButtons(state constant.TransitionState, breakoutInstanceId string) [][]gotgbot.InlineKeyboardButton {
	var buttons [][]gotgbot.InlineKeyboardButton = make([][]gotgbot.InlineKeyboardButton, 0)
	options := []constant.SevaWing{
		constant.BalBalikaMandal,
		constant.KishoreKishoriMandal,
		constant.YuvakYuvatiMandal,
	}

	const numOptionsPerRow int = 1
	numGroups := int(math.Ceil(float64(len(options)) / float64(numOptionsPerRow)))

	// Group the options in pairs. So each row has 2
	for i := 0; i < numGroups; i++ {
		curOptions := options[i*numOptionsPerRow : int(math.Min(float64(len(options)), float64(i*numOptionsPerRow+numOptionsPerRow)))]
		curRow := make([]gotgbot.InlineKeyboardButton, 0)

		for _, option := range curOptions {
			curRow = append(curRow, gotgbot.InlineKeyboardButton{
				Text:         option.String(),
				CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":" + strconv.Itoa(int(option)),
			})
		}
		buttons = append(buttons, curRow)
	}

	// Add back button
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{
			Text:         "⏪ Back",
			CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":-1",
		},
	})
	return buttons
}

func generateSevaStepButtons(state constant.TransitionState, breakoutInstanceId string) [][]gotgbot.InlineKeyboardButton {
	var buttons [][]gotgbot.InlineKeyboardButton = make([][]gotgbot.InlineKeyboardButton, 0)
	options := []constant.SevaRole{
		constant.RCT,
		constant.Coordinator,
		constant.RC,
		constant.SabhaKaryakar,
		constant.NDC,
		constant.PC,
		constant.GC,
		constant.LocalAnalyst,
		constant.SevaTraining,
		constant.NetworkingKaryakar,
		constant.CampusSabhaKaryakar,
	}

	const numOptionsPerRow int = 2
	numGroups := int(math.Ceil(float64(len(options)) / float64(numOptionsPerRow)))

	// Group the options in pairs. So each row has 2
	for i := 0; i < numGroups; i++ {
		curOptions := options[i*numOptionsPerRow : int(math.Min(float64(len(options)), float64(i*numOptionsPerRow+numOptionsPerRow)))]
		curRow := make([]gotgbot.InlineKeyboardButton, 0)

		for _, option := range curOptions {
			curRow = append(curRow, gotgbot.InlineKeyboardButton{
				Text:         option.String(),
				CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":" + strconv.Itoa(int(option)),
			})
		}
		buttons = append(buttons, curRow)
	}

	// Add back button
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{
			Text:         "⏪ Back",
			CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":-1",
		},
	})
	return buttons
}

func generateResultStepButtons(state constant.TransitionState, breakoutInstanceId string) [][]gotgbot.InlineKeyboardButton {
	var buttons [][]gotgbot.InlineKeyboardButton = make([][]gotgbot.InlineKeyboardButton, 0)

	// Add back button
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{
			Text:         "⏪ Back",
			CallbackData: "BREAKOUT:" + breakoutInstanceId + ":" + strconv.Itoa(int(state)) + ":-1",
		},
	})
	return buttons
}
