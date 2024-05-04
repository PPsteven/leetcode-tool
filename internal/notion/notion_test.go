package notion

import (
	"github.com/ppsteven/leetcode-tool/internal/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotion(t *testing.T) {
	notion := NewNotion("secret_xxxxxx")
	notion.WithConfig("xxx", "xxxx")

	m := &meta.Meta{
		"465b111ca3e74113a0fc382aa1d3dfa6",
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

	t.Run("Insert", func(t *testing.T) {
		err := notion.Insert(record)
		assert.NoError(t, err)
	})

	t.Run("Query", func(t *testing.T) {
		ret, err := notion.Query()
		assert.NoError(t, err)
		t.Log(ret)
	})

	t.Run("Init", func(t *testing.T) {
		err := notion.Init()
		assert.NoError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		err := notion.Update("465b111ca3e74113a0fc382aa1d3dfa6", MetaToRecord(m))
		assert.NoError(t, err)
	})

	t.Run("InsertOrUpdate", func(t *testing.T) {
		_ = notion.Init()

		m.Title = "test"
		err := notion.InsertOrUpdate(MetaToRecord(m))
		assert.NoError(t, err)
	})

	t.Run("CreateDB", func(t *testing.T) {
		db, err := notion.CreateDB()
		assert.NoError(t, err)
		t.Log(db)
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
		{Type: "text", Name: "_id", Content: e.ID},
		{Type: "text", Name: "ID", Content: e.Index},
		{Type: "title", Name: "Name", Content: e.Title},
		{Type: "url", Name: "Link", Content: e.Link},
		{Type: "select", Name: "Difficulty", Content: e.Difficulty},
		{Type: "multi_select", Name: "Tags", Content: e.Tags},
		{Type: "select", Name: "Solved", Content: solved},
	}
	return &Record{Fields: fields}
}
