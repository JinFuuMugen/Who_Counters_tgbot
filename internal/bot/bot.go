package bot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"who-counters-bot/internal/dataparse"
	"who-counters-bot/internal/models"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	Herald   int = 0
	Crusader     = 1
	Guardian     = 2
	Archon       = 3
	Legend       = 4
	Ancient      = 5
	Divine       = 6
	Immortal     = 7
)

func GetUserResponse(updates tgbotapi.UpdatesChannel) string {
	responseChan := make(chan string) //goroutine to wait user's response

	go func() {
		for u := range updates {
			if u.Message == nil {
				continue
			}
			if reflect.TypeOf(u.Message.Text).Kind() == reflect.String && u.Message.Text != "" {
				responseChan <- u.Message.Text
				return
			}
		}
	}()

	return <-responseChan //get response from channel
}

func TelegramBot(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	log.Println("Bot started.")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {

			case "/whowon":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Provide match ID:"))

				matchID := GetUserResponse(updates)

				response, err := http.Get("https://api.opendota.com/api/matches/" + matchID)
				if err != nil || response.StatusCode != 200 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Bad match ID."))
					continue
				}

				data, err := io.ReadAll(response.Body)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error."))
					continue
				}

				var allData map[string]interface{}
				err = json.Unmarshal(data, &allData)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error."))
					continue
				}

				if allData["radiant_win"].(bool) {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Radiant claimed victory at game "+matchID))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Dire claimed victory at game "+matchID))
				}

			case "/showhero":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Provide hero ID:"))

				heroID := GetUserResponse(updates)

				response, err := http.Get("https://api.opendota.com/api/heroStats")
				if err != nil || response.StatusCode != http.StatusOK {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Bad hero ID."))
					continue
				}
				heroIDNum, err := strconv.Atoi(heroID)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Bad hero ID."))
					continue
				}

				data, err := io.ReadAll(response.Body)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error."))
					continue
				}

				var allHeroes []models.HeroModel
				err = json.Unmarshal(data, &allHeroes)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error."))
					continue
				}
				flag := false
				for _, hero := range allHeroes {
					if hero.ID == heroIDNum {
						err := dataparse.DownloadPhoto("https://cdn.cloudflare.steamstatic.com"+hero.Image, hero.LocalizedName)
						if err != nil {
							break
						}
						photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "img/"+hero.LocalizedName+".png")
						photo.Caption = "Hero name is " + hero.LocalizedName
						bot.Send(photo)

						winrates := hero.GetWinrates()
						winratesString := "Its winrates: "
						winratesString += "\nAt herald:" + strconv.FormatFloat(100*winrates[Herald], 'f', 2, 64) + "%"
						winratesString += "\nAt crusader: " + strconv.FormatFloat(100*winrates[Crusader], 'f', 2, 64) + "%"
						winratesString += "\nAt guardian: " + strconv.FormatFloat(100*winrates[Guardian], 'f', 2, 64) + "%"
						winratesString += "\nAt archon:" + strconv.FormatFloat(100*winrates[Archon], 'f', 2, 64) + "%"
						winratesString += "\nAt legend:" + strconv.FormatFloat(100*winrates[Legend], 'f', 2, 64) + "%"
						winratesString += "\nAt ancient:" + strconv.FormatFloat(100*winrates[Ancient], 'f', 2, 64) + "%"
						winratesString += "\nAt divine:" + strconv.FormatFloat(100*winrates[Divine], 'f', 2, 64) + "%"
						winratesString += "\nAt immortal:" + strconv.FormatFloat(100*winrates[Immortal], 'f', 2, 64) + "%"

						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, winratesString))

						flag = true
						break
					}
				}
				if !flag {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hero not found."))
				}
			}
		}
	}
	return nil
}
