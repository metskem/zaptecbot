package util

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang-jwt/jwt/v4"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Broadcast -send message to all admins
func Broadcast(message string) {
	for _, chat := range conf.ChatIDs {
		SendMessage(chat, message, false)
	}
}

func SendMessage(chatId int64, message string, logToStdout bool) {
	msgConfig := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: chatId, ReplyToMessageID: 0}, Text: message, DisableWebPagePreview: true}
	_, err := conf.Bot.Send(msgConfig)
	if err != nil {
		log.Printf("failed sending message to chat %d, error is %v", chatId, err)
	}
	if logToStdout {
		log.Println(message)
	}
}

func IsAuthorized(chatId int64) bool {
	for _, allowedChatId := range conf.ChatIDs {
		if allowedChatId == chatId {
			return true
		}
	}
	return false
}

func IsTokenValid(token string) bool {
	if token == "" {
		return false
	}
	//fmt.Println("validating token ", token)
	jwToken, _ := jwt.Parse(token, nil)
	return jwToken.Valid
}

func GetToken() string {
	if !IsTokenValid(conf.CachedToken) {
		transport := http.Transport{IdleConnTimeout: time.Second}
		client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
		userEscaped := url.QueryEscape(conf.UserName)
		passwordEscaped := url.QueryEscape(conf.Password)
		postData := fmt.Sprintf("grant_type=password&username=%s&password=%s", userEscaped, passwordEscaped)
		resp, err := client.Post(conf.GetTokenUrl, "application/x-www-form-urlencoded", strings.NewReader(postData))
		if err == nil && resp != nil {
			respBody, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusOK {
				loginResponse := model.LoginResponse{}
				if err := json.Unmarshal(respBody, &loginResponse); err != nil {
					log.Printf("failed to decode the login response: %s\n", err)
				}
				conf.CachedToken = loginResponse.AccessToken
				log.Printf("succesfull login, token will expire in %d hours\n", loginResponse.ExpiresIn/3600)
				_ = resp.Body.Close()
				return loginResponse.AccessToken
			} else {
				log.Printf("response (%d) from login failed:%s\n", resp.StatusCode, respBody)
			}
		} else {
			log.Printf("response from login failed:%s\n", err)
		}
	}
	return ""
}

func ParseChargerState(rawStates model.ChargerStatesRaw) model.ChargerState {
	chargerState := model.ChargerState{}
	for _, rawState := range rawStates {
		switch rawState.StateID {
		case 150:
			chargerState.CommunicationMode = rawState.ValueAsString
		case 151:
			chargerState.PermanentCableLock = rawState.ValueAsString
		case 201:
			chargerState.TemperatureInternal5 = rawState.ValueAsString
		case 270:
			chargerState.Humidity = rawState.ValueAsString
		case 501:
			chargerState.VoltagePhase1 = rawState.ValueAsString
		case 502:
			chargerState.VoltagePhase2 = rawState.ValueAsString
		case 503:
			chargerState.VoltagePhase3 = rawState.ValueAsString
		case 507:
			chargerState.CurrentPhase1 = rawState.ValueAsString
		case 508:
			chargerState.CurrentPhase2 = rawState.ValueAsString
		case 509:
			chargerState.CurrentPhase3 = rawState.ValueAsString
		case 548:
			chargerState.PhaseRotation = rawState.ValueAsString
		case 702:
			chargerState.ChargeMode = rawState.ValueAsString
		case 710:
			switch rawState.ValueAsString {
			case "0":
				chargerState.ChargerOperationMode = model.ChargerOperationMode0
			case "1":
				chargerState.ChargerOperationMode = model.ChargerOperationMode1
			case "2":
				chargerState.ChargerOperationMode = model.ChargerOperationMode2
			case "3":
				chargerState.ChargerOperationMode = model.ChargerOperationMode3
			case "5":
				chargerState.ChargerOperationMode = model.ChargerOperationMode5
			}
		}
	}
	return chargerState
}

func ParseSchedule(updateText string) (schedule model.Schedule, err error) {
	var durationStr, parsedTime string
	words := strings.Split(updateText, " ")
	now := time.Now()
	schedRegex := regexp.MustCompile(conf.SchedulePattern1)
	if !schedRegex.MatchString(updateText) {
		schedRegex = regexp.MustCompile(conf.SchedulePattern2)
		if !schedRegex.MatchString(updateText) {
			return schedule, errors.New(fmt.Sprintf("failed to parse schedule %s", updateText))
		} else {
			parsedTime = words[1] + " " + words[2]
			durationStr = words[3]
		}
	} else {
		//we add the current year/month/day
		parsedTime = fmt.Sprintf("%d-%d-%d %s", now.Year(), now.Month(), now.Day(), words[1])
		durationStr = words[2]
	}

	if len(words) != 3 && len(words) != 4 {
		return schedule, errors.New(fmt.Sprintf("failed to parse request, we expected 3 or 4 words, but got %d", len(words)))
	}

	if schedule.StartTime, err = time.ParseInLocation("2006-1-2 15:04", parsedTime, now.Location()); err != nil {
		return schedule, errors.New(fmt.Sprintf("failed to parse time %s: %s", words[1], err))
	}

	// when you request a time that is less than current time, we assume that you meant that time for tomorrow:
	if schedule.StartTime.Before(now) {
		schedule.StartTime = schedule.StartTime.Add(time.Hour * 24)
	}

	var duration int
	if duration, err = strconv.Atoi(durationStr); err != nil {
		return schedule, errors.New(fmt.Sprintf("failed to parse duration %s: %s", durationStr, err))
	}
	if duration <= 0 || duration > 24 {
		return schedule, errors.New(fmt.Sprintf("duration %d is <0 or >24 hours", duration))
	}
	schedule.ChargeDuration = time.Duration(duration) * time.Hour
	return
}
