package meme_generator

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"time"
)

type Image struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}

func MemeGenerator(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	urlReddit := "https://meme-api.herokuapp.com/gimme"
	img := new(Image)
	err := getJson(urlReddit, img)
	if err != nil {
		log.Panic(err)
	}
	msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, img.URL)
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
