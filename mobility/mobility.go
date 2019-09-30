package mobility

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"strings"
)

var rest_api_users_url = "http://localhost:8080/api/users"

type GetRequest struct {
	Embedded struct {
		Users []struct {
			Username     string `json:"username"`
			From         string `json:"from"`
			To           string `json:"to"`
			Date         string `json:"date"`
			TrasportType string `json:"trasport_type"`
			TransportID  string `json:"transport_id"`
			Links        struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				User struct {
					Href string `json:"href"`
				} `json:"user"`
			} `json:"_links"`
		} `json:"users"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"self"`
		Profile struct {
			Href string `json:"href"`
		} `json:"profile"`
		Search struct {
			Href string `json:"href"`
		} `json:"search"`
	} `json:"_links"`
	Page struct {
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
		Number        int `json:"number"`
	} `json:"page"`
}

type User struct {
	Username     string `json:"username"`
	From         string `json:"from"`
	To           string `json:"to"`
	Date         string `json:"date"`
	TranspotType string `json:"transport_type"`
	TransportId  string `json:"transport_id"`
}

func Mobility(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi 😉 Specify whether you want to pick up somebody with your ticket using: \n\n"+" /shareRide {username} {startPoint} {endPoint} {date} {transport type} {transport number or identifier} \n\n or to find rideshare opportunities using: \n\n /findRide {startPoint} {endPoint} {date} \n\n 🚌🚎🚃🚝🚊")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func ShareRide(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	arguments := update.Message.CommandArguments()
	s := strings.Split(arguments, " ")
	if len(s) != 6 {
		//
	}
	username, from, to, date, transport_type, transport_id := s[0], s[1], s[2], s[3], s[4], s[5]

	body := &User{
		Username:     username,
		From:         from,
		To:           to,
		Date:         date,
		TranspotType: transport_type,
		TransportId:  transport_id,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	req, _ := http.NewRequest("POST", rest_api_users_url, buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thank you, "+"@"+username+" for your contribution. You'll get contacted soon ☺️")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func FindRide(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	arguments := update.Message.CommandArguments()
	s := strings.Split(arguments, " ")
	from, to, date := s[0], s[1], s[2]
	url_get := rest_api_users_url + "?date=" + date + "&from=" + from + "&to=" + to
	res, _ := http.Get(url_get)
	body, _ := ioutil.ReadAll(res.Body)
	data := GetRequest{}
	json.Unmarshal(body, &data)
	username := data.Embedded.Users[0].Username
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hey, we found something for you, contact  "+"@"+username+" to ride together from "+from+" to "+to+" on "+date+" 🤗")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

}
