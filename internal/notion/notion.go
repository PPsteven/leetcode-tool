package notion

import (
	"context"
	"github.com/jomei/notionapi"
)

type Notion struct {
	DatabaseID notionapi.DatabaseID
	PageID     notionapi.PageID
	client     *notionapi.Client
}

type Field struct {
	Type    string
	Name    string
	Content interface{}
}

type Record struct {
	Fields []*Field
}

func (r *Record) MakeProperties() notionapi.Properties {
	properties := make(notionapi.Properties)
	for _, field := range r.Fields {
		switch field.Type {
		case "title":
			properties[field.Name] = &notionapi.TitleProperty{
				Type: "title",
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: field.Content.(string)}},
				},
			}
		case "text":
			properties[field.Name] = &notionapi.RichTextProperty{
				Type: "rich_text",
				RichText: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: field.Content.(string),
						},
					},
				},
			}
		case "url":
			properties[field.Name] = &notionapi.URLProperty{
				Type: "url",
				URL:  field.Content.(string),
			}
		case "select":
			properties[field.Name] = &notionapi.SelectProperty{
				Type:   "select",
				Select: notionapi.Option{Name: field.Content.(string)},
			}
		case "multi_select":
			values := field.Content.([]string)
			options := make([]notionapi.Option, len(values))
			for i, val := range values {
				options[i] = notionapi.Option{Name: val}
			}

			properties[field.Name] = &notionapi.MultiSelectProperty{
				Type:        "multi_select",
				MultiSelect: options,
			}
		}
	}
	return properties
}

func NewNotion(token string) *Notion {
	client := notionapi.NewClient(notionapi.Token(token))
	return &Notion{
		client: client,
	}
}

func (n *Notion) WithConfig(pageID, databaseID string) *Notion {
	n.DatabaseID = notionapi.DatabaseID(databaseID)
	n.PageID = notionapi.PageID(pageID)
	return n
}

func (n *Notion) Insert(record *Record) error {
	ctx := context.Background()

	pageReq := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: n.DatabaseID,
		},
		Properties: record.MakeProperties(),
	}

	_, err := n.client.Page.Create(ctx, pageReq)
	return err
}
