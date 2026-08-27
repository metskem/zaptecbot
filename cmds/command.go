package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/model"
	"github.com/metskem/zaptecbot/util"
)

func ChargerState() (model.ChargerState, error) {
	jwToken := util.GetToken()
	if jwToken == "" {
		return model.ChargerState{}, fmt.Errorf("failed to get token")
	}
	transport := http.Transport{IdleConnTimeout: time.Second}
	client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(conf.ChargerStateUrl, conf.ChargerId), nil)
	if err != nil {
		return model.ChargerState{}, fmt.Errorf("failed to create http request: %s", err)
	}
	req.Header = map[string][]string{"Accept": {"*/*"}, "Authorization": {fmt.Sprintf("bearer %s", jwToken)}}
	resp, err := client.Do(req)
	if err != nil {
		return model.ChargerState{}, fmt.Errorf("response from charger state failed: %s", err)
	}
	if resp != nil {
		respBody, _ := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusOK {
			stateResponse := model.ChargerStatesRaw{}
			if err := json.Unmarshal(respBody, &stateResponse); err != nil {
				return model.ChargerState{}, fmt.Errorf("failed to decode the charger state response: %s", err)
			}
			return util.ParseChargerState(stateResponse), nil
		}
		return model.ChargerState{}, fmt.Errorf("response (%d) from charge state failed: %s", resp.StatusCode, respBody)
	}
	return model.ChargerState{}, fmt.Errorf("response from charger state was nil")
}

func InstallationDetails() (model.InstallationDetails, error) {
	jwToken := util.GetToken()
	if jwToken == "" {
		return model.InstallationDetails{}, fmt.Errorf("failed to get token")
	}
	transport := http.Transport{IdleConnTimeout: time.Second}
	client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
	req, err := http.NewRequest(http.MethodGet, conf.InstallationDetailsUrl, nil)
	if err != nil {
		return model.InstallationDetails{}, fmt.Errorf("failed to create http request: %s", err)
	}
	req.Header = map[string][]string{"Accept": {"*/*"}, "Authorization": {fmt.Sprintf("bearer %s", jwToken)}}
	resp, err := client.Do(req)
	if err != nil {
		return model.InstallationDetails{}, fmt.Errorf("response from installation details failed: %s", err)
	}
	if resp != nil {
		respBody, _ := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusOK {
			detailsResponse := model.InstallationDetails{}
			if err := json.Unmarshal(respBody, &detailsResponse); err != nil {
				return model.InstallationDetails{}, fmt.Errorf("failed to decode the installation details response: %s", err)
			}
			return detailsResponse, nil
		}
		return model.InstallationDetails{}, fmt.Errorf("response (%d) from installation details failed: %s", resp.StatusCode, respBody)
	}
	return model.InstallationDetails{}, fmt.Errorf("response from installation details was nil")
}

func InstallationSetMaxCurrent(maxCurrent int) error {
	jwToken := util.GetToken()
	if jwToken == "" {
		return fmt.Errorf("failed to get token")
	}
	installationDetails, err := InstallationDetails()
	if err != nil {
		return fmt.Errorf("failed to get charger details: %s", err)
	}
	if len(installationDetails.Data) == 0 {
		return fmt.Errorf("no installation data available")
	}

	transport := http.Transport{IdleConnTimeout: time.Second}
	client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
	body := strings.NewReader(fmt.Sprintf(`{"maxCurrent":%d,"availableCurrent":%d}`, maxCurrent, maxCurrent))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(conf.InstallationUpdateUrl, installationDetails.Data[0].ID), body)
	if err != nil {
		return fmt.Errorf("failed to create http request: %s", err)
	}
	req.Header = map[string][]string{"Accept": {"*/*"}, "Content-Type": {"application/json"}, "Authorization": {fmt.Sprintf("bearer %s", jwToken)}}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("response from installation SetMaxCurrent failed: %s", err)
	}
	if resp != nil {
		respBody, _ := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		if resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("response (%d) from installation SetMaxCurrent failed: %s", resp.StatusCode, respBody)
	}
	return errors.New("response from installation SetMaxCurrent failed, nil response")
}

func ShowState(update tgbotapi.Update) {
	chargerState, err := ChargerState()
	if err != nil {
		log.Println(err)
		util.Broadcast(err.Error())
		return
	}
	chargerDetails, err := InstallationDetails()
	if err != nil {
		log.Println(err)
		util.Broadcast(err.Error())
		return
	}
	util.SendMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%s", chargerState, chargerDetails), false)
}

