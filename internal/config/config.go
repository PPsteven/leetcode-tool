package config

import (
	"github.com/spf13/viper"
	"log"
)

type GptCfg struct {
	ApiKey string `json:"api_key" mapstructure:"api_key"` // eg. sk-xxxxxxxxxx
	Model  string `json:"model" mapstructure:"model"`     // eg. gpt-3.5-turbo
	Prompt string `json:"prompt" mapstructure:"prompt"`   // optional
}

type NotionCfg struct {
	Token      string `json:"token" mapstructure:"token"`
	DatabaseID string `json:"database_id" mapstructure:"database_id"`
}

type Config struct {
	Lang   string     `json:"lang" mapstructure:"lang"`
	Env    string     `json:"env" mapstructure:"env"` // eg. en, cn
	Gpt    *GptCfg    `json:"gpt" mapstructure:"gpt"`
	Notion *NotionCfg `json:"notion" mapstructure:"notion"`
}

const configPath = ".leetcode.json"

func NewConfig() *Config {
	c := loadConfig()
	return &c
}

func loadConfig() Config {
	var c Config
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s, err: %v", configPath, err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %s, err: %v", configPath, err)
	}
	return c
}
