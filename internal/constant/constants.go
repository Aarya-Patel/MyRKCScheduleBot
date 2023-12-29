package constant

import (
	"encoding/json"
	"fmt"
)

type EISide int8

const (
	ESide EISide = iota
	ISide
)

var EISideToDisplayString = map[EISide]string{
	ESide: `eSide`,
	ISide: `iSide`,
}

func (e EISide) String() string {
	return EISideToDisplayString[e]
}

func (e *EISide) UnmarshalJSON(data []byte) error {
	var str string = string(data)
	for k, v := range EISideToDisplayString {
		if v == str {
			*e = k
			return nil
		}
	}
	return fmt.Errorf("Cannot unmarshal into EISide")
}

type Centre int16

const (
	Brandon Centre = iota
	Calgary
	Cambridge
	Edmonton
	FortMcMurray
	Montreal
	Ottawa
	Regina
	Saskatoon
	Scarborough
	Toronto
	Vancouver
	Windsor
	Winnipeg
)

var CentreToDisplayString = map[Centre]string{
	Brandon:      "Brandon",
	Calgary:      "Calgary",
	Cambridge:    "Cambridge",
	Edmonton:     "Edmonton",
	FortMcMurray: "Fort McMurray",
	Montreal:     "Montreal",
	Ottawa:       "Ottawa",
	Regina:       "Regina",
	Saskatoon:    "Saskatoon",
	Scarborough:  "Scarborough",
	Toronto:      "Toronto",
	Vancouver:    "Vancouver",
	Windsor:      "Windsor",
	Winnipeg:     "Winnipeg",
}

func (e Centre) String() string {
	return CentreToDisplayString[e]
}

func (e *Centre) UnmarshalJSON(data []byte) error {
	var str string = string(data)
	for k, v := range CentreToDisplayString {
		if v == str {
			*e = k
			return nil
		}
	}
	return fmt.Errorf("Cannot unmarshal into Centre")
}

type SevaWing int8

const (
	BalBalikaMandal SevaWing = iota
	KishoreKishoriMandal
	YuvakYuvatiMandal
)

var SevaWingToDisplayString = map[SevaWing]string{
	BalBalikaMandal:      "Bal/Balika Mandal",
	KishoreKishoriMandal: "Kishore/Kishori Mandal",
	YuvakYuvatiMandal:    "Yuvak/Yuvati Mandal",
}

func (e SevaWing) String() string {
	return SevaWingToDisplayString[e]
}

func (e *SevaWing) UnmarshalJSON(data []byte) error {
	var str string = string(data)
	for k, v := range SevaWingToDisplayString {
		if v == str {
			*e = k
			return nil
		}
	}
	return fmt.Errorf("Cannot unmarshal into SevaWing")
}

type SevaRole int16

const (
	RCT SevaRole = iota
	Coordinator
	RC
	SabhaKaryakar
	NDC
	PC
	GC
	LocalAnalyst
	SevaTraining
	CampusSabhaKaryakar
	NetworkingKaryakar
)

var SevaRoleToDisplayString = map[SevaRole]string{
	RCT:                 "RCT",
	Coordinator:         "Coordinator",
	RC:                  "RC",
	SabhaKaryakar:       "Sabha Karyakar",
	NDC:                 "NDC",
	PC:                  "PC",
	GC:                  "GC",
	LocalAnalyst:        "Local Analyst",
	SevaTraining:        "Seva Training",
	CampusSabhaKaryakar: "Campus Sabha Karyakar",
	NetworkingKaryakar:  "Networking Karyakar",
}

func (e SevaRole) String() string {
	return SevaRoleToDisplayString[e]
}

func (e *SevaRole) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	for k, v := range SevaRoleToDisplayString {
		if v == str {
			*e = k
			return nil
		}
	}
	return fmt.Errorf("Cannot unmarshal into SevaRole")
}

type TransitionState int8

//go:generate stringer -type=TransitionState
const (
	EISideStep TransitionState = iota
	CentreStep
	WingStep
	SevaStep
	ResultStep
)
