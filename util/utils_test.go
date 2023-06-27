package util

import (
	"fmt"
	"github.com/metskem/zaptecbot/model"
	"testing"
)

func TestParseChargerState(t *testing.T) {
	rawStates := []model.ChargerStateRaw{
		{ChargerID: "fake charger id", StateID: 150, ValueAsString: "s150"},
		{ChargerID: "fake charger id", StateID: 151, ValueAsString: "s151"},
		{ChargerID: "fake charger id", StateID: 201, ValueAsString: "s201"},
		{ChargerID: "fake charger id", StateID: 270, ValueAsString: "s270"},
		{ChargerID: "fake charger id", StateID: 501, ValueAsString: "s501"},
		{ChargerID: "fake charger id", StateID: 502, ValueAsString: "s502"},
		{ChargerID: "fake charger id", StateID: 503, ValueAsString: "s503"},
		{ChargerID: "fake charger id", StateID: 507, ValueAsString: "s507"},
		{ChargerID: "fake charger id", StateID: 508, ValueAsString: "s508"},
		{ChargerID: "fake charger id", StateID: 509, ValueAsString: "s509"},
		{ChargerID: "fake charger id", StateID: 548, ValueAsString: "s548"},
		{ChargerID: "fake charger id", StateID: 702, ValueAsString: "s702"},
		{ChargerID: "fake charger id", StateID: 710, ValueAsString: "3"},
	}
	state := ParseChargerState(rawStates)
	//t.Log(state)
	expectedState := "CommunicationMode:s150, PermanentCableLock:s151, Humidity:s270, TemperatureInternal5:s201, VoltagePhase1:s501, VoltagePhase2:s502, VoltagePhase3:s503, CurrentPhase1:s507, CurrentPhase2:s508, CurrentPhase3:s509, PhaseRotation:s548, ChargeMode:s702, ChargerOperationMode:Connected_Charging"
	{
		if fmt.Sprintf("%s", state) != expectedState {
			t.Errorf("unexpected state (expected / actual) : \n%s\n%s\n", expectedState, state)
		}
	}
}

func TestValidateScheduleAdd(t *testing.T) {
	validStrings := []string{"/sa 12:34 5", "/sa 12:34 15"}
	for _, testString := range validStrings {
		if _, err := ParseSchedule(testString); err != nil {
			t.Errorf(err.Error())
		}
	}
	invalidStrings := []string{"/sa 12:34 xx", "/sa 12:34 ", "/sa 12:3a 1", "/sa  12:34 9", "/sa 25:12 1", "/sa 12:34 0", "/sa 12:34 99"}
	for _, testString := range invalidStrings {
		if _, err := ParseSchedule(testString); err == nil {
			t.Errorf("ParseSchedule should have failed for \"%s\"", testString)
		}
	}
}
