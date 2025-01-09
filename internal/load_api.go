package internal

import (
	"context"
	"fmt"

	"github.com/MustafaLo/Noted/models"
	"github.com/dstotijn/go-notion"
)


func InitService(API_KEY string)(*models.APIService, error){
	return &models.APIService{Client: notion.NewClient(API_KEY)}, nil
}

func IntializeDatabase(s *models.APIService, DB_ID string, PAGE_ID string)(string, error){
	if DB_ID != ""{
		_, err := s.Client.FindDatabaseByID(context.Background(), DB_ID)
		if err == nil{
			return "", nil
		}
		return DB_ID, nil
	}

	database_params := createDatabaseParams(PAGE_ID)
	if err := database_params.Validate(); err != nil {
		return "", fmt.Errorf("invalid database parameters: %w", err)
	}

	database, err := s.Client.CreateDatabase(context.Background(), database_params)
	if err != nil{
		return "", fmt.Errorf("unable to validate database parameters: %w", err)
	}

	//Rewrite new environment variables with database id appended
	envMap, err := LoadEnv()
	if err != nil{
		return "", err
	}

	envMap["NOTION_DATABASE_ID"] = database.ID
	if err := UpdateEnv(envMap); err != nil{
		return "", err
	}

	return envMap["NOTION_DATABASE_ID"], nil
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