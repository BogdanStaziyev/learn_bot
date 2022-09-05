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
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

//err = database.Migrate(conf)
//if err != nil {
//	log.Fatalf("Unable to apply migration: %q\n", err)
//}
//
//ses, err := postgresql.Open(
//	postgresql.ConnectionURL{
//		User:     conf.DatabaseUser,
//		Host:     conf.DatabaseHost,
//		Password: conf.DatabasePassword,
//		Database: conf.DatabaseName,
//	})
//if err != nil {
//	log.Fatalf("Unable to create new DB session %q\n: ", err)
//}
//defer func(ses db.Session) {
//	err = ses.Close()
//	if err != nil {
//
//	}
//}(ses)
//
//_, err = os.Stat(conf.FileStorageLocation)
//if err != nil {
//	err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
//}
//if err != nil {
//	log.Fatalf("Storage folder is not available %s", err)
//}
//
//bot.Debug = true
//
//log.Printf("Authorized on account %s", bot.Self.UserName)
