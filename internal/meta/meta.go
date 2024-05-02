package meta

import (
	"github.com/bmatcuk/doublestar/v2"
	"github.com/ppsteven/leetcode-tool/internal/helper"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	indexRegex      = regexp.MustCompile("@index (.+)")
	titleRegex      = regexp.MustCompile("@title (.+)")
	difficultyRegex = regexp.MustCompile("@difficulty (.+)")
	tagsRegex       = regexp.MustCompile("@tags (.+)")
	draftRegex      = regexp.MustCompile("@draft (.+)")
	linkRegex       = regexp.MustCompile("@link (.+)")
	frontendIdRegex = regexp.MustCompile("@frontendId (.+)")
	solvedRegex     = regexp.MustCompile("@solved (.+)")
)

type Meta struct {
	Index      string
	Title      string
	Difficulty string
	Tags       []string
	Draft      bool
	Fp         string
	Link       string
	FrontendId string
	Ext        string
	Solved     bool
	Completed  string
}

type TagMetas map[string](Metas)

type Metas []*Meta

func (a Metas) Len() int      { return len(a) }
func (a Metas) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Metas) Less(i, j int) bool {
	iIndex, _ := strconv.Atoi(a[i].Index)
	jIndex, _ := strconv.Atoi(a[j].Index)
	return iIndex < jIndex
}

func GetTagMetas() TagMetas {
	files, err := doublestar.Glob("./solve/**/*")
	if err != nil {
		log.Fatal(err)
	}
	tagMetas := make(TagMetas, 0)
	tagMetas["all"] = make(Metas, 0)
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for _, fp := range files {
		if isFolder, _ := helper.IsDirectory(fp); isFolder {
			continue
		}

		if strings.HasSuffix(fp, ".md") {
			continue
		}
		wg.Add(1)
		fp := fp
		go func() {
			content, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			meta := findMeta(content, fp)
			if meta != nil {
				lock.Lock()
				addMeta(tagMetas, meta)
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return tagMetas
}

func addMeta(tagMetas TagMetas, meta *Meta) {
	if meta == nil {
		return
	}
	for _, tag := range meta.Tags {
		if _, ok := tagMetas[tag]; !ok {
			tagMetas[tag] = make(Metas, 0)
		}
		tagMetas[tag] = append(tagMetas[tag], meta)
	}
	tagMetas["all"] = append(tagMetas["all"], meta)
}

func findTag(content []byte, r *regexp.Regexp) string {
	res := r.FindSubmatch(content)
	if len(res) < 2 {
		return ""
	}
	return string(res[1])
}

func findMeta(content []byte, fp string) *Meta {
	draft := findTag(content, draftRegex) == "" || findTag(content, draftRegex) == "true"
	if draft {
		return nil
	}
	tags := strings.Split(findTag(content, tagsRegex), ",")
	solved := false
	if strings.ToLower(findTag(content, solvedRegex)) == "true" {
		solved = true
	}

	return &Meta{
		Index:      findTag(content, indexRegex),
		Title:      findTag(content, titleRegex),
		Difficulty: findTag(content, difficultyRegex),
		Tags:       tags,
		Draft:      draft,
		Fp:         filepath.Dir(fp),
		Link:       findTag(content, linkRegex),
		FrontendId: findTag(content, frontendIdRegex),
		Ext:        filepath.Ext(fp),
		Solved:     solved,
		Completed:  genCompleted(solved, filepath.Ext(fp)),
	}
}

func genCompleted(isCompleted bool, ext string) string {
	if isCompleted {
		return ext[1:] + " ✅"
	}
	return ext[1:] + " ❌"
}
