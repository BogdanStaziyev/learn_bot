package telegram

import (
	"bot_money/internal/app"
	"bot_money/internal/domain"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	commandStart = "start"
	commandCount = "count"
)

type coordinate map[string]float64

var dataBase = map[int64]coordinate{}

func (b Bot) handelCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Невідома команда")

	switch message.Command() {
	case commandStart:
		dataBase = map[int64]coordinate{}
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Привіт %s. Радий тебе вітати в своєму боті. Тут ти зможеш порахувати зворотню геодезичну задачу не виходячи з телеграму, тобі лише потрібно ввести координати в форматі: X1 252525.252 і так далі, окремими повідомленнями. Я збережу іх та за командою 'Порахувати' зроблю все за тебе", message.Chat.UserName))
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	case commandCount:
		var co domain.Coordinate
		for key, value := range dataBase[message.Chat.ID] {
			switch key {
			case "X1":
				co.X1 = value
			case "X2":
				co.X2 = value
			case "Y1":
				co.Y1 = value
			case "Y2":
				co.Y2 = value
			}
		}
		g, m, s := app.AtaNumber(co.X1, co.Y1, co.X2, co.Y2)
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Координати для перевірки:\nX1: %f Y1: %f \nX2: %f Y2: %f\nРезультат обчислення зворотньої геодезичної задачі %d° %d′ %d″", co.X1, co.Y1, co.X2, co.Y2, g, m, s))
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
	command := strings.Split(message.Text, " ")

	switch command[0] {
	case "X1":
		if len(command) != 2 {
			_, _ = b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error not enough arguments"))
		}
		x1, err := strconv.ParseFloat(command[1], 64)
		if err != nil {
			_, _ = b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error add"))
		}

		if _, ok := dataBase[message.Chat.ID]; !ok {
			dataBase[message.Chat.ID] = coordinate{}
		}
		dataBase[message.Chat.ID][command[0]] = x1

		coordinateX1 := fmt.Sprintf("X1: %f", dataBase[message.Chat.ID][command[0]])
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, coordinateX1))
	case "Y1":
		if len(command) != 2 {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error not enough arguments"))
		}
		y1, err := strconv.ParseFloat(command[1], 64)
		if err != nil {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error add"))
		}

		if _, ok := dataBase[message.Chat.ID]; !ok {
			dataBase[message.Chat.ID] = coordinate{}
		}
		dataBase[message.Chat.ID][command[0]] = y1

		coordinateY1 := fmt.Sprintf("Y1: %f", dataBase[message.Chat.ID][command[0]])
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, coordinateY1))
	case "X2":
		if len(command) != 2 {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error not enough arguments"))
		}
		x2, err := strconv.ParseFloat(command[1], 64)
		if err != nil {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error add"))
		}

		if _, ok := dataBase[message.Chat.ID]; !ok {
			dataBase[message.Chat.ID] = coordinate{}
		}
		dataBase[message.Chat.ID][command[0]] = x2

		coordinateX2 := fmt.Sprintf("X2: %f", dataBase[message.Chat.ID][command[0]])
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, coordinateX2))
	case "Y2":
		if len(command) != 2 {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error not enough arguments"))
		}
		y2, err := strconv.ParseFloat(command[1], 64)
		if err != nil {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "error add"))
		}

		if _, ok := dataBase[message.Chat.ID]; !ok {
			dataBase[message.Chat.ID] = coordinate{}
		}
		dataBase[message.Chat.ID][command[0]] = y2
		coordinateY2 := fmt.Sprintf("Y2: %f", dataBase[message.Chat.ID][command[0]])
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, coordinateY2))
	}
	return nil
}
