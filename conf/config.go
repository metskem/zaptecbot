package conf

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/model"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	SchedulePattern1 = "^/s[ad] \\d{2}:\\d{2} \\d{1,2}.*"                          // /s[ad] 12:34 5
	SchedulePattern2 = "^/s[ad] \\d{4}\\-\\d{2}\\-\\d{2} \\d{2}:\\d{2} \\d{1,2}.*" // /s[ad] 2021-01-01 12:34 5
)

var (
	CommitHash string
	BuildTime  string

	BotToken   = os.Getenv("BOT_TOKEN")
	ChargerId  = os.Getenv("ZAPTEC_CHARGER_ID")
	UserName   = os.Getenv("ZAPTEC_USERNAME")
	Password   = os.Getenv("ZAPTEC_PASSWORD")
	DebugStr   = os.Getenv("DEBUG")
	ChatIDsStr = os.Getenv("CHAT_IDS")
	ChatIDs    = make(map[int]int64)
	Debug      bool

	Me                   tgbotapi.User
	Bot                  *tgbotapi.BotAPI
	HttpTimeout          = 5
	GetTokenUrl          = "https://api.zaptec.com/oauth/token"
	StateUrl             = "https://api.zaptec.com/api/chargers/%s/state"
	StopStartChargingUrl = "https://api.zaptec.com/api/chargers/%s/sendCommand/%d"
	CachedToken          string
	ChargeSchedules      = make(map[string]model.Schedule) // key is the schedule time+ duration, i.e.: "12:34 5"
)

func EnvironmentComplete() {
	envComplete := true

	if len(BotToken) == 0 {
		log.Print("missing envvar \"BOT_TOKEN\"")
		envComplete = false
	}
	if len(ChargerId) == 0 {
		log.Print("missing envvar \"ZAPTEC_CHARGER_ID\"")
		envComplete = false
	}
	if len(UserName) == 0 {
		log.Print("missing envvar \"ZAPTEC_USERNAME\"")
		envComplete = false
	}
	if len(Password) == 0 {
		log.Print("missing envvar \"ZAPTEC_PASSWORD\"")
		envComplete = false
	}

	Debug = false
	if DebugStr == "true" {
		Debug = true
	}

	chatIDsString := strings.Split(ChatIDsStr, ",")
	var chatids string
	for i := 0; i < len(chatIDsString); i++ {
		ChatIDs[i], _ = strconv.ParseInt(chatIDsString[i], 0, 64)
		chatids = fmt.Sprintf("%s %d", chatids, ChatIDs[i])
	}

	if !envComplete {
		log.Fatal("one or more required envvars missing, aborting...")
	}
}
