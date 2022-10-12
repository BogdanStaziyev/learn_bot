package main

import (
	"bot_money/config"
	"bot_money/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	var conf = config.GetConfiguration()

	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot)
	if err = telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
