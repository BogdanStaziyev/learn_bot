package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type wallet map[string]float64

var db = map[int64]wallet{}

func main() {

	bot, err := tgbotapi.NewBotAPI("5407561453:AAE3topXqyAy-zZ8eIJkD0uZkNBcJlaXfUI")
	if err != nil {
		log.Panic(err)
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

				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}
				db[update.Message.Chat.ID][command[1]] += amount

				balance := fmt.Sprintf("%f", db[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balance))
			case "Sub":
				if len(command) != 3 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error money"))
				}
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
				}

				if _, ok := db[update.Message.Chat.ID]; !ok {
					continue
				}
				db[update.Message.Chat.ID][command[1]] -= amount

				balance := fmt.Sprintf("%f", db[update.Message.Chat.ID][command[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balance))
			case "Delete":
				delete(db[update.Message.Chat.ID], command[1])
			case "Show":
				msg := ""
				var sum float64
				for key, value := range db[update.Message.Chat.ID] {
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
