package model

import (
	"fmt"
	"time"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ChargerStatesRaw []ChargerStateRaw

type ChargerStateRaw struct {
	ChargerID     string `json:"ChargerId"`
	StateID       int    `json:"StateId"`
	Timestamp     string `json:"Timestamp"`
	ValueAsString string `json:"ValueAsString,omitempty"`
}

type ChargerState struct {
	CommunicationMode    string // 150
	PermanentCableLock   string // 151
	TemperatureInternal5 string // 201
	Humidity             string // 270
	VoltagePhase1        string // 501
	VoltagePhase2        string // 502
	VoltagePhase3        string // 503
	CurrentPhase1        string // 507
	CurrentPhase2        string // 508
	CurrentPhase3        string // 509
	ChargerMaxCurrent    string // 510
	TotalChargePower     string // 513
	PhaseRotation        string // 548
	ChargeMode           string // 702
	ChargerOperationMode string // 710
	MainboardVersion     string // 708
	ComputerVersion      string // 911
	SourceVersion        string // 916
}

var (
	ChargerOperationMode0 = "Unknown"
	ChargerOperationMode1 = "Disconnected"
	ChargerOperationMode2 = "Connected_Requesting"
	ChargerOperationMode3 = "Connected_Charging"
	ChargerOperationMode5 = "Connected_Finished"
)

func (state ChargerState) String() string {
	return fmt.Sprintf("CommunicationMode: %s\nPermanentCableLock: %s\nHumidity: %s\nTemperatureInternal5: %s\nPhase1: %sV (%sA)\nPhase2: %sV (%sA)\nPhase3: %sV (%sA)\nPhaseRotation: %s\nChargeMode: %s\nChargerOperationMode: %s\nChargerMaxCurrent: %s\nMainboardVersion: %s\nComputerVersion: %s\nSourceVersion: %s",
		state.CommunicationMode, state.PermanentCableLock, state.Humidity, state.TemperatureInternal5, state.VoltagePhase1, state.CurrentPhase1, state.VoltagePhase2, state.CurrentPhase2, state.VoltagePhase3, state.CurrentPhase3, state.PhaseRotation, state.ChargeMode, state.ChargerOperationMode, state.ChargerMaxCurrent, state.MainboardVersion, state.ComputerVersion, state.SourceVersion)
}

type Schedule struct {
	StartTime      time.Time
	ChargeDuration time.Duration
	InProgress     bool
}

func (schedule Schedule) String() string {
	return fmt.Sprintf("startTime:%s, duration:%d, inProgress:%t", schedule.StartTime.Format("15:04"), int(schedule.ChargeDuration.Hours()), schedule.InProgress)
}
func (schedule Schedule) Key() string {
	return fmt.Sprintf("%s %d", schedule.StartTime.Format("15:04"), int(schedule.ChargeDuration.Hours()))
}
