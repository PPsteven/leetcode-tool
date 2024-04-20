package gpt

import (
	"github.com/ppsteven/leetcode-tool/internal/gpt"
	"github.com/ppsteven/leetcode-tool/pkg/leetcode"
	"log"
)

func Run(lc *leetcode.Leetcode, number string) {
	if lc.Config.Gpt == nil || lc.Config.Gpt.ApiKey == "" || lc.Config.Gpt.Model == "" {
		log.Fatal("please config gpt api key and model in .leetcode.json")
	}

	client := gpt.NewOpenai(lc.Config.Gpt.ApiKey, lc.Config.Gpt.Model)
	_, err := client.Hint(lc, number)
	if err != nil {
		log.Fatal(err)
	}
}
