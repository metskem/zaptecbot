package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/conf"
	"github.com/metskem/zaptecbot/util"
	"log"
	"os"
	"strings"
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
							util.SendMessage(chat.ID, fmt.Sprintf("sorry, %s is not allowed to send me that command", fromUser))
						}
					}
				}

			}
			fmt.Println("")
		}
	} else {
		log.Printf("failed getting Bot updatesChannel, error: %v", err)
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
				if strings.Contains(update.Message.Text, fmt.Sprintf("@%s", conf.Me.UserName)) {
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
	chatId := update.Message.Chat.ID
	if strings.HasPrefix(update.Message.Text, "/list") {
		msg := fmt.Sprintf("list requested by %s", update.Message.From.UserName)
		log.Println(msg)
		util.SendMessage(chatId, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/start") {
		util.SendMessage(chatId, fmt.Sprintf("Hi %s (%s %s), the bot admin needs to add the chatId first before you can command me", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName))
	}

	if strings.HasPrefix(update.Message.Text, "/debug") {
		if strings.Contains(update.Message.Text, " on") {
			conf.Bot.Debug = true
			util.SendMessage(chatId, "debug turned on")
		} else {
			if strings.Contains(update.Message.Text, " off") {
				conf.Bot.Debug = false
				util.SendMessage(chatId, "debug turned off")
			} else {
				util.SendMessage(chatId, "please specify /debug on  or  /debug off")
			}
		}
	}

}
