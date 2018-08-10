package main

import (
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

func main() {

	conf, err := getConfig()
	if err != nil {
		logrus.Error(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/push")
	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
		log.Println("Someone pushed to channel!")
	}
}

// Get method processes env variables and fills Config struct
func getConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("bot", &c); err != nil {
		return c, err
	}
	return c, nil
}
