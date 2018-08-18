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
	bot.SendUpdate("test")
}

func TestPullRequestCreated(t *testing.T) {
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

	reader, err := os.Open("prcreated.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_created", reader)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestCreated(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPullRequestApproved(t *testing.T) {
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

	reader, err := os.Open("prmerged.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/pull_request_merged", reader)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	bot.PullRequestMerged(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
