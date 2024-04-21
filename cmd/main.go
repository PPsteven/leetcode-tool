package main

import (
	"fmt"
	"github.com/ppsteven/leetcode-tool/cmd/gpt"
	"github.com/ppsteven/leetcode-tool/internal/config"
	"log"
	"os"

	"github.com/ppsteven/leetcode-tool/cmd/new"
	"github.com/ppsteven/leetcode-tool/cmd/tags"
	"github.com/ppsteven/leetcode-tool/cmd/update"
	"github.com/ppsteven/leetcode-tool/pkg/leetcode"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version = "master"
	Commit  = ""
	Date    = ""
	BuiltBy = ""
)

var (
	app = kingpin.New("leetcode-tool", "一个让你更方便刷题的工具.")

	updateCmd = app.Command("update", "Update readme.")

	newCmd = app.Command("new", "Init a new leetcode problem.")
	number = newCmd.Arg("number", "problem number").Required().String()
	lang   = newCmd.Flag("lang", "language").String()

	metaCmd    = app.Command("meta", "Show problem meta by number.")
	metaNumber = metaCmd.Arg("number", "problem number").Required().String()

	tagsCmd   = app.Command("tags", "Update tag toc files.")
	tagsForce = tagsCmd.Flag("force", "force update file").Short('f').Bool()

	gptCmd    = app.Command("gpt", "Use gpt to solve problem.")
	gptNumber = gptCmd.Arg("number", "problem number").Required().String()
)

func showMeta(lc *leetcode.Leetcode, number string) {
	meta, err := lc.GetMetaByNumber(number)
	if err != nil {
		log.Fatal(err)
	}
	if meta == nil {
		log.Fatal("mate not found")
	}
	meta.Content = ""
	fmt.Printf("%+v\n", meta)
}

func main() {
	app.Version(buildVersion(Version, Commit, Date, BuiltBy))
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	lc := leetcode.NewLeetcode(config.NewConfig())

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case updateCmd.FullCommand():
		update.Run()
	case newCmd.FullCommand():
		new.Run(lc, *number, *lang)
	case metaCmd.FullCommand():
		showMeta(lc, *metaNumber)
	case tagsCmd.FullCommand():
		tags.Run(lc, *tagsForce)
	case gptCmd.FullCommand():
		gpt.Run(lc, *gptNumber)
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
