package main

import (
	"net/http"

	"github.com/anatoliyfedorenko/bitbucketbot/chat"
	"github.com/anatoliyfedorenko/bitbucketbot/config"
	"github.com/sirupsen/logrus"
)

func main() {

	c, err := config.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	bot, err := chat.NewBot(c)

	logrus.Printf("Authorized on account %s", bot.API.Self.UserName)
	logrus.Printf("Configure to send messages to chat: %v", c.Chat)

	http.HandleFunc("/pull_request_created", bot.PullRequestCreated)
	http.HandleFunc("/pull_request_commented", bot.PullRequestCommented)
	http.HandleFunc("/pull_request_approved", bot.PullRequestApproved)
	http.HandleFunc("/pull_request_merged", bot.PullRequestMerged)
	http.HandleFunc("/pull_request_declined", bot.PullRequestDeclined)

	http.ListenAndServe(":8080", nil)
}
