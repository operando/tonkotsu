package main

import (
	"github.com/BurntSushi/toml"
	"github.com/operando/golack"
)

type Config struct {
	Slack     golack.Slack   `toml:"slack"`
	Webhook   golack.Webhook `toml:"webhook"`
	Payload   golack.Payload `toml:"playload"`
	Log       string         `toml:"log"`
	SleepTime int            `toml:"sleeptime"`
	Android   Android        `toml:"android"`
	Ios       Ios            `toml:"ios"`
}

type Android struct {
	Package string `toml:"package"`
}

type Ios struct {
	Country string `toml:"country"`
	AppId   string `toml:"app_id"`
}

func LoadConfig(configPath string, config *Config) (*Config, error) {
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return config, err
	}
	return config, nil
}
