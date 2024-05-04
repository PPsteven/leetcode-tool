package notion

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"log"
	"strings"
)

type PageUID string

type SigAndID struct {
	Signature string
	PageID    notionapi.PageID
}

type Notion struct {
	DatabaseID notionapi.DatabaseID
	PageID     notionapi.PageID
	client     *notionapi.Client
	PageSig    map[PageUID]*SigAndID
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

type PageData struct {
	PageID notionapi.PageID
	Data   map[string]string
}

func ParsePage(page *notionapi.Page) (PageUID, *PageData) {
	var pageUID PageUID
	pageID := GetStandardID(page.ID)
	data := make(map[string]string)
	for name, property := range page.Properties {
		data[name] = ParseProperty(property)
	}
	if v, ok := data["_id"]; ok {
		pageUID = PageUID(v)
	} else {
		log.Fatalf("_id not found: %v", data)
	}
	return pageUID, &PageData{PageID: notionapi.PageID(pageID), Data: data}
}

func ParseProperty(property notionapi.Property) string {
	switch property.GetType() {
	case "title":
		p, ok := property.(*notionapi.TitleProperty)
		if !ok {
			log.Fatalf("title parsed failed: %v", property)
		}
		if len(p.Title) == 0 {
			return ""
		}
		return p.Title[0].Text.Content
	case "rich_text":
		p, ok := property.(*notionapi.RichTextProperty)
		if !ok {
			log.Fatalf("rich_text parsed failed: %v", property)
		}
		if len(p.RichText) == 0 {
			return ""
		}
		return p.RichText[0].Text.Content
	case "url":
		p, ok := property.(*notionapi.URLProperty)
		if !ok {
			log.Fatalf("url parsed failed: %v", property)
		}
		return p.URL
	case "select":
		p, ok := property.(*notionapi.SelectProperty)
		if !ok {
			log.Fatalf("select parsed failed: %v", property)
		}
		return p.Select.Name
	case "multi_select":
		p, ok := property.(*notionapi.MultiSelectProperty)
		if !ok {
			log.Fatalf("multi_select parsed failed: %v", property)
		}
		optionStr := ""
		for _, option := range p.MultiSelect {
			optionStr += option.Name + ","
		}
		if len(optionStr) > 0 {
			optionStr = optionStr[:len(optionStr)-1]
		}
		return optionStr
	default:
		log.Fatalf("%s not supported current", property.GetType())
	}
	return ""
}

func NewNotion(token string) *Notion {
	client := notionapi.NewClient(notionapi.Token(token), notionapi.WithRetry(10))
	return &Notion{
		client: client,
	}
}

func (n *Notion) WithConfig(pageID, databaseID string) *Notion {
	n.PageID = notionapi.PageID(pageID)
	n.DatabaseID = notionapi.DatabaseID(databaseID)
	return n
}

func (n *Notion) Init() error {
	if n.DatabaseID == "" && n.PageID == "" {
		return fmt.Errorf("both database_id and page_id are empty")
	}
	// 1. 未填写database_id则创建数据库
	if n.DatabaseID == "" && n.PageID != "" {
		db, err := n.CreateDB()
		if err != nil {
			return fmt.Errorf("create db failed: %v", err)
		}
		n.DatabaseID = notionapi.DatabaseID(GetStandardID(db.ID))

		log.Printf("Create Database: %s, Please add the database_id in the config file\n", n.DatabaseID)
		log.Printf("Visited Link: %s", fmt.Sprintf("https://www.notion.so/%s", n.DatabaseID))
	}

	// 2. 创建完成database_id则添加记录
	if n.DatabaseID != "" {
		records, err := n.Query()
		if err != nil {
			return fmt.Errorf("query failed: %v", err)
		}
		n.PageSig = make(map[PageUID]*SigAndID)
		for pageUID, pageData := range records {
			n.PageSig[pageUID] = &SigAndID{PageID: pageData.PageID, Signature: GetPageSig(pageData.Data)}
		}
	}
	return nil
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

func (n *Notion) Update(pageID notionapi.PageID, record *Record) error {
	ctx := context.Background()

	updateReq := &notionapi.PageUpdateRequest{
		Properties: record.MakeProperties(),
	}
	_, err := n.client.Page.Update(ctx, pageID, updateReq)

	return err
}

func (n *Notion) Query() (records map[PageUID]*PageData, err error) {
	ctx := context.Background()

	dbQueryReq := &notionapi.DatabaseQueryRequest{}

	queryResp, err := n.client.Database.Query(ctx, n.DatabaseID, dbQueryReq)

	if err != nil {
		return nil, err
	}

	records = make(map[PageUID]*PageData)
	for _, result := range queryResp.Results {
		pageUID, pageData := ParsePage(&result)
		records[pageUID] = pageData
	}
	return records, nil
}

func (n *Notion) InsertOrUpdate(record *Record) error {
	pageUID, pageData := ParsePage(&notionapi.Page{Properties: record.MakeProperties()})

	if pageSig, ok := n.PageSig[pageUID]; ok {
		// 签名不一致，更新行信息
		if pageSig.Signature != GetPageSig(pageData.Data) {
			return n.Update(n.PageSig[pageUID].PageID, record)
		} else {
			return nil
		}
	}

	return n.Insert(record)
}

var EmptySelect = notionapi.Select{Options: make([]notionapi.Option, 0)}

func (n *Notion) CreateDB() (db *notionapi.Database, err error) {
	ctx := context.Background()

	dbReq := &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   "page_id",
			PageID: n.PageID,
		},
		Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "Leetcode"}}},
		Properties: notionapi.PropertyConfigs{
			"Link":       &notionapi.URLPropertyConfig{Type: "url"},
			"Tags":       &notionapi.MultiSelectPropertyConfig{Type: "multi_select", MultiSelect: EmptySelect},
			"Solved":     &notionapi.SelectPropertyConfig{Type: "select", Select: EmptySelect},
			"Difficulty": &notionapi.SelectPropertyConfig{Type: "select", Select: EmptySelect},
			"ID":         &notionapi.RichTextPropertyConfig{Type: "rich_text"},
			"_id":        &notionapi.RichTextPropertyConfig{Type: "rich_text"},
			"Name":       &notionapi.TitlePropertyConfig{Type: "title"},
		},
		IsInline: true, // show database inline in the parent page
	}

	return n.client.Database.Create(ctx, dbReq)
}

func GetPageSig(v map[string]string) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", v["Difficulty"], v["ID"], v["Link"], v["Name"], v["Solved"], v["Tags"])
}

func GetStandardID(stringer fmt.Stringer) string {
	return strings.ReplaceAll(stringer.String(), "-", "")
}
