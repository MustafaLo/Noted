package internal

import (
	"time"

	"github.com/dstotijn/go-notion"
)

func CreateDatabaseEntry(s *APIService, DB_ID string, fileMetaData map[string]interface{}, note string, category string)(error){

}

func createPageParams()(notion.CreatePageParams){

}

func createDatabasePageProperties(fileMetaData map[string]interface{}, note string, category string)(notion.DatabasePageProperties){
	return notion.DatabasePageProperties{
		"File Name": notion.DatabasePageProperty{
			Title: []notion.RichText{{Text:&notion.Text {Content: fileMetaData["fileName"].(string)}}},
		},
		"Note": notion.DatabasePageProperty{
			RichText: []notion.RichText{{Text: &notion.Text{Content: note}}},
		},
		"Line Numbers": notion.DatabasePageProperty{
			RichText: []notion.RichText{{Text: &notion.Text{Content: fileMetaData["lines"].(string)}}},
		},
	
	}
}
