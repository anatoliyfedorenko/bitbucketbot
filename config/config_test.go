package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	conf, err := GetConfig()
	assert.Error(t, err)

	os.Setenv("BOT_TELEGRAM_TOKEN", "testToken")
	os.Setenv("BOT_CHAT", "-12345")

	conf, err = GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, conf.TelegramToken, "testToken")
	assert.Equal(t, conf.Chat, int64(-12345))
}
