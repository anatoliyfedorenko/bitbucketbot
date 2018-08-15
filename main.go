package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
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

	http.HandleFunc("/merge_created", bot.mergeCreated)
	http.HandleFunc("/merge_commented", bot.mergeCommented)
	http.HandleFunc("/merge_approved", bot.mergeApproved)
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

func (bot Bot) mergeCreated(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пользователь %s создал пул реквест! Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href)
	bot.sendUpdate(text)
}

func (bot Bot) mergeCommented(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCommentCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%s написал комментарий к пул реквесту (%v). Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href, pr.Comment.Links.HTML.Href)
	bot.sendUpdate(text)
}

func (bot Bot) mergeApproved(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestApprovedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пул реквест был одобрен %v! Посмотреть => %v", pr.Approval.User.DisplayName, pr.PullRequest.Links.HTML.Href)
	bot.sendUpdate(text)
}

func (bot Bot) mergeAccepted(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestMergedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пул реквест был мержнут пользователем %v! Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href)
	bot.sendUpdate(text)
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	m.DisableWebPagePreview = true
	bot.API.Send(m)
}
