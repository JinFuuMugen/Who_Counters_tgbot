package bot

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func TelegramBot(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {

			case "/whowon":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Provide match ID:"))

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

				matchID := <-responseChan //get mathId from channel

				response, err := http.Get("https://api.opendota.com/api/matches/" + matchID)
				if err != nil {
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

			}
		}
	}
	return nil
}
