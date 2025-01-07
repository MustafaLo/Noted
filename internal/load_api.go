package internal

import (
	"context"
	"fmt"
	"github.com/dstotijn/go-notion"
)

type APIService struct{
	Client *notion.Client
}


func InitService(API_KEY string)(*APIService, error){
	return &APIService{Client: notion.NewClient(API_KEY)}, nil
}

func IntializeDatabase(s *APIService, DB_ID string, PAGE_ID string)(error){
	if DB_ID != ""{
		_, err := s.Client.FindDatabaseByID(context.Background(), DB_ID)
		if err == nil{
			return nil
		}
	}

	database_params := createDatabaseParams(PAGE_ID)
	if err := database_params.Validate(); err != nil {
		return fmt.Errorf("invalid database parameters: %w", err)
	}

	database, err := s.Client.CreateDatabase(context.Background(), database_params)
	if err != nil{
		return fmt.Errorf("unable to validate database parameters: %w", err)
	}

	//Rewrite new environment variables with database id appended
	envMap, err := LoadEnv()
	if err != nil{
		return err
	}

	envMap["NOTION_DATABASE_ID"] = database.ID
	if err := UpdateEnv(envMap); err != nil{
		return err
	}

	return nil
}

func createDatabaseParams(pageID string)(notion.CreateDatabaseParams){
	return notion.CreateDatabaseParams{
		ParentPageID: pageID,
		Title: []notion.RichText{
			{
				Text: &notion.Text {
					Content: "Noted CLI Database",
				},
			},
		},
		Description: []notion.RichText{
			{
				Text: &notion.Text{
					Content: "Database to store notes for your project!",
				},
			},
		},
		Properties: notion.DatabaseProperties{
			"File Name": notion.DatabaseProperty{
				Type: notion.DBPropTypeTitle,
				Title: nil,
			},
			"Note": notion.DatabaseProperty{
				Type: notion.DBPropTypeRichText,
				RichText: nil,
			},
			"Line Numbers": notion.DatabaseProperty{
				Type: notion.DBPropTypeRichText,
				RichText: nil,
			},
			"Timestamp":notion.DatabaseProperty{
				Type: notion.DBPropTypeCreatedTime,
				CreatedTime: nil,
			},
			"Category":notion.DatabaseProperty{
				Type: notion.DBPropTypeSelect,
				Select: &notion.SelectMetadata{
					Options: []notion.SelectOptions{
						{Name: "Bug", Color: "red"},
						{Name: "Feature", Color: "blue"},
						{Name: "Improvement", Color: "green"},
					},
				},
			},
		},
		IsInline: true,
	}
}