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
	ID                                         string  `json:"Id"`
	Name                                       string  `json:"Name"`
	Address                                    string  `json:"Address"`
	ZipCode                                    string  `json:"ZipCode"`
	City                                       string  `json:"City"`
	CountryID                                  string  `json:"CountryId"`
	InstallationType                           int     `json:"InstallationType"`
	MaxCurrent                                 int     `json:"MaxCurrent"`
	AvailableCurrent                           int     `json:"AvailableCurrent"`
	AvailableCurrentPhase1                     int     `json:"AvailableCurrentPhase1"`
	AvailableCurrentPhase2                     int     `json:"AvailableCurrentPhase2"`
	AvailableCurrentPhase3                     int     `json:"AvailableCurrentPhase3"`
	AvailableCurrentMode                       int     `json:"AvailableCurrentMode"`
	AvailableCurrentScheduleWeekendActive      bool    `json:"AvailableCurrentScheduleWeekendActive"`
	DefaultThreeToOneSwitchCurrent             int     `json:"DefaultThreeToOneSwitchCurrent"`
	InstallationCategoryID                     string  `json:"InstallationCategoryId"`
	InstallationCategory                       string  `json:"InstallationCategory"`
	UseLoadBalancing                           bool    `json:"UseLoadBalancing"`
	IsRequiredAuthentication                   bool    `json:"IsRequiredAuthentication"`
	Latitude                                   float64 `json:"Latitude"`
	Longitude                                  float64 `json:"Longitude"`
	Active                                     bool    `json:"Active"`
	NetworkType                                int     `json:"NetworkType"`
	AvailableInternetAccessPLC                 bool    `json:"AvailableInternetAccessPLC"`
	AvailableInternetAccessWiFi                bool    `json:"AvailableInternetAccessWiFi"`
	CreatedOnDate                              string  `json:"CreatedOnDate"`
	UpdatedOn                                  string  `json:"UpdatedOn"`
	CurrentUserRoles                           int     `json:"CurrentUserRoles"`
	AuthenticationType                         int     `json:"AuthenticationType"`
	MessagingEnabled                           bool    `json:"MessagingEnabled"`
	RoutingID                                  string  `json:"RoutingId"`
	OcppCloudURLVersion                        int     `json:"OcppCloudUrlVersion"`
	TimeZoneName                               string  `json:"TimeZoneName"`
	TimeZoneIanaName                           string  `json:"TimeZoneIanaName"`
	IsSubscriptionsAvailableForCurrentUser     bool    `json:"IsSubscriptionsAvailableForCurrentUser"`
	AvailableFeatures                          int     `json:"AvailableFeatures"`
	EnabledFeatures                            int     `json:"EnabledFeatures"`
	ActiveChargerCount                         int     `json:"ActiveChargerCount"`
	FeaturePowerManagementEcoModeDepartureTime int     `json:"Feature_PowerManagement_EcoMode_DepartureTime"`
	FeaturePowerManagementEcoModeMinEnergy     int     `json:"Feature_PowerManagement_EcoMode_MinEnergy"`
	FeaturePowerManagementEcoModeDeliveryArea  int     `json:"Feature_PowerManagement_EcoMode_DeliveryArea"`
	PropertyIsMinimumPowerOfflineMode          bool    `json:"PropertyIsMinimumPowerOfflineMode"`
	PropertyOfflineModeAllowAnonymous          bool    `json:"PropertyOfflineModeAllowAnonymous"`
	PropertyEnergySensorRippleEnabled          bool    `json:"PropertyEnergySensorRippleEnabled"`
	PropertyEnergySensorRippleNumBits          int     `json:"PropertyEnergySensorRippleNumBits"`
	Tic                                        struct {
		Enabled bool `json:"Enabled"`
	} `json:"Tic"`
	SurplusMode struct {
		Active   bool `json:"Active"`
		Strategy int  `json:"Strategy"`
	} `json:"SurplusMode"`
	PropertyFirmwareAutomaticUpdates    bool `json:"PropertyFirmwareAutomaticUpdates"`
	PropertyMaxSinglePhaseChargeCurrent int  `json:"PropertyMaxSinglePhaseChargeCurrent"`
	PropertySessionMaxStopCount         int  `json:"PropertySessionMaxStopCount"`
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
	return fmt.Sprintf("Active: %t\nMaxCurrent:%d\nAvailableCurrent:%d", detail.Active, detail.MaxCurrent, detail.AvailableCurrent)
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
