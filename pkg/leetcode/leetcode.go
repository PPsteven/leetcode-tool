package leetcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ppsteven/leetcode-tool/internal/config"
	"github.com/sashabaranov/go-openai"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

const RemoteProblems = "https://raw.githubusercontent.com/PPsteven/leetcode-tool/master/problems.json"

type Meta struct {
	Index      string
	Title      string
	Slug       string
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
	Config    *config.Config
	GptClient *openai.Client
	Problems  []byte
}

func NewLeetcode(config *config.Config) *Leetcode {
	client := openai.NewClient("your token")
	return &Leetcode{Config: config, GptClient: client}
}

func DownloadFile(remoteFile string) error {
	// Create the file
	out, err := os.Create("problems.json")
	if err != nil {
		return err
	}
	defer out.Close()

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
	}
	resp, err := client.Get(remoteFile)
	if err != nil {
		_ = os.Remove("problems.json")
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
	file, err := ioutil.ReadFile("problems.json")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println(fmt.Sprintf("file problems.json not exists, start downloading from %s", RemoteProblems))

		err = DownloadFile(RemoteProblems)
		if err != nil {
			log.Fatal(fmt.Errorf("download file failed: %v", err))
		}

		file, err := ioutil.ReadFile("problems.json")
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
	host := "https://leetcode.com"
	if l.Config.Env == "cn" {
		title = "titleCn"
		content = "content.cn"
		difficulty = difficultyMap[strings.ToLower(difficulty)]
		host = "https://leetcode.cn"
	}
	title = problem.Get(title).String()
	content = problem.Get(content).String()

	titleSlug := problem.Get("titleSlug").String()

	return &Meta{
		Index:      number,
		Title:      title,
		Difficulty: difficulty,
		Tags:       tags,
		Link:       fmt.Sprintf("%s/problems/%s/description/", host, titleSlug),
		Slug:       titleSlug,
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
