package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gen1us2k/go-translit"
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
		logrus.Fatal(err)
	}
	bot.API = b
	bot.c = conf

	logrus.Printf("Authorized on account %s", bot.API.Self.UserName)
	logrus.Printf("Configure to send messages to chat: %v", conf.Chat)

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
	logrus.Println("PR Created!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%s создал пул реквест: [%v](%v)", translit.Translit(pr.Actor.DisplayName), pr.PullRequest.Title, pr.PullRequest.Links.HTML.Href)
	logrus.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) pullRequestCommented(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Commented!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCommentCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf(`%s написал комментарий к пул реквесту "[%v](%v)" : %v.`, translit.Translit(pr.Actor.DisplayName), pr.PullRequest.Title, pr.PullRequest.Links.HTML.Href, pr.Comment.Content.Markup)
	logrus.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) pullRequestApproved(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Approved!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestApprovedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf(`%v одобрил ПР "[%v](%v)"`, translit.Translit(pr.Approval.User.DisplayName), pr.PullRequest.Title, pr.PullRequest.Links.HTML.Href)
	logrus.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) pullRequestMerged(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Merged!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestMergedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf(`%v смержил ПР "[%](%v)" в ветку %v`, translit.Translit(pr.Actor.DisplayName), pr.PullRequest.Title, pr.PullRequest.Links.HTML.Href, pr.PullRequest.Destination.Branch.Name)
	logrus.Println(text)
	bot.sendUpdate(text)
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	m.DisableWebPagePreview = true
	bot.API.Send(m)
	logrus.Println("Message Send!")
}
