package gpt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ppsteven/leetcode-tool/pkg/leetcode"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"text/template"
)

type Openai struct {
	model  string
	client *openai.Client
}

func NewOpenai(apiKey, model string) *Openai {
	client := openai.NewClient(apiKey)
	return &Openai{
		client: client,
		model:  model,
	}
}

func (o *Openai) Chat(content string) (string, error) {
	stream, err := o.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: o.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	defer stream.Close()

	ans := ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return "", fmt.Errorf("Stream error: %v", err)
		}

		words := response.Choices[0].Delta.Content
		ans += words

		fmt.Printf(words)
	}

	return ans, nil
}

func (o *Openai) Hint(lc *leetcode.Leetcode, number string) (string, error) {
	meta, err := lc.GetMetaByNumber(number)
	if err != nil {
		return "", err
	}

	textLang := "中文"
	if lc.Config.Lang == "en" {
		textLang = "English"
	}

	var content bytes.Buffer
	err = hitTpl.Execute(&content, &HintData{
		Lang:     lc.Config.Lang,
		TextLang: textLang,
		Problem:  meta.Content,
	})
	if err != nil {
		log.Fatal(err)
	}
	return o.Chat(content.String())
}

type HintData struct {
	Lang     string
	TextLang string
	Problem  string
}

var hitTpl = template.Must(template.New("hint").Parse(hitStr))

var hitStr = `
您是一个算法专家，请基于下面的算法题目，给出该算法的思路和复杂度, 使用 {{ .TextLang }} 回答
SETP1. 给出算法的归类，如递归，栈
SETP2. 若是存在暴力解法，给出思路和复杂度
SETP3. 给出最优解法和复杂度
SETP4. 代码实现，使用 {{ .Lang }} 语言，代码带注释和测试样例。

{{ .Problem }}
`
