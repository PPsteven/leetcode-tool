package sync

import (
	"context"
	"fmt"
	"github.com/ppsteven/leetcode-tool/internal/meta"
	"github.com/ppsteven/leetcode-tool/internal/notion"
	"github.com/ppsteven/leetcode-tool/pkg/leetcode"
	"golang.org/x/sync/errgroup"
	"log"
)

func Run(lc *leetcode.Leetcode, isNotion bool) {
	tagMetas := meta.GetTagMetas()
	metas, ok := tagMetas["all"]
	if !ok {
		return
	}

	if isNotion {
		nc := notion.NewNotion(lc.Config.Notion.Token).
			WithConfig(lc.Config.Notion.PageID, lc.Config.Notion.DatabaseID)

		err := nc.Init()
		if err != nil {
			log.Fatalf("notion init failed: %v", err)
		}

		g, _ := errgroup.WithContext(context.TODO())
		nThread := 5

		in := make(chan *notion.Record, 10)
		go func() {
			for _, m := range metas {
				record := MetaToRecord(m)
				in <- record
			}
			close(in)
		}()

		progress := make(chan struct{}, 0)
		for i := 0; i < nThread; i++ {
			g.Go(func() error {
				for record := range in {
					err := nc.InsertOrUpdate(record)
					if err != nil {
						return err
					}
					progress <- struct{}{}
				}
				return nil
			})
		}

		go func() {
			cur := 0
			for range progress {
				cur++
				fmt.Printf("\rsync leetcode record to notion, progress: %d/%d", cur, len(metas))
			}
			fmt.Printf("\rsync leetcode record to notion, progress: %d/%d done.\n", cur, len(metas))
		}()

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}
		close(progress)
	}
}

func MetaToRecord(e *meta.Meta) *notion.Record {
	var solved string
	if e.Solved {
		solved = "Yes"
	} else {
		solved = "No"
	}

	fields := []*notion.Field{
		{Type: "text", Name: "_id", Content: e.ID},
		{Type: "text", Name: "ID", Content: e.Index},
		{Type: "title", Name: "Name", Content: e.Title},
		{Type: "url", Name: "Link", Content: e.Link},
		{Type: "select", Name: "Difficulty", Content: e.Difficulty},
		{Type: "multi_select", Name: "Tags", Content: e.Tags},
		{Type: "select", Name: "Solved", Content: solved},
	}
	return &notion.Record{Fields: fields}
}
