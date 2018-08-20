package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anatoliyfedorenko/bitbucketbot/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"gopkg.in/telegram-bot-api.v4"
)

//Bot defines bot struct
type Bot struct {
	Chat
	API    *tgbotapi.BotAPI
	config config.Config
}

//NewBot creates new Telegram Bot instance
func NewBot(c config.Config) (*Bot, error) {
	bot := &Bot{}
	a, err := tgbotapi.NewBotAPI(c.TelegramToken)
	if err != nil {
		logrus.Errorf("telegram: NewBotAPI failed: %v\n", err)
		return nil, err
	}
	bot.API = a
	bot.config = c
	return bot, nil
}

//SendUpdate sends message to channel
func (bot *Bot) SendUpdate(text string) {
	m := tgbotapi.NewMessage(bot.config.Chat, text)
	m.DisableWebPagePreview = true
	m.ParseMode = "Markdown"
	bot.API.Send(m)
	logrus.Println("Message Send!")
}

//PullRequestCreated handles PR Created Webhook POST requests
func (bot *Bot) PullRequestCreated(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Created!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%s создал пул реквест [#%v](%v): %v", pr.Actor.DisplayName, pr.PullRequest.ID, pr.PullRequest.Links.HTML.Href, pr.PullRequest.Title)
	logrus.Println(text)
	if (pr.Actor.DisplayName != "") && (pr.PullRequest.ID != 0) && (pr.PullRequest.Links.HTML.Href != "") && (pr.PullRequest.Title != "") {
		bot.SendUpdate(text)
	}

}

//PullRequestCommented handles PR Commented Webhook POST requests
func (bot *Bot) PullRequestCommented(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Commented!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestCommentCreatedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%s написал комментарий к пул реквесту [#%v](%v): %v.", pr.Actor.DisplayName, pr.PullRequest.ID, pr.PullRequest.Links.HTML.Href, pr.Comment.Content.Raw)
	logrus.Println(text)
	if pr.Actor.DisplayName != "" && pr.PullRequest.ID != 0 && pr.PullRequest.Links.HTML.Href != "" && pr.Comment.Content.HTML != "" {
		bot.SendUpdate(text)
	}

}

//PullRequestApproved handles PR Approved Webhook POST requests
func (bot *Bot) PullRequestApproved(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Approved!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestApprovedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%v одобрил пул реквест [#%v](%v)", pr.Approval.User.DisplayName, pr.PullRequest.ID, pr.PullRequest.Links.HTML.Href)
	logrus.Println(text)
	if pr.Approval.User.DisplayName != "" && pr.PullRequest.ID != 0 && pr.PullRequest.Links.HTML.Href != "" {
		bot.SendUpdate(text)
	}

}

//PullRequestMerged handles PR Merged Webhook POST requests
func (bot *Bot) PullRequestMerged(w http.ResponseWriter, r *http.Request) {
	logrus.Println("PR Merged!")
	decoder := json.NewDecoder(r.Body)
	var pr bitbucket.PullRequestMergedPayload
	err := decoder.Decode(&pr)
	if err != nil {
		logrus.Errorf("Decode failed: %v", err)
	}
	text := fmt.Sprintf("%v смержил пул реквест [#%v](%v) в ветку %v", pr.Actor.DisplayName, pr.PullRequest.ID, pr.PullRequest.Links.HTML.Href, pr.PullRequest.Destination.Branch.Name)
	logrus.Println(text)
	if pr.Actor.DisplayName != "" && pr.PullRequest.ID != 0 && pr.PullRequest.Links.HTML.Href != "" && pr.PullRequest.Destination.Branch.Name != "" {
		bot.SendUpdate(text)
	}

}
