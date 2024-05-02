package update

import (
	"bytes"
	"fmt"
	"github.com/ppsteven/leetcode-tool/internal/helper"
	"github.com/ppsteven/leetcode-tool/internal/meta"
	"github.com/ppsteven/leetcode-tool/pkg/leetcode"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"
)

const (
	toc = "toc"
)

var (
	tableTpl = template.Must(template.New("table").Parse(tableStr))
	tagTpl   = template.Must(template.New("tag").Parse(tagStr))
)

type (
	Meta  = meta.Meta
	Metas = meta.Metas
)

type TableData struct {
	Metas Metas
	Total int
}

func genTable(data *TableData) string {
	var bf bytes.Buffer
	sort.Sort(data.Metas)
	tableTpl.Execute(&bf, data)
	return bf.String()
}

func Run() {
	tagMetas := meta.GetTagMetas()

	if !helper.FileExists(toc) {
		_ = os.MkdirAll(toc, 0755)
	}

	wg := sync.WaitGroup{}

	for tag, metas := range tagMetas {
		fp := filepath.Join(toc, fmt.Sprintf("%s.md", tag))
		wg.Add(1)
		metas := metas
		tag := tag
		go func() {
			if !helper.FileExists(fp) {
				var content bytes.Buffer
				err := tagTpl.Execute(&content, &leetcode.Tag{Name: tag})
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(fp, content.Bytes(), 0644)
				if err != nil {
					log.Printf("write file %s error, %s\n", fp, err)
				}
			}
			content, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			table := genTable(&TableData{
				Metas: metas,
				Total: len(metas),
			})
			contents := strings.Split(string(content), "<!--- table -->")
			contents[1] = "<!--- table -->\n" + table
			newContent := strings.Join(contents, "")
			ioutil.WriteFile(fp, []byte(newContent), 0644)
			wg.Done()
		}()
	}
	wg.Wait()
}

var tableStr = `

总计: {{ .Total }}

| 网页序号 | 序号 | 难度 | 题目                    | 解答                      | 完成 |
| ---- | ---- | ---- | ------------------ | ---------------- | -------- | {{ range .Metas }}
| {{ .FrontendId }} | {{ .Index }} | {{ .Difficulty }} | [{{ .Title }}]({{ .Link }}) | [{{ .Fp }}](../{{ .Fp }})| {{ .Completed }} |{{ end }}
`

var tagStr = `# {{ .Name }}

<!--- table -->

`
