package command

import (
	"Kibibytes/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type Command struct {
	Name        string
	Description string
	// for simple "arguments -> response" type commands
	// takes precedence over Function if defined
	TextFunction func(args string) string
	// for more complex commands, sending images etc
	Function func(update *tgbotapi.Update, bot *tgbotapi.BotAPI)
}

type Manager struct {
	Commands               []Command
	CommandNotFoundMessage string
}

func (c Command) Execute(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if c.TextFunction != nil {
		utils.SendMessage(c.TextFunction(update.Message.CommandArguments()), update, bot)
	} else {
		c.Function(update, bot)
	}
}

func (cm Manager) GetHelpText() string {
	var text = ""

	for _, cmd := range cm.Commands {
		text = text + fmt.Sprintf("/%s %s\n", cmd.Name, cmd.Description)
	}

	return text
}

func (cm Manager) HandleCommand(cmdName string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	cmdName = strings.ToLower(cmdName)
	for _, cmd := range cm.Commands {
		if cmdName == cmd.Name {
			// command found
			cmd.Execute(update, bot)
			return
		}
	}
	// command not found
	utils.SendMessage(cm.CommandNotFoundMessage, update, bot)
}
