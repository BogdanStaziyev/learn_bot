package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b Bot) Start() error {

	updates := b.initUpdatesChanel()

	b.handleUpdates(updates)

	return nil
}

func (b Bot) initUpdatesChanel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // non massage
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handelCommand(update.Message); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if err := b.handelMessage(update.Message); err != nil {
			log.Fatal(err)
		}
	}
}

func (b Bot) sentCommand(text string) error {
	fmt.Print(text)
	return nil
}
