package internal

import (
	"context"
	"fmt"
	"os"
	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
)

type APIService struct{
	Client *notion.Client
}

var s APIService

func InitService()(error){
	//Client is already authenticated
	if s.Client != nil{
		return nil
	}

	err := godotenv.Load()
	if err != nil{
		return fmt.Errorf("error loading env file: %w", err)
	} 
	client := notion.NewClient(os.Getenv("NOTION_API_KEY"))
	s.Client = client
	return nil
}

func IntializeDatabase()(error){
	databaseID := os.Getenv("NOTION_DATABASE_ID")
	if databaseID != ""{
		_, err := s.Client.FindDatabaseByID(context.Background(), databaseID)
		if err == nil{
			return nil
		}
	}

	database_params := createDatabaseParams()

	if err := database_params.Validate(); err != nil {
		return fmt.Errorf("invalid database parameters: %w", err)
	}

	database, err := s.Client.CreateDatabase(context.Background(), database_params)
	if err != nil{
		return fmt.Errorf("unable to validate database parameters: %w", err)
	}

	os.Setenv("NOTION_DATABASE_ID", database.ID)
	return nil
}

func createDatabaseParams()(notion.CreateDatabaseParams){
	return notion.CreateDatabaseParams{
		ParentPageID: os.Getenv("NOTION_PAGE_ID"),
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
						{Name: "Bug"},
						{Name: "Feature"},
						{Name: "Improvement"},
					},
				},
			},
		},
		IsInline: true,
	}
}