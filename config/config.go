package config

import "github.com/kelseyhightower/envconfig"

//Config defines config struct
type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	Chat          int64  `envconfig:"CHAT" required:"true"`
}

// GetConfig method processes env variables and fills Config struct
func GetConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("bot", &c); err != nil {
		return c, err
	}
	return c, nil
}
