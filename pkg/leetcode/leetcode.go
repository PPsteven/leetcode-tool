package leetcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zcong1993/leetcode-tool/internal/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

const RemoteProblems = "https://raw.githubusercontent.com/PPsteven/leetcode-tool/master/data/problems.json"

type Meta struct {
	Index      string
	Title      string
	Difficulty string
	Tags       []string
	Link       string
	Content    string
	Code       string
	Solved     bool
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

func DownloadFile(remoteFile string) error {
	// Create the file
	out, err := os.Create("data/problems.json")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(remoteFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (l *Leetcode) getAllProblem() ([]byte, error) {
	file, err := ioutil.ReadFile("data/problems.json")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println(fmt.Sprintf("file problems.json not exists, start downloading from %s", RemoteProblems))
		err = DownloadFile(RemoteProblems)

		if err != nil {
			return nil, fmt.Errorf("download file failed: %v", err)
		}
		file, err := ioutil.ReadFile("data/problems.json")
		if err != nil {
			return nil, fmt.Errorf("read file failed: %v", err)
		}
		return file, nil
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
	if l.Problems == nil {
		l.Problems, _ = l.getAllProblem()
	}

	tags := make([]Tag, 0)
	tagsMap := make(map[string]Tag)
	for _, problem := range gjson.ParseBytes(l.Problems).Map() {
		_ = json.Unmarshal([]byte(problem.Get("topicTags").Raw), &tags)
		for _, tag := range tags {
			tagsMap[tag.Slug] = tag
		}
	}
	tags = make([]Tag, 0, len(tagsMap))
	for _, tag := range tagsMap {
		tags = append(tags, tag)
	}
	return tags, nil
}
