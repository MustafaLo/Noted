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

func IntializeDatabase(s *APIService)(*notion.Database, error){
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
			Title: &notion.EmptyMetadata{},
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
		return nil, fmt.Errorf("invalid database parameters: %w", err)
	}

	database, ok := s.Client.CreateDatabase(context.Background(), database_params)
	if ok != nil{
		return nil, fmt.Errorf("unable to validate database parameters: %w", ok)
	}

	return &database, nil
}