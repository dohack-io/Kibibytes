package googleMapsUrlGenerator

import (
	"Kibibytes/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

var isRunning = false

func SetNotification(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	args := strings.Split(update.Message.CommandArguments(), " ")
	delay, err := strconv.ParseInt(args[0], 10, 64)

	if len(args) == 2 && err == nil {
		delay = delay * int64(60)

		utils.InsertNotify(utils.Notify{
			UserId:        update.Message.Chat.ID,
			Context:       args[1],
			Executiontime: utils.GetUnixtimestamp(delay),
		})

		if !isRunning {
			isRunning = true

			go func() {
				for range time.Tick(time.Minute) {
					if checkNotification(update, bot) {
						isRunning = false
						break
					}
				}
			}()

			utils.SendMessage("notification scheduled.", update, bot)
		}
	} else {
		utils.SendMessage("Missing Parameters", update, bot)
	}
}

func checkNotification(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	notes := utils.GetNextNotifies(update.Message.Chat.ID)

	for _, note := range notes {
		if note.Context == "home" {
			home(&note, update, bot)
		}

		utils.DeleteNotify(note.Id)
	}

	return notes == nil
}

func home(note *utils.Notify, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	utils.SendMessage(ToHome(note.UserId), update, bot)

}
