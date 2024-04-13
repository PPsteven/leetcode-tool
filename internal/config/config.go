package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Lang string `json:"lang"`
	Env  string `json:"env"` // eg. en, cn
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
