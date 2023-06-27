package cmds

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/model"
	"github.com/metskem/zaptecbot/util"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func State(update tgbotapi.Update) {
	//msg := fmt.Sprintf("list requested by %s", update.Message.From.UserName)
	//log.Println(msg)
	//util.SendMessage(update.Message.Chat.ID, msg, false)
	if jwToken := util.GetToken(); jwToken != "" {
		transport := http.Transport{IdleConnTimeout: time.Second}
		client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
		if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(conf.StateUrl, conf.ChargerId), nil); err != nil {
			log.Printf("failed to create http request: %s\n", err)
		} else {
			req.Header = map[string][]string{"Accept": {"*/*"}, "Authorization": {fmt.Sprintf("bearer %s", jwToken)}}
			resp, err := client.Do(req)
			if err == nil && resp != nil {
				respBody, _ := io.ReadAll(resp.Body)
				if resp.StatusCode == http.StatusOK {
					stateResponse := model.ChargerStatesRaw{}
					if err := json.Unmarshal(respBody, &stateResponse); err != nil {
						log.Printf("failed to decode the charger state response: %s\n", err)
					}
					chargerState := util.ParseChargerState(stateResponse)
					util.SendMessage(update.Message.Chat.ID, fmt.Sprintf("%s", chargerState), false)
					_ = resp.Body.Close()
				} else {
					log.Printf("response (%d) from charge state failed:%s\n", resp.StatusCode, respBody)
				}
			} else {
				log.Printf("response from charger state failed:%s\n", err)
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
	if schedule, err = util.ValidateSchedule(update.Message.Text); err != nil {
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
	if schedule, err = util.ValidateSchedule(update.Message.Text); err != nil {
		util.SendMessage(chatId, err.Error(), true)
		return
	}
	scheduleFound := false
	for _, chargeSchedule := range conf.ChargeSchedules {
		if schedule.ChargeDuration == chargeSchedule.ChargeDuration && schedule.StartTime == chargeSchedule.StartTime {
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
			msg = fmt.Sprintf("%sstartTime: %s, duration: %d hours\n", msg, chargeSchedule.StartTime.Format("15:04"), int(chargeSchedule.ChargeDuration.Hours()))
		}
		util.SendMessage(chatId, msg, false)
	}
}
