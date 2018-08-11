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
	bot.sendUpdate("New PR created! \n")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	logrus.Infof("Full Info: %v", pr)
	text := fmt.Sprintf("There is a new pull request in %v created by %v! Please, click [here](%v) to view details!", pr.Repository.FullName, pr.Actor.DisplayName, pr.PullRequest.Links.HTML.Href)
	bot.sendUpdate(text)

	prDetails := "Pull Request data: \n"
	prDetails += fmt.Sprintf("title: %v \n", pr.PullRequest.Title)
	prDetails += fmt.Sprintf("description: %v \n", pr.PullRequest.Description)
	prDetails += fmt.Sprintf("from %v to %v \n", pr.PullRequest.Source.Branch.Name, pr.PullRequest.Destination.Branch.Name)
	prDetails += fmt.Sprintf("Reviews needed from: %v \n", pr.PullRequest.Reviewers)

	bot.sendUpdate(prDetails)

}

func (bot Bot) mergeAccepted(w http.ResponseWriter, r *http.Request) {
	bot.sendUpdate("PR merged! \n")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestMergedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	logrus.Infof("Full Info: %v", pr)
	bot.sendUpdate(pr.PullRequest.Title)
}

func (bot Bot) sendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.c.Chat, text)
	bot.API.Send(m)
}
