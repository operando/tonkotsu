package main

import (
	"github.com/BurntSushi/toml"
	"github.com/operando/golack"
)

type Config struct {
	SlackUpdatePost golack.Slack   `toml:"slack_update_post"`
	SlackErrorPost  golack.Slack   `toml:"slack_error_post"`
	SlackStartPost  golack.Slack   `toml:"slack_start_post"`
	Webhook         golack.Webhook `toml:"webhook"`
	Log             string         `toml:"log"`
	SleepTime       int            `toml:"sleeptime"`
	Android         Android        `toml:"android"`
	Ios             Ios            `toml:"ios"`
	ErrorPost       bool           `toml:"error_post"`
}

type Android struct {
	Package string `toml:"package"`
}

type Ios struct {
	Country string `toml:"country"`
	AppID   string `toml:"app_id"`
}

func LoadConfig(configPath string, config *Config) (*Config, error) {
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return config, err
	}
	return config, nil
}
