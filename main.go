package main

import (
	"encoding/json"
	"fmt"
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

	User struct {
		userType    string `json: type`
		userName    string `json: username`
		displayName string `json: display_name`
	}

	Push struct {
		user       string `json:"actor"`
		repository string `json:"repository"`
	}

	PullRequest struct {
		id          int64  `json:"id"`
		title       string `json:"title"`
		description string `json:"description"`
	}

	MergeCreated struct {
		owner       User        `json:"actor"`
		pullRequest PullRequest `json:"pullrequest"`
	}

	MergeAccepted struct {
		owner       User        `json:"actor"`
		pullRequest PullRequest `json:"pullrequest"`
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
	fmt.Println("New push to repo, begin decoding...")

	decoder := json.NewDecoder(r.Body)
	var p Push
	err := decoder.Decode(&p)
	if err != nil {
		log.Println("Error decoding!")
	}

	text := fmt.Sprintf("User: %v pushed to repo: %v", p.user, p.repository)

	bot.sendUpdate(text)
}

func (bot Bot) mergeCreated(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MR to repo!")
	bot.sendUpdate("MR to repo!")
}

func (bot Bot) mergeAccepted(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MR accepted!")
	bot.sendUpdate("MR accepted!")
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	bot.API.Send(m)
}
