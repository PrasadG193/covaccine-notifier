package notify

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Telegram struct {
	ChatID int64
	Bot    *tgbotapi.BotAPI
}

func NewTelegram(token, username string) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find bot for the given botAPI token %s", token))
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdates(u)
	if err != nil {
		return nil, errors.New("Unable to get channel for bot updates")
	}
	for _, update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if strings.ToLower(update.Message.Chat.UserName) == strings.ToLower(username) {
			chatID := update.Message.Chat.ID
			log.Printf("chatID for the conversation between username: %s and bot: %s is %d", username, bot.Self.UserName, chatID)
			return &Telegram{
				ChatID: chatID,
				Bot:    bot,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to get the chatID \n Send message to the bot %s", bot.Self.UserName))
}

func (t *Telegram) SendMessage(body string) error {
	if len(body) > 4096 {
		log.Printf("message body too long, trimming it")
		body = body[:4095]
	}
	msg := tgbotapi.NewMessage(t.ChatID, body)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to send message to telegram")
	}
	return nil
}