func StartStopCharger(cmd string) {
	var cmdCode int
	switch cmd {
	case "start":
		cmdCode = 507
	case "stop":
		cmdCode = 506
	}

	// Check current charger state before proceeding
	chargerState, err := ChargerState()
	if err != nil {
		msg := fmt.Sprintf("failed to get charger state before %s command: %s", cmd, err)
		log.Println(msg)
		util.Broadcast(msg)
		return
	}
	currentMode := chargerState.ChargerOperationMode
	log.Printf("charger current operation mode: %s (requested command: %s)", currentMode, cmd)

	if cmd == "start" && (currentMode == model.ChargerOperationModeConnectedCharging || currentMode == model.ChargerOperationModeDisconnected) {
		msg := fmt.Sprintf("charger is in %s mode, not sending start command", currentMode)
		log.Println(msg)
		util.Broadcast(msg)
		return
	}
	if cmd == "stop" && currentMode != model.ChargerOperationModeConnectedCharging {
		msg := fmt.Sprintf("charger is in %s mode (not charging), not sending stop command", currentMode)
		log.Println(msg)
		util.Broadcast(msg)
		return
	}

	util.Broadcast(fmt.Sprintf("charger is in %s mode, sending %s command", currentMode, cmd))

	if jwToken := util.GetToken(); jwToken != "" {
		transport := http.Transport{IdleConnTimeout: time.Second}
		client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
		if req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(conf.StopStartChargingUrl, conf.ChargerId, cmdCode), nil); err != nil {
			log.Printf("failed to create http request: %s\n", err)
		} else {
			maxAttempts := 5
			for attempt := range maxAttempts {
				req.Header = map[string][]string{"Accept": {"*/*"}, "Authorization": {fmt.Sprintf("bearer %s", jwToken)}}
				resp, err := client.Do(req)
				if err == nil && resp != nil {
					respBody, _ := io.ReadAll(resp.Body)
					_ = resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						log.Printf("%s charger succeeded", cmd)
						return
					}
					util.Broadcast(fmt.Sprintf("(attempt %d) failed to %s charger, %d response was returned: %s", attempt, cmd, resp.StatusCode, respBody))
				} else {
					util.Broadcast(fmt.Sprintf("(attempt %d) failed to %s charger: %s", attempt, cmd, err))
				}
				time.Sleep(5 * time.Duration(attempt) * time.Second)
			}
		}
	}
}

func Debug(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	if strings.Contains(update.Message.Text, " on") {
		conf.Bot.Debug = true
		conf.Debug = true
		util.SendMessage(chatId, "debug turned on", true)
	} else {
		if strings.Contains(update.Message.Text, " off") {
			conf.Bot.Debug = false
			conf.Debug = false
			util.SendMessage(chatId, "debug turned off", true)
		} else {
			util.SendMessage(chatId, "please specify /debug on  or  /debug off", false)
		}
	}
}

func ScheduleAdd(update tgbotapi.Update) (schedule model.Schedule) {
	var err error
	chatId := update.Message.Chat.ID
	// first validate/parse the given string, we expect "/sa HH:mm n"
	if schedule, err = util.ParseSchedule(update.Message.Text); err != nil {
		util.SendMessage(chatId, err.Error(), true)
		return
	}
	// next is validating if there is overlap with existing schedules
	overlapExists := false
	for _, chargeSchedule := range conf.ChargeSchedules {
		if !(schedule.StartTime.After(chargeSchedule.StartTime.Add(chargeSchedule.ChargeDuration)) || schedule.StartTime.Add(schedule.ChargeDuration).Before(chargeSchedule.StartTime)) {
			util.SendMessage(chatId, fmt.Sprintf("requested schedule overlaps with existing schedule %s", chargeSchedule.String()), true)
			overlapExists = true
		}
	}
	if !overlapExists {
		conf.ChargeSchedules[schedule.Key()] = schedule
		util.SendMessage(chatId, fmt.Sprintf("charge schedule (%d) %s added", len(conf.ChargeSchedules), schedule.Key()), true)
	}
	return
}

func ScheduleDelete(update tgbotapi.Update) (schedule model.Schedule) {
	var err error
	chatId := update.Message.Chat.ID
	// first validate/parse the given string, we expect "/sd HH:mm n"
	if schedule, err = util.ParseSchedule(update.Message.Text); err != nil {
		util.SendMessage(chatId, err.Error(), true)
		return
	}
	scheduleFound := false
	for _, chargeSchedule := range conf.ChargeSchedules {
		// when deleting a schedule, we don't specify the day (only hours/mins), we should also be able to delete schedules that are already running, so we try both the given time and the given time minus 24 hours
		if (schedule.ChargeDuration == chargeSchedule.ChargeDuration && schedule.StartTime == chargeSchedule.StartTime) || (schedule.ChargeDuration == chargeSchedule.ChargeDuration && schedule.StartTime.Add(-24*time.Hour) == chargeSchedule.StartTime) {
			delete(conf.ChargeSchedules, schedule.Key())
			scheduleFound = true
			break
		}
	}
	if scheduleFound {
		util.SendMessage(chatId, fmt.Sprintf("charge schedule %s deleted, %d schedules left", schedule.Key(), len(conf.ChargeSchedules)), true)
	} else {
		util.SendMessage(chatId, fmt.Sprintf("charge schedule %s not found, %d schedules left", schedule.Key(), len(conf.ChargeSchedules)), true)
	}
	return
}

func ScheduleList(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	if len(conf.ChargeSchedules) == 0 {
		util.SendMessage(chatId, "no charge schedules found", false)
	} else {
		var msg string
		for _, chargeSchedule := range conf.ChargeSchedules {
			msg = fmt.Sprintf("%s%s: %d hours\n", msg, chargeSchedule.StartTime.Format("2006-01-02T15:04Z07:00"), int(chargeSchedule.ChargeDuration.Hours()))
		}
		util.SendMessage(chatId, msg, false)
	}
}

func SetMaxCurrentLow(update tgbotapi.Update) {
	SetMaxCurrent(update, conf.InstallationMaxCurrentLow)
}

func SetMaxCurrentHigh(update tgbotapi.Update) {
	SetMaxCurrent(update, conf.InstallationMaxCurrentHigh)
}

func SetMaxCurrent(update tgbotapi.Update, maxCurrent int) {
	err := InstallationSetMaxCurrent(maxCurrent)
	if err != nil {
		util.SendMessage(update.Message.Chat.ID, fmt.Sprintf("failed to set max current to %d Amps: %s", maxCurrent, err), true)
	} else {
		util.SendMessage(update.Message.Chat.ID, fmt.Sprintf("max current set to %d Amps", maxCurrent), true)
	}
}
