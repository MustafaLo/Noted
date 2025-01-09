/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func printQueryResponse(dbQueryResponse *notion.DatabaseQueryResponse) {
	if dbQueryResponse == nil || len(dbQueryResponse.Results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Println("Query Results:")
	fmt.Println("====================")

	for i, page := range dbQueryResponse.Results {
		pageProperties := page.Properties.(notion.DatabasePageProperties)

		note := pageProperties["Note"].RichText[0].Text.Content
		timestamp := pageProperties["Timestamp"].Date.Start
		category := pageProperties["Category"].Select.Name
		url := page.URL

		// Print out the result neatly
		fmt.Printf("%d.\n", i+1)
		fmt.Printf("  Note:      %s\n", note)
		fmt.Printf("  Timestamp: %s\n", timestamp)
		fmt.Printf("  Category:  %s\n", category)
		fmt.Printf("  URL:       %s\n", url)
		fmt.Println("--------------------")
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "View notes on active file",
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
		}
		printQueryResponse(dbQueryResponse)
		
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
