/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/MustafaLo/noted/internal"
	"github.com/MustafaLo/noted/models"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"

	client "github.com/cohere-ai/cohere-go/v2/client"
	cohere "github.com/cohere-ai/cohere-go/v2"
)

func getAllNotes(dbQueryResponse *notion.DatabaseQueryResponse)(string){
	var note_block string
	for i, page := range dbQueryResponse.Results{
		pageProperties := page.Properties.(notion.DatabasePageProperties)
		note_block += fmt.Sprintf("Note %d: %s. ", i, pageProperties["Note"].RichText[0].Text.Content)
	}
	return note_block
}

func generateInsights(note_block string)(error){
	fmt.Println("Generating Insights...")
	envMap, err := internal.LoadEnv(); 
	if err != nil{
		return err
	}

	co := client.NewClient(client.WithToken(envMap["COHERE_API_KEY"]))	
	message := fmt.Sprintf("Summarize these notes in 2-3 sentences or fewer, highlighting key insights:\n%s", note_block)
	model_name := "command-r-plus-08-2024"

	resp, err := co.Chat(
		context.TODO(),
		&cohere.ChatRequest{
			Model: &model_name,
			Message: message,
		},
	)
	if err != nil{
		return fmt.Errorf("error making Cohere chat request: %w", err)
	}
	if resp == nil || resp.Text == ""{
		return fmt.Errorf("empty or invalid request from Cohere")
	}

	fmt.Println("\n========== Insights Generated ==========")
	fmt.Println("\n" + resp.Text)
	fmt.Println("\n========================================")

	return nil
}	

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "A brief description of your command",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		activeFileMetaData, err:= internal.GetCurrentFileMetadata()
		if err != nil{
			fmt.Printf("Error %s", err)
			return
		}
		client := cmd.Context().Value("client").(*models.APIService)
		database_id := cmd.Context().Value("databaseID").(string)

		dbQueryResponse, err := internal.FilterDatabase(client, database_id, activeFileMetaData.FileName)
		if err != nil{
			fmt.Printf("Error %s", err)
			return
		} else if (dbQueryResponse == nil || len(dbQueryResponse.Results) == 0){
			fmt.Printf("No insights found")
			return
		}
		notes := getAllNotes(dbQueryResponse)
		if strings.TrimSpace(notes) == ""{
			fmt.Printf("No insights found")
			return
		}

		if err := generateInsights(notes); err != nil{
			fmt.Printf("Error %s", err)
			return
		}
	},
}

func init() {
	insightsCmdHelpTemplate := internal.CreateHelpTemplate(
		"~~~~~~~~~~~~~~~~~ Insights Help ~~~~~~~~~~~~~~~~~",
		"./noted insights",
		"Use the 'insights' command to see AI generated summarizations of your notes for your current active file",
		[]string{},
		[]string{
			"./noted insights",
		},
		[]string{
			"You MUST create your own Cohere API Key in order to use this flag",
			"Insert your api key into your env file as \"COHERE_API_KEY=\"",
		},
	)
	insightsCmd.SetHelpTemplate(internal.GenerateHelpMessage(insightsCmdHelpTemplate))
	rootCmd.AddCommand(insightsCmd)
}
