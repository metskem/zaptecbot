package cmds

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/types"
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
					stateResponse := types.ChargerStatesRaw{}
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
		util.SendMessage(chatId, "debug turned on", true)
	} else {
		if strings.Contains(update.Message.Text, " off") {
			conf.Bot.Debug = false
			util.SendMessage(chatId, "debug turned off", true)
		} else {
			util.SendMessage(chatId, "please specify /debug on  or  /debug off", false)
		}
	}

}