package config

import (
	"encoding/json"
	"io/ioutil"
)

type Gpt struct {
	ApiKey string `json:"api_key"` // eg. sk-xxxxxxxxxx
	Model  string `json:"model"`   // eg. gpt-3.5-turbo
}

type Config struct {
	Lang string `json:"lang"`
	Env  string `json:"env"` // eg. en, cn
	Gpt  *Gpt   `json:"gpt"`
}

const configPath = ".leetcode.json"

func NewConfig() *Config {
	c := loadConfig()
	return &c
}

func loadConfig() Config {
	var c Config
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return c
	}
	return c
}
