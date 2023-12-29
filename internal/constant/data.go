package constant

import (
	"encoding/json"
	"os"
	"time"
)

const BreakoutProgressDataKind = "BreakoutProgressData"

type CustomBreakoutTime struct {
	time.Time
}

func (t CustomBreakoutTime) DisplayDate() string {
	return t.Time.Format("2006-01-02")
}

func (t CustomBreakoutTime) DisplayTime() string {
	return t.Time.Format("3:04PM")
}

func (t *CustomBreakoutTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	time, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}

	*t = CustomBreakoutTime{time}
	return nil
}

// Data regarding the progess made during the breakout workflow.
// This will be stored in the database.
type BreakoutProgressData struct {
	State              TransitionState
	BreakoutInstanceId string
	UserId             int64
	Side               EISide
	Centre             Centre
	Wing               SevaWing
	Role               SevaRole
}

// End data of each breakout room
type BreakoutRoomData struct {
	Time     CustomBreakoutTime `json:",omitempty"`
	Room     string             `json:",omitempty"`
	Link     string             `json:",omitempty"`
	Password string             `json:",omitempty"`
}

type BreakoutSevaRoleData map[string][]BreakoutRoomData

type BreakoutSevaWingData map[string]BreakoutSevaRoleData

type BreakoutCentreData map[string]BreakoutSevaWingData

type BreakoutEISideData map[string]BreakoutCentreData

var breakoutData BreakoutEISideData

func GetBreakoutData() (BreakoutEISideData, error) {
	if breakoutData == nil {
		data, err := os.ReadFile("./internal/constant/data.json")
		if err != nil {
			return nil, err
		}

		breakoutData = make(BreakoutEISideData)
		err = json.Unmarshal(data, &breakoutData)
		if err != nil {
			return nil, err
		}
	}

	return breakoutData, nil
}
