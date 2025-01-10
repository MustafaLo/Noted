/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func getAllNotes(dbQueryResponse *notion.DatabaseQueryResponse)(string){
	var note_block string
	for i, page := range dbQueryResponse.Results{
		pageProperties := page.Properties.(notion.DatabasePageProperties)
		note_block += fmt.Sprintf("Note %d: %s. ", i, pageProperties["Note"].RichText[0].Text.Content)
	}
	return note_block
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
		fmt.Println(notes)
	},
}

func init() {
	rootCmd.AddCommand(insightsCmd)
}
