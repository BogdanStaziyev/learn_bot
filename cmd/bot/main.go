package main

import (
	"bot_money/config"
	"bot_money/internal/database"
	"encoding/json"
	"fmt"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type wallet map[string]float64

var dataBase = map[int64]wallet{}

func main() {
	var conf = config.GetConfiguration()

	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	err = database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migration: %q\n", err)
	}

	ses, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session %q\n: ", err)
	}
	defer func(ses db.Session) {
		err = ses.Close()
		if err != nil {

		}
	}(ses)

	_, err = os.Stat(conf.FileStorageLocation)
	if err != nil {
		err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
	}
	if err != nil {
		log.Fatalf("Storage folder is not available %s", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			command := strings.Split(update.Message.Text, " ")

			switch command[0] {
			case "Add":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, " money"))
				}
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = wallet{}
				}
				dataBase[update.Message.Chat.ID][command[1]] += amount

				balance := fmt.Sprintf("%f", dataBase[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balance))
			case "Sub":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error money"))
				}
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					continue
				}
				dataBase[update.Message.Chat.ID][command[1]] -= amount

				balance := fmt.Sprintf("%f", dataBase[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balance))
			case "Delete":
				delete(dataBase[update.Message.Chat.ID], command[1])
			case "Show":
				msg := ""
				var sum float64
				for key, value := range dataBase[update.Message.Chat.ID] {
					price, _ := getPrice(key)
					valPrice := price * value
					sum += valPrice
					msg += fmt.Sprintf("%s : %f, %f \n", key, value, valPrice)
				}

				msg += fmt.Sprintf("total price: %f", sum)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
			}
		}
	}
}

type Binance struct {
	Price float64 `json:"price,string"`
}

func getPrice(symbol string) (price float64, err error) {
	log.Print(symbol)
	resp, _ := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/avgPrice?symbol=%sUSDT", symbol))
	defer resp.Body.Close()

	var jResp Binance

	err = json.NewDecoder(resp.Body).Decode(&jResp)
	if err != nil {
		return
	}
	return jResp.Price, nil
}
