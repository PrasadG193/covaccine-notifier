//Package notify has functions and types used for sending notifications on different communication channel
package notify

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

const timeout = 60

type Telegram struct {
	ChatID int64
	Bot    *tgbotapi.BotAPI
}

// NewTelegram returns new instance of Telegram service that has new BotAPI instance and chatID.
//
// It requires a username to fetch the appropriate chatID and
// a token provided by @BotFather on Telegram to create bot Instance
func NewTelegram(username, token string) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find bot for the given botAPI token %s", token))
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

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
	return nil, errors.New(fmt.Sprintf("Unable to get the chatID, Send message to the bot %s", bot.Self.UserName))
}

// SendMessage takes message body and send it to the given chatID as text message or file
func (t *Telegram) SendMessage(body string) error {
	if len(body) > 4096 {
		log.Printf("Message body too long, Message will be sent as file ")
		fileBytes := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("slots-available-%d.txt", time.Now().Unix()),
			Bytes: []byte(body),
		}
		documentConfig := tgbotapi.NewDocumentUpload(t.ChatID, fileBytes)
		if _, err := t.Bot.Send(documentConfig); err != nil {
			return errors.Wrap(err, "Unable to send message to telegram")
		}
	} else {
		msg := tgbotapi.NewMessage(t.ChatID, body)
		_, err := t.Bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "Unable to send message to telegram")
		}
	}
	return nil
}
