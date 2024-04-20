package config

import (
	"github.com/spf13/viper"
	"log"
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
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s, err: %v", configPath, err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %s, err: %v", configPath, err)
	}
	return c
}
