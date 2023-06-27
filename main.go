package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/cmds"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/util"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	//  used for memory profiling, import net/http/pprof
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	conf.EnvironmentComplete()
	log.SetOutput(os.Stdout)

	var err error

	conf.Bot, err = tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err.Error())
	}

	conf.Bot.Debug = conf.Debug

	conf.Me, err = conf.Bot.GetMe()
	meDetails := "unknown"
	if err == nil {
		meDetails = fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", conf.Me.ID, conf.Me.UserName, conf.Me.FirstName, conf.Me.LastName)
		log.Printf("Bot started: %s, buildtime:%s, commit hash:%s", meDetails, conf.BuildTime, conf.CommitHash)
		log.Printf("Configured chat ids: %v", conf.ChatIDs)
	} else {
		log.Printf("Bot.GetMe() failed: %v", err)
	}

	newUpdate := tgbotapi.NewUpdate(0)
	newUpdate.Timeout = 60

	updatesChan, err := conf.Bot.GetUpdatesChan(newUpdate)
	if err == nil {

		// start schedule-handler
		go func() {
			for range time.Tick(5 * time.Second) {
				for _, schedule := range conf.ChargeSchedules {
					if schedule.StartTime.Before(time.Now()) && schedule.InProgress == false {
						cmds.StartStopCharger("start")
						schedule.InProgress = true
						conf.ChargeSchedules[schedule.Key()] = schedule
						msg := fmt.Sprintf("schedule \"%s\" started", schedule.Key())
						util.Broadcast(msg)
						log.Println(msg)
					}

					if schedule.StartTime.Add(schedule.ChargeDuration).Before(time.Now()) && schedule.InProgress == true {
						cmds.StartStopCharger("stop")
						delete(conf.ChargeSchedules, schedule.Key()) // delete the schedule
						msg := fmt.Sprintf("schedule \"%s\" has ended, %d schedules left", schedule.Key(), len(conf.ChargeSchedules))
						util.Broadcast(msg)
						log.Println(msg)
					}
				}
			}
		}()

		// announce that we are live again
		util.Broadcast(fmt.Sprintf("%s has been (re)started, buildtime: %s", conf.Me.UserName, conf.BuildTime))

		// start listening for messages, and optionally respond
		for update := range updatesChan {
			if update.Message == nil { // ignore any non-Message Updates
				log.Println("ignored null update")
			} else {
				chat := update.Message.Chat
				mentionedMe, cmdMe := TalkOrCmdToMe(update)

				// check if someone is talking to me:
				if (chat.IsPrivate() || (chat.IsGroup() && mentionedMe)) && update.Message.Text != "/start" {
					log.Printf("[%s] [chat:%d] %s\n", update.Message.From.UserName, chat.ID, update.Message.Text)
					if cmdMe {
						fromUser := update.Message.From.UserName
						if chat.IsPrivate() {
							fromUser = chat.UserName
						}
						// /status can be done by anyone, for the other cmds you need admin role
						if util.IsAuthorized(chat.ID) {
							HandleCommand(update)
						} else {
							util.SendMessage(chat.ID, fmt.Sprintf("sorry, %s is not allowed to send me that command", fromUser), true)
						}
					}
				}

			}
		}
	} else {
		log.Printf(
			"failed getting Bot updatesChannel, error: %v", err)
		os.Exit(8)
	}
}

// TalkOrCmdToMe - Returns if we are mentioned and if we were commanded
func TalkOrCmdToMe(update tgbotapi.Update) (bool, bool) {
	entities := update.Message.Entities
	var mentioned = false
	var botCmd = false
	if entities != nil {
		for _, entity := range *entities {
			if entity.Type == "mention" {
				if strings.HasPrefix(update.Message.Text, fmt.Sprintf("@%s", conf.Me.UserName)) {
					mentioned = true
				}
			}
			if entity.Type == "bot_command" {
				botCmd = true
				if strings.Contains(update.Message.Text, fmt.Sprintf("@%s", conf.Me.UserName)) || update.Message.Chat.Type == "private" {
					mentioned = true
				}
			}
		}
	}
	// if another bot was mentioned, the cmd is not for us
	if update.Message.Chat.IsGroup() && mentioned == false {
		botCmd = false
	}
	return mentioned, botCmd
}

func HandleCommand(update tgbotapi.Update) {
	if strings.HasPrefix(update.Message.Text, "/start") {
		util.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Hi %s (%s %s), the bot admin needs to add the chatId first before you can command me", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName), false)
	}

	if strings.HasPrefix(update.Message.Text, "/state") {
		cmds.State(update)
	}

	if strings.HasPrefix(update.Message.Text, "/debug") {
		cmds.Debug(update)
	}

	if strings.HasPrefix(update.Message.Text, "/sa") {
		cmds.ScheduleAdd(update)
	}

	if strings.HasPrefix(update.Message.Text, "/sd") {
		cmds.ScheduleDelete(update)
	}

	if strings.HasPrefix(update.Message.Text, "/sl") {
		cmds.ScheduleList(update)
	}
}
