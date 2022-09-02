package main

import (
	"bot_money/config"
	"bot_money/internal/database"
	"fmt"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type coordinate map[string]float64

var dataBase = map[int64]coordinate{}

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

	type coord struct {
		X1 float64
		Y1 float64
		X2 float64
		Y2 float64
	}

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			command := strings.Split(update.Message.Text, " ")

			var upname = update.Message.Chat.UserName

			var HelloMessage = fmt.Sprintf("Привіт %s. Радий тебе вітати в своєму боті. Тут ти зможеш порахувати зворотню геодезичну задачу не виходячи з телеграму, тобі лише потрібно ввести координати в форматі: X1 252525.252 і так далі, окремими повідомленнями. Я збережу іх та за командою 'Порахувати' зроблю все за тебе", upname)

			switch command[0] {
			case "/start":
				if len(command) != 1 {
					_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error not found command"))
					continue
				}
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, HelloMessage))

			case "X1":
				if len(command) != 2 {
					_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error not enough arguments"))
					continue
				}
				x1, err := strconv.ParseFloat(command[1], 64)
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
					continue
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = coordinate{}
				}
				dataBase[update.Message.Chat.ID][command[0]] = x1

				coordinateX1 := fmt.Sprintf("X1: %f", dataBase[update.Message.Chat.ID][command[0]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, coordinateX1))
			case "Y1":
				if len(command) != 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error not enough arguments"))
					continue
				}
				y1, err := strconv.ParseFloat(command[1], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
					continue
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = coordinate{}
				}
				dataBase[update.Message.Chat.ID][command[0]] = y1

				coordinateY1 := fmt.Sprintf("Y1: %f", dataBase[update.Message.Chat.ID][command[0]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, coordinateY1))
			case "X2":
				if len(command) != 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error not enough arguments"))
					continue
				}
				x2, err := strconv.ParseFloat(command[1], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
					continue
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = coordinate{}
				}
				dataBase[update.Message.Chat.ID][command[0]] = x2

				coordinateX2 := fmt.Sprintf("X2: %f", dataBase[update.Message.Chat.ID][command[0]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, coordinateX2))
			case "Y2":
				if len(command) != 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error not enough arguments"))
					continue
				}
				y2, err := strconv.ParseFloat(command[1], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error add"))
					continue
				}

				if _, ok := dataBase[update.Message.Chat.ID]; !ok {
					dataBase[update.Message.Chat.ID] = coordinate{}
				}
				dataBase[update.Message.Chat.ID][command[0]] = y2
				coordinateY2 := fmt.Sprintf("Y2: %f", dataBase[update.Message.Chat.ID][command[0]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, coordinateY2))
			case "Порахувати":
				var co coord
				msg := ""
				for key, value := range dataBase[update.Message.Chat.ID] {
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

				g, m, s := atanNumber(co.X1, co.Y1, co.X2, co.Y2)
				msg += fmt.Sprintf("Координати для перевірки:\nX1: %f Y1: %f \nX2: %f Y2: %f\nРезультат обчислення зворотньої геодезичної задачі %d° %d′ %d″", co.X1, co.Y1, co.X2, co.Y2, g, m, s)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
			}
		}
	}
}

func atanNumber(x1, y1, x2, y2 float64) (int, int, int) {
	const radius float64 = 180
	const degree, minutes, seconds int = 180, 60, 60
	var deg, min, sec int
	x := x2 - x1
	y := y2 - y1
	subtractionCoordinate := y / x
	atanResult := math.Atan(subtractionCoordinate)
	atanResult *= radius / math.Pi
	deg = int(atanResult)
	minute := (atanResult - float64(deg)) * 60
	min = int(minute)
	sec = int((minute - float64(min)) * 60)
	if x < 0 && y > 0 {
		deg = (degree - 1) + deg
		min = (minutes - 1) + min
		sec = seconds + sec
	} else if x < 0 && y < 0 {
		deg = degree + deg
	} else if x > 0 && y < 0 {
		deg = ((degree * 2) - 1) + deg
		min = (minutes - 1) + min
		sec = seconds + sec
	}
	return deg, min, sec
}
