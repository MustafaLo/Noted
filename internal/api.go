package internal

import (
	"context"
	"fmt"
	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
)

type APIService struct{
	Client *notion.Client
}

var s APIService
var envMap map[string]string

func LoadEnv()(error){
	if envMap != nil{
		return nil
	}

	data, err := godotenv.Read()
	if err != nil{
		return fmt.Errorf("error loading env file: %w", err)
	}

	if _, exists := data["NOTION_API_KEY"]; !exists{
		return fmt.Errorf("env file needs authentication key")
	}
	envMap = data
	return nil
}

func InitService()(error){
	//Client is already authenticated
	if s.Client != nil{
		return nil
	}

	client := notion.NewClient(envMap["NOTION_API_KEY"])
	s.Client = client
	return nil
}

func IntializeDatabase()(error){
	if databaseID, exists := envMap["NOTION_DATABASE_ID"]; exists{
		_, err := s.Client.FindDatabaseByID(context.Background(), databaseID)
		if err == nil{
			return nil
		}
	}

	database_params := createDatabaseParams(envMap["NOTION_PAGE_ID"])
	if err := database_params.Validate(); err != nil {
		return fmt.Errorf("invalid database parameters: %w", err)
	}

	database, err := s.Client.CreateDatabase(context.Background(), database_params)
	if err != nil{
		return fmt.Errorf("unable to validate database parameters: %w", err)
	}

	envMap["NOTION_DATABASE_ID"] = database.ID
	if err := godotenv.Write(envMap, ".env"); err != nil {
		return fmt.Errorf("error writing to .env file: %w", err)
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
				Title: &notion.EmptyMetadata{},
			},
			"Note": notion.DatabaseProperty{
				Type: notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			},
			"Line Numbers": notion.DatabaseProperty{
				Type: notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			},
			"Timestamp":notion.DatabaseProperty{
				Type: notion.DBPropTypeCreatedTime,
				CreatedTime: &notion.EmptyMetadata{},
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