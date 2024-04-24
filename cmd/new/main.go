package new

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ppsteven/leetcode-tool/pkg/leetcode"
)

type TplFile struct {
	Name     string
	FileName string
	TplStr   string
}

type LanguageConfig struct {
	LeetcodeLang string
	TplFiles     []TplFile
}

const (
	folder = "solve"
	prefix = "solve"
)

var (
	languageConfigs = map[string]LanguageConfig{
		"go": {
			LeetcodeLang: "Go",
			TplFiles:     []TplFile{{"code", "%s.%s.go", codeStrGo}, {"test", "%s.%s_test.go", testCodeStrGo}},
		},
		"ts": {
			LeetcodeLang: "TypeScript",
			TplFiles:     []TplFile{{"code", "%s.%s.ts", codeStrTs}, {"test", "%s.%s.test.ts", testCodeStrTs}},
		},
		"js": {
			LeetcodeLang: "JavaScript",
			TplFiles:     []TplFile{{"code", "%s.%s.js", codeStrJs}, {"test", "%s.%s.test.js", testCodeStrJs}},
		},
		"py3": {
			LeetcodeLang: "Python3",
			TplFiles:     []TplFile{{"code", "%s.%s.py", codeStrPy3}, {"test", "%s.%s_test.py", testCodeStrPy3}, {"__init__", "__init__.py", ""}},
		},
		"java": {
			LeetcodeLang: "Java",
			TplFiles:     []TplFile{{"code", "%s.%s.java", codeStrJava}, {"test", "test_%s_%s.java", testCodeStrJava}},
		},
		"cpp": {
			LeetcodeLang: "C++",
			TplFiles:     []TplFile{{"code", "%s.%s.cpp", codeStrCpp}},
		},
	}
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func normalizeNumber(number string) string {
	if len(number) >= 4 {
		return number
	}
	return strings.Repeat("0", 4-len(number)) + number
}

func mustExecuteTemplate(name string, str string, data interface{}) []byte {
	tpl := template.Must(template.New(name).Parse(str))
	var bf bytes.Buffer
	err := tpl.Execute(&bf, data)
	if err != nil {
		log.Fatalf("mustExecuteTemplate %s error: %s\n", name, err.Error())
	}
	return bf.Bytes()
}

type MetaWithFolder struct {
	leetcode.Meta
	Folder     string
	TagStr     string
	FrontendId string
}

func Run(lc *leetcode.Leetcode, n string, lang string) {
	if lang == "" {
		lang = lc.Config.Lang
	}

	config, ok := languageConfigs[lang]
	if !ok {
		supportLangs := make([]string, len(languageConfigs))
		i := 0
		for l := range languageConfigs {
			supportLangs[i] = l
			i++
		}
		log.Fatalf("invalid lang, now support %s\n", strings.Join(supportLangs, ","))
	}
	meta, err := lc.GetMetaByNumber(n)
	if err != nil || meta == nil {
		log.Fatal(err, meta)
	}
	number := normalizeNumber(meta.Index)
	folderName := number + "." + meta.Slug
	fp := filepath.Join(folder, folderName)
	_ = os.MkdirAll(fp, 0755)
	metaf := &MetaWithFolder{
		*meta,
		folderName,
		strings.Join(meta.Tags, ","),
		n,
	}
	metaf.Meta.Content = strings.ReplaceAll(metaf.Meta.Content, "↵", "")
	//metaf.Meta.Code = gjson.Get(metaf.CodeSnippets, fmt.Sprintf("#(lang=%s).code", config.LeetcodeLang)).String()

	problemFp := filepath.Join(fp, "problem.md")
	if !fileExists(problemFp) {
		bf := mustExecuteTemplate("problem", problemStr, metaf)
		ioutil.WriteFile(problemFp, bf, 0644)
	}

	for _, tpl := range config.TplFiles {
		fileName := tpl.FileName
		if strings.Count(tpl.FileName, "%s") > 1 {
			fileName = fmt.Sprintf(tpl.FileName, number, meta.Slug)
		}
		fp := filepath.Join(fp, fileName)
		if !fileExists(fp) {
			bf := mustExecuteTemplate(tpl.Name, tpl.TplStr, metaf)
			ioutil.WriteFile(fp, bf, 0644)
		}
	}
	fmt.Printf("Done: %s\n", fp)
}

var (
	codeStrGo = `package {{ .Folder }}

/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
 * @solved {{ .Solved }}
*/

{{ .Code }}
`

	testCodeStrGo = `package {{ .Folder }}_test

`

	problemStr = `# [{{ .Index }}. {{ .Title }}]({{ .Link }})

{{ .Content }}
`
)

var (
	codeStrTs = `/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
 * @solved {{ .Solved }}
*/

export {{ .Code }}
`
	testCodeStrTs = `
it('solve_{{ .Index }} should pass', () => {})
`
)

var (
	codeStrJs = `/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
 * @solved {{ .Solved }}
*/

{{ .Code }}
`
	testCodeStrJs = `
it('solve_{{ .Index }} should pass', () => {})
`
)

var (
	codeStrPy3 = `'''
@index {{ .Index }}
@title {{ .Title }}
@difficulty {{ .Difficulty }}
@tags {{ .TagStr }}
@draft false
@link {{ .Link }}
@frontendId {{ .FrontendId }}
@solved {{ .Solved }}
'''

{{ .Code }}
`
	testCodeStrPy3 = `def test_solve():
	pass
`
)

var (
	codeStrJava = `package {{ .Folder }};

/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
 * @solved {{ .Solved }}
*/
{{ .Code }}
`

	testCodeStrJava = `package {{ .Folder }};

public class test_{{ printf "%04s" .Index }} {
	public static void main(String[] args) {
		Solution solution = new Solution();
		// do some test
	}
}
`
)

var (
	codeStrCpp = `#include <algorithm>
#include <unordered_map>
#include <string>
#include <vector>
#include <iostream>
#include <math.h>
using namespace std;

/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
 * @solved {{ .Solved }}
*/
{{ .Code }}
`
)
