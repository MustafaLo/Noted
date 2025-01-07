package internal

import (
	"github.com/dstotijn/go-notion"
)

func CreateDatabaseEntry(s *APIService, DB_ID string, fileMetaData map[string]interface{}, note string, category string)(error){

}

func createPageParams()(notion.CreatePageParams){

}

func createDatabasePageProperties(fileMetaData map[string]interface{}, note string, category string)(notion.DatabaseProperties){
	return notion.DatabaseProperties{
		"File Name": notion.DatabaseProperty{
			Type: notion.DBPropTypeTitle,
			Title: []notion.RichText{{Text: &notion.Text{Content: fileMetaData["fileName"].(string)}}},
		},
		"Note": notion.DatabaseProperty{
			Type: notion.DBPropTypeRichText,
			RichText: []notion.RichText{{Text: &notion.Text{Content: note}}},
		},
	}
}
