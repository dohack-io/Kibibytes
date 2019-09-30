package main

import (
	"Kibibytes/command"
	"Kibibytes/meme_generator"
	"Kibibytes/mobility"
	googleMapsUrlGenerator "Kibibytes/navigation"
	"Kibibytes/twitter"
	"Kibibytes/utils"
	"Kibibytes/utils/secrets"
	"Kibibytes/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	DEFAULT    = iota
	REGISTER   = iota
	SETNAME    = iota
	TRAVELMODE = iota
	READY      = iota
)

var commands = []command.Command{
	{
		Name:        "start",
		Description: "gets you started!",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			utils.SendMessage("Hi there, I'm DohackChatBot coded in Go! /help to get more information. Use /register to fill some default information", update, bot)
			utils.InsertUser(utils.User{
				Id:       update.Message.Chat.ID,
				Username: "",
				Location: "",
				State:    DEFAULT,
			})
		},
	},
	{
		Name:        "register",
		Description: "Save some default informations",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			utils.UpdateUser(utils.User{
				Id:         update.Message.Chat.ID,
				Username:   "",
				Location:   "",
				State:      REGISTER,
				Travelmode: "",
			})
			utils.SendMessage("Welcome, how should we name you?", update, bot)
		},
	},
	{
		Name:        "help",
		Description: "displays this help message",
		TextFunction: func(args string) string {
			return commandManager.GetHelpText()
		},
	},
	{
		Name:        "meme",
		Description: "to get random meme",
		Function:    meme_generator.MemeGenerator,
	},
	{
		Name:        "repeat",
		Description: "what did you say?",
		TextFunction: func(args string) string {
			if args != "" {
				return args
			}

			return "missing parameter"
		},
	},
	{
		Name: "navigate",
		Description: "Generates you a Google Maps URL between 2 points. Please insert to Adresses splitted by an semicolon (;) " +
			"Example: navigate Hamm; Dortmund",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			var msg string

			args := strings.Split(update.Message.CommandArguments(), ";")
			if len(args) == 2 {
				msg = googleMapsUrlGenerator.FromTo(args[0], args[1], update.Message.Chat.ID)
			} else {
				msg = "Please insert to Adresses splitted by an semicolon (;)"
			}

			utils.SendMessage(msg, update, bot)
		},
	},
	{
		Name: "find",
		Description: "Find some location" +
			"Exmaple: find Pizza",
		TextFunction: func(args string) string {
			return googleMapsUrlGenerator.Find(args)
		},
	},
	{
		Name:        "home",
		Description: "Way to home",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			msg := googleMapsUrlGenerator.ToHome(update.Message.Chat.ID)
			utils.SendMessage(msg, update, bot)
		},
	},
	{
		Name:        "tellajoke",
		Description: "tells a random bad joke",
		TextFunction: func(args string) string {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
			req.Header.Add("Accept", "text/plain")
			res, _ := client.Do(req)
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			return string(body)
		},
	},
	{
		Name:         "weathernow",
		Description:  "/weatherNow location-shows weather",
		TextFunction: weather.GetWeatherNow,
	},
	{
		Name:         "weatherforecast",
		Description:  "/weatherForecast location-shows weather forecast",
		TextFunction: weather.GetWeatherForecast,
	},
	{
		Name:         "police",
		Description:  "get last news from Dortmund police department",
		TextFunction: twitter.Command,
	},
	{
		Name:        "notify",
		Description: "schedule a notification",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			googleMapsUrlGenerator.SetNotification(update, bot)
		},
	},
	{
		Name:        "mobility",
		Description: "use ridesharing to ride for free",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			mobility.Mobility(update, bot)
		},
	},
	{
		Name:        "shareride",
		Description: "ridesharing use /mobility first",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			mobility.ShareRide(update, bot)
		},
	},
	{
		Name:        "findride",
		Description: "find ride opportinities use /mobility first",
		Function: func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			mobility.FindRide(update, bot)
		},
	},
}

var commandManager = command.Manager{
	CommandNotFoundMessage: "I don't know this command. I'm not smart enough (I was coded in Go), /help to find out what I can do",
}

// go _really_ hates initialization loops.
func init() {
	commandManager.Commands = commands
}

func main() {
	token := secrets.Get("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token) // main token "906502226:AAGzs1h-h_dZzVpQETgbtim4IySgnBatvCU"
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			commandManager.HandleCommand(update.Message.Command(), &update, bot)
		} else {
			user := utils.GetUser(update.Message.Chat.ID)

			if user.State == REGISTER || user.State == SETNAME || user.State == TRAVELMODE {
				register(user, update.Message.Text, bot, update)
			}
		}
	}
}

func register(user utils.User, text string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var travelmodeKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("walking"),
			tgbotapi.NewKeyboardButton("bicycling"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("transit"),
			tgbotapi.NewKeyboardButton("driving"),
		),
	)

	travelmodeKeyboard.OneTimeKeyboard = true
	travelmodeKeyboard.ResizeKeyboard = true

	if user.State == REGISTER {
		user.Username = text
		user.State = SETNAME

		utils.UpdateUser(user)
		utils.SendMessage(fmt.Sprintf("Hello %s! Where is your place like 127.0.0.1?", user.Username), &update, bot)

	} else if user.State == SETNAME {
		user.Location = text
		user.State = TRAVELMODE

		utils.UpdateUser(user)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What type of traveler are you mostly?")
		msg.ReplyMarkup = travelmodeKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

	} else if user.State == TRAVELMODE {
		user.Travelmode = text
		user.State = READY
		utils.UpdateUser(user)
		utils.SendMessage("That's it have fun!", &update, bot)
	}
}
