package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

type (
	//Config defines config struct
	Config struct {
		TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
		Chat          int64  `envconfig:"CHAT" required:"true"`
	}

	//Bot defines bot struct
	Bot struct {
		API *tgbotapi.BotAPI
		c   Config
	}
)

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
	http.HandleFunc("/merge_created", bot.mergeCreated)
	http.HandleFunc("/merge_accepted", bot.mergeAccepted)
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
	log.Println("New push created!")
}

func (bot Bot) mergeCreated(w http.ResponseWriter, r *http.Request) {
	bot.sendUpdate("New PR created! Here is the data: \n")
	bot.sendUpdate(getResponseString(r))
}

func (bot Bot) mergeAccepted(w http.ResponseWriter, r *http.Request) {
	bot.sendUpdate("PR merged! Here is the data: \n")
	bot.sendUpdate(getResponseString(r))
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	bot.API.Send(m)
}

func getResponseString(r *http.Request) string {
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(responseData)
}
