package chat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anatoliyfedorenko/bitbucketbot/config"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// Just change these constants to test bot in real chat.
// Make sure you add a bot to the chat
const BotToken = "638667863:AAFLWagdVgmNJFbwo_nZ-paKFMBAQj9Fo74"
const BotChat = "-275166411"

func TestSendUpdate(t *testing.T) {
	bot := setupTestBot(t)
	bot.SendUpdate("test")
}

func TestNewBotReturnsError(t *testing.T) {
	os.Setenv("BOT_TELEGRAM_TOKEN", BotToken)
	os.Setenv("BOT_CHAT", BotChat)
	conf, err := config.GetConfig()
	assert.NoError(t, err)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	r := httpmock.NewStringResponder(200, `{
		"ok": false,
		"error_code": 404,
		"description": "Not Found"
	  }`)

	url := fmt.Sprintf("https://api.telegram.org/bot%v/getMe", conf.TelegramToken)
	httpmock.RegisterResponder("POST", url, r)

	_, err = NewBot(conf)
	assert.Error(t, err)
	assert.Equal(t, "Not Found", err.Error())
}

func TestPullRequestsSucceed(t *testing.T) {
	bot := setupTestBot(t)

	var testCases = []struct {
		title string
		file  string
		link  string
	}{
		{"PR Created", "prcreated.json", "/pull_request_created"},
		{"PR Commented", "prcommented.json", "/pull_request_commented"},
		{"PR Approved", "prapproved.json", "/pull_request_approved"},
		{"PR Merged", "prmerged.json", "/pull_request_merged"},
		{"PR Declined", "prdeclined.json", "/pull_request_declined"},
	}

	for _, tt := range testCases {
		reader, err := os.Open(tt.file)
		assert.NoError(t, err)
		req, err := http.NewRequest("POST", tt.link, reader)
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		switch tt.title {
		case "PR Created":
			bot.PullRequestCreated(rr, req)
		case "PR Commented":
			bot.PullRequestCommented(rr, req)
		case "PR Approved":
			bot.PullRequestApproved(rr, req)
		case "PR Merged":
			bot.PullRequestMerged(rr, req)
		case "PR Declined":
			bot.PullRequestDeclined(rr, req)
		}
	}
}
func TestPullRequestsFail(t *testing.T) {
	bot := setupTestBot(t)

	var testCases = []struct {
		title string
		file  string
		link  string
	}{
		{"PR Created", "", "/pull_request_created"},
		{"PR Commented", "", "/pull_request_commented"},
		{"PR Approved", "", "/pull_request_approved"},
		{"PR Merged", "", "/pull_request_merged"},
		{"PR Declined", "", "/pull_request_declined"},
	}

	for _, tt := range testCases {
		reader, err := os.Open(tt.file)
		assert.Error(t, err)
		req, err := http.NewRequest("POST", tt.link, reader)
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		switch tt.title {
		case "PR Created":
			bot.PullRequestCreated(rr, req)
		case "PR Commented":
			bot.PullRequestCommented(rr, req)
		case "PR Approved":
			bot.PullRequestApproved(rr, req)
		case "PR Merged":
			bot.PullRequestMerged(rr, req)
		case "PR Declined":
			bot.PullRequestDeclined(rr, req)
		}
	}
}

func setupTestBot(t *testing.T) *Bot {
	os.Setenv("BOT_TELEGRAM_TOKEN", BotToken)
	os.Setenv("BOT_CHAT", BotChat)
	conf, err := config.GetConfig()
	assert.NoError(t, err)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	r := httpmock.NewStringResponder(200, `{
		"ok": true,
		"result": {
		  "id": 1000,
		  "is_bot": true,
		  "first_name": "testBot",
		  "username": "testbot_bot"
		}
	  }`)

	url := fmt.Sprintf("https://api.telegram.org/bot%v/getMe", conf.TelegramToken)
	httpmock.RegisterResponder("POST", url, r)

	bot, err := NewBot(conf)
	assert.NoError(t, err)
	return bot
}
