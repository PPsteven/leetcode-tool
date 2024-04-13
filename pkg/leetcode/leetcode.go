package leetcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zcong1993/leetcode-tool/internal/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

type Meta struct {
	Index      string
	Title      string
	Difficulty string
	Tags       []string
	Link       string
	Content    string
	//Code         string
	//CodeSnippets string
}

type Tag struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	TranslatedName string `json:"translatedName"`
}

var (
	difficultyMap = map[string]string{
		"easy":   "简单",
		"medium": "中等",
		"hard":   "困难",
	}
)

type Leetcode struct {
	Config   *config.Config
	Problems []byte
}

func NewLeetcode(config *config.Config) *Leetcode {
	return &Leetcode{Config: config}
}

func (l *Leetcode) getAllProblem() ([]byte, error) {
	file, err := ioutil.ReadFile("/Users/ppsteven/Projects/leetcode-tool/data/problems.json")
	if err == os.ErrNotExist {
		return nil, errors.New("234324")
	}
	return file, nil
}

func (l *Leetcode) getDetail(number string) (*Meta, error) {
	if number == "" {
		return nil, nil
	}

	problem := gjson.GetBytes(l.Problems, fmt.Sprintf("%s", number))

	tagsResult := problem.Get("topicTags.#.slug").Array()
	tags := make([]string, len(tagsResult))
	for i, t := range tagsResult {
		tags[i] = t.String()
	}

	title := "title"
	difficulty := problem.Get("difficulty").String()
	content := "content.en"
	if l.Config.Env == "cn" {
		title = "titleCn"
		content = "content.cn"
		difficulty = difficultyMap[strings.ToLower(difficulty)]
	}
	title = problem.Get(title).String()
	content = problem.Get(content).String()

	return &Meta{
		Index:      number,
		Title:      title,
		Difficulty: difficulty,
		Tags:       tags,
		Link:       fmt.Sprintf("https://leetcode.cn/problems/%s/description/", problem.Get("titleSlug").String()),
		Content:    content,
	}, nil
}

func (l *Leetcode) GetMetaByNumber(number string) (*Meta, error) {
	if l.Problems == nil {
		l.Problems, _ = l.getAllProblem()
	}
	return l.getDetail(number)
}

func (l *Leetcode) GetTags() ([]Tag, error) {
	resp, err := http.Get("https://leetcode-cn.com/problems/api/tags/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := make([]Tag, 0)
	err = json.Unmarshal([]byte(gjson.GetBytes(bt, "topics").Raw), &res)
	return res, err
}
