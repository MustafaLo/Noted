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

func InitService()(*APIService, error){
	err := godotenv.Load()
	if err != nil{
		return nil, fmt.Errorf("Error loading env file: %w", err)
	} 
	client := notion.NewClient(os.Getenv("NOTION_API_KEY"))
	return &APIService{Client: client}, nil
}

func CreateDatabase(s *APIService)(error){
	title := []notion.RichText{
		{
			Text: &notion.Text {
				Content: "Noted CLI Database",
			},
		},
	}

	description := []notion.RichText{
		{
			Text: &notion.Text{
				Content: "Description of Noted CLI Database",
			},
		},
	}

	properties := notion.DatabaseProperties{
		"MyProperty": notion.DatabaseProperty{
			Type: notion.DBPropTypeTitle,
		},
	}

	database_params := notion.CreateDatabaseParams{
		ParentPageID: "172bbbe2-e342-8061-b911-dac8ff678c19",
		Title: title,
		Description: description,
		Properties: properties,
		Icon: nil,
		Cover: nil,
		IsInline: true,
	}

	// Validate the parameters before sending
	if err := database_params.Validate(); err != nil {
		return fmt.Errorf("invalid database parameters: %w", err)
	}

	s.Client.CreateDatabase(context.Background(), database_params)

	return nil
}