package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	commandStart = "start"
)

func (b Bot) handelCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Невідома команда")

	switch message.Command() {
	case commandStart:
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Привіт %s. Радий тебе вітати в своєму боті. Тут ти зможеш порахувати зворотню геодезичну задачу не виходячи з телеграму, тобі лише потрібно ввести координати в форматі: X1 252525.252 і так далі, окремими повідомленнями. Я збережу іх та за командою 'Порахувати' зроблю все за тебе", message.Chat.UserName))
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	default:
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (b Bot) handelMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}
