package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

//Config defines config struct
type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	Chat          int64  `envconfig:"CHAT" required:"true"`
}

type Bot struct {
	API *tgbotapi.BotAPI
	c   Config
}

func main() {

	conf, err := getConfig()
	if err != nil {
		logrus.Error(err)
	}

	bot := &Bot{}

	b, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.API = b
	bot.c = conf

	log.Printf("Authorized on account %s", bot.API.Self.UserName)

	http.HandleFunc("/push", bot.push)
	http.ListenAndServe(":8080", nil)
}

// Get method processes env variables and fills Config struct
func getConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("bot", &c); err != nil {
		return c, err
	}
	return c, nil
}

func (bot Bot) push(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Someone just pushed to repo!")
	m := tgbotapi.NewMessage(bot.c.Chat, "Someone just pushed to repo!")
	bot.API.Send(m)
}
