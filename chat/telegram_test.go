package chat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/bitbucketbot/config"
)

func TestNewBotReturnsError(t *testing.T) {
	os.Setenv("BOT_TELEGRAM_TOKEN", "testToken")
	os.Setenv("BOT_CHAT", "-12345")
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

func TestSendUpdate(t *testing.T) {
	bot := setupTestBot(t)
	bot.SendUpdate("test")
}

func TestPullRequestCreated(t *testing.T) {
	bot := setupTestBot(t)

	reader, err := os.Open("prcreated.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_created", reader)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestCreated(rr, req)

}

func TestPullRequestMerged(t *testing.T) {
	bot := setupTestBot(t)

	reader, err := os.Open("prmerged.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_merged", reader)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestMerged(rr, req)

}

func TestPullRequestCommented(t *testing.T) {
	bot := setupTestBot(t)

	reader, err := os.Open("prcommented.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_commented", reader)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestCommented(rr, req)

}

func TestPullRequestApproved(t *testing.T) {
	bot := setupTestBot(t)

	reader, err := os.Open("prapproved.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_approved", reader)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestApproved(rr, req)

}

func setupTestBot(t *testing.T) *Bot {
	os.Setenv("BOT_TELEGRAM_TOKEN", "testToken")
	os.Setenv("BOT_CHAT", "-12345")
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
