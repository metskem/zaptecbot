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

type InstallationDetails struct {
	Pages      int `json:"Pages"`
	TotalCount int `json:"TotalCount"`
	Data       []struct {
		ID                                     string  `json:"Id"`
		Name                                   string  `json:"Name"`
		Address                                string  `json:"Address"`
		ZipCode                                string  `json:"ZipCode"`
		City                                   string  `json:"City"`
		CountryID                              string  `json:"CountryId"`
		InstallationType                       int     `json:"InstallationType"`
		MaxCurrent                             float32 `json:"MaxCurrent"`
		AvailableCurrentMode                   int     `json:"AvailableCurrentMode"`
		AvailableCurrentScheduleWeekendActive  bool    `json:"AvailableCurrentScheduleWeekendActive"`
		DefaultThreeToOneSwitchCurrent         float32 `json:"DefaultThreeToOneSwitchCurrent"`
		InstallationCategoryID                 string  `json:"InstallationCategoryId"`
		InstallationCategory                   string  `json:"InstallationCategory"`
		UseLoadBalancing                       bool    `json:"UseLoadBalancing"`
		IsRequiredAuthentication               bool    `json:"IsRequiredAuthentication"`
		Latitude                               float64 `json:"Latitude"`
		Longitude                              float64 `json:"Longitude"`
		Active                                 bool    `json:"Active"`
		NetworkType                            int     `json:"NetworkType"`
		AvailableInternetAccessPLC             bool    `json:"AvailableInternetAccessPLC"`
		AvailableInternetAccessWiFi            bool    `json:"AvailableInternetAccessWiFi"`
		CreatedOnDate                          string  `json:"CreatedOnDate"`
		UpdatedOn                              string  `json:"UpdatedOn"`
		CurrentUserRoles                       int     `json:"CurrentUserRoles"`
		AuthenticationType                     int     `json:"AuthenticationType"`
		MessagingEnabled                       bool    `json:"MessagingEnabled"`
		RoutingID                              string  `json:"RoutingId"`
		OcppCloudURLVersion                    int     `json:"OcppCloudUrlVersion"`
		IsSubscriptionsAvailableForCurrentUser bool    `json:"IsSubscriptionsAvailableForCurrentUser"`
		AvailableFeatures                      int     `json:"AvailableFeatures"`
		EnabledFeatures                        int     `json:"EnabledFeatures"`
		PropertyFirmwareAutomaticUpdates       bool    `json:"PropertyFirmwareAutomaticUpdates"`
	} `json:"Data"`
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
	StandAlone           string // 712
	MainboardVersion     string // 708
	ComputerVersion      string // 911
	SourceVersion        string // 916
}

var (
	ChargerOperationModeUnknown             = "Unknown"
	ChargerOperationModeDisconnected        = "Disconnected"
	ChargeroperationmodeconnectedRequesting = "Connected_Requesting"
	ChargeroperationmodeconnectedCharging   = "Connected_Charging"
	ChargeroperationmodeconnectedFinished   = "Connected_Finished"
)

func (state ChargerState) String() string {
	return fmt.Sprintf("CommunicationMode: %s\nPermanentCableLock: %s\nHumidity: %s\nTemperatureInternal5: %s\nPhase1: %sV (%sA)\nPhase2: %sV (%sA)\nPhase3: %sV (%sA)\nPhaseRotation: %s\nChargeMode: %s\nChargerOperationMode: %s\nIsStandAlone: %s\nChargerMaxCurrent: %sA\nMainboardVersion: %s\nComputerVersion: %s\nSourceVersion: %s",
		state.CommunicationMode, state.PermanentCableLock, state.Humidity, state.TemperatureInternal5, state.VoltagePhase1, state.CurrentPhase1, state.VoltagePhase2, state.CurrentPhase2, state.VoltagePhase3, state.CurrentPhase3, state.PhaseRotation, state.ChargeMode, state.ChargerOperationMode, state.StandAlone, state.ChargerMaxCurrent, state.MainboardVersion, state.ComputerVersion, state.SourceVersion)
}

func (detail InstallationDetails) String() string {
	if len(detail.Data) == 0 {
		return "no installation data available"
	}
	return fmt.Sprintf("Active: %t\nMaxCurrent:%4.1f", detail.Data[0].Active, detail.Data[0].MaxCurrent)
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
