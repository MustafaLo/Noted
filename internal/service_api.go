package internal

import (
	"context"
	"fmt"
	"github.com/dstotijn/go-notion"
)

func CreateDatabaseEntry(s *APIService, DB_ID string, fileMetaData map[string]interface{}, note string, category string)(error){
	database_page_properties := createDatabasePageProperties(fileMetaData, note, category)
	page_parameters := createPageParams(DB_ID, database_page_properties)
	page, err := s.Client.CreatePage(context.Background(), page_parameters)
	if err != nil{
		return fmt.Errorf("Error creating page: ", err)
	}
	fmt.Println("Page created successfully: ", page.URL)
	return nil
}

func createPageParams(DB_ID string, db_page_props notion.DatabasePageProperties)(notion.CreatePageParams){
	return notion.CreatePageParams{
		ParentType: "database_id",
		ParentID: DB_ID,
		DatabasePageProperties: &db_page_props,
	}
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
			RichText: []notion.RichText{{Text: &notion.Text{Content: "Test"}}},
		},
		"Category": notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: category},
		},
	}
}
