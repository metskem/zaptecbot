package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/zaptecbot/conf"
	"log"
)

// Broadcast -send message to all admins
func Broadcast(message string) {
	for _, chat := range conf.ChatIDs {
		SendMessage(chat, message)
	}
}

func SendMessage(chatId int64, message string) {
	msgConfig := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: chatId, ReplyToMessageID: 0}, Text: message, DisableWebPagePreview: true}
	_, err := conf.Bot.Send(msgConfig)
	if err != nil {
		log.Printf("failed sending message to chat %d, error is %v", chatId, err)
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
