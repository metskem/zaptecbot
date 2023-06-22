package types

import "fmt"

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
	PhaseRotation        string // 548
	ChargeMode           string // 702
	ChargerOperationMode string // 710
}

func (state ChargerState) String() string {
	return fmt.Sprintf("CommunicationMode:%s, PermanentCableLock:%s, Humidity:%s, TemperatureInternal5:%s, VoltagePhase1:%s, VoltagePhase2:%s, VoltagePhase3:%s, CurrentPhase1:%s, CurrentPhase2:%s, CurrentPhase3:%s, PhaseRotation:%s, ChargeMode:%s, ChargerOperationMode:%s",
		state.CommunicationMode, state.PermanentCableLock, state.Humidity, state.TemperatureInternal5, state.VoltagePhase1, state.VoltagePhase2, state.VoltagePhase3, state.CurrentPhase1, state.CurrentPhase2, state.CurrentPhase3, state.PhaseRotation, state.ChargeMode, state.ChargerOperationMode)
}
