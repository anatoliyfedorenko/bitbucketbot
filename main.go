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

	http.HandleFunc("/pull_request_created", bot.pullRequestCreated)
	http.HandleFunc("/pull_request_commented", bot.pullRequestCommented)
	http.HandleFunc("/pull_request_approved", bot.pullRequestApproved)
	http.HandleFunc("/pull_request_merged", bot.pullRequestMerged)
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

func (bot Bot) pullRequestCreated(w http.ResponseWriter, r *http.Request) {
	log.Println("PR Created!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пользователь %s создал пул реквест! Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href)
	log.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) pullRequestCommented(w http.ResponseWriter, r *http.Request) {
	log.Println("PR Commented!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCommentCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%s написал комментарий к пул реквесту (%v). Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href, pr.Comment.Links.HTML.Href)
	log.Println(text)
	bot.sendUpdate(text)
}
log.Println("PR Merged!")

func (bot Bot) pullRequestApproved(w http.ResponseWriter, r *http.Request) {
	log.Println("PR Approved!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestApprovedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пул реквест был одобрен %v! Посмотреть => %v", pr.Approval.User.DisplayName, pr.PullRequest.Links.HTML.Href)
	log.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) pullRequestMerged(w http.ResponseWriter, r *http.Request) {
	log.Println("PR Merged!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestpullRequestdPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("Пул реквест был мержнут пользователем %v! Посмотреть => %v", pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href)
	log.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	m.DisableWebPagePreview = true
	bot.API.Send(m)
	log.Println("Message Send!")
}
