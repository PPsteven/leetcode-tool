package notion

import (
	"github.com/ppsteven/leetcode-tool/internal/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotion(t *testing.T) {
	notion := NewNotion("secret_xxxxxx")
	notion.WithConfig("", "xxxx")

	m := &meta.Meta{
		"1",
		"Binary Search",
		"Medium",
		[]string{"Arrays", "Binary Search"},
		false,
		"1",
		"https://leetcode.com/problems/binary-search/",
		"1",
		"md",
		true,
		""}
	record := MetaToRecord(m)

	t.Run("", func(t *testing.T) {
		err := notion.Insert(record)
		assert.NoError(t, err)
	})
}

func MetaToRecord(e *meta.Meta) *Record {
	var solved string
	if e.Solved {
		solved = "Yes"
	} else {
		solved = "No"
	}

	fields := []*Field{
		{Type: "title", Name: "ID", Content: e.Index},
		{Type: "text", Name: "Name", Content: e.Title},
		{Type: "url", Name: "Link", Content: e.Link},
		{Type: "select", Name: "Difficulty", Content: e.Difficulty},
		{Type: "multi_select", Name: "Tags", Content: e.Tags},
		{Type: "select", Name: "Solved", Content: solved},
	}
	return &Record{Fields: fields}
}
