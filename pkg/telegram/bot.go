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

//
//func atanNumber(x1, y1, x2, y2 float64) (int, int, int) {
//	const radius float64 = 180
//	const degree, minutes, seconds int = 180, 60, 60
//	var deg, min, sec int
//	x := x2 - x1
//	y := y2 - y1
//	subtractionCoordinate := y / x
//	atanResult := math.Atan(subtractionCoordinate)
//	atanResult *= radius / math.Pi
//	deg = int(atanResult)
//	minute := (atanResult - float64(deg)) * 60
//	min = int(minute)
//	sec = int((minute - float64(min)) * 60)
//	if x < 0 && y > 0 {
//		deg = (degree - 1) + deg
//		min = (minutes - 1) + min
//		sec = seconds + sec
//	} else if x < 0 && y < 0 {
//		deg = degree + deg
//	} else if x > 0 && y < 0 {
//		deg = ((degree * 2) - 1) + deg
//		min = (minutes - 1) + min
//		sec = seconds + sec
//	}
//	return deg, min, sec
//}
