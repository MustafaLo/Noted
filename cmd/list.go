/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/MustafaLo/noted/internal"
	"github.com/MustafaLo/noted/models"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func printQueryResponse(dbQueryResponse *notion.DatabaseQueryResponse) {
	if dbQueryResponse == nil || len(dbQueryResponse.Results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Println("\nQuery Results:")
	fmt.Println(strings.Repeat("=", 50))

	for i, page := range dbQueryResponse.Results {
		pageProperties := page.Properties.(notion.DatabasePageProperties)

		note := pageProperties["Note"].RichText[0].Text.Content
		timestamp := pageProperties["Timestamp"].CreatedTime.Local().Format("Jan 2, 2006 3:04 PM")
		category := pageProperties["Category"].Select.Name
		url := page.URL

		// Print out the result neatly with alignment and borders
		fmt.Printf("\n%d. %-10s\n", i+1, strings.Repeat("-", 40))
		fmt.Printf("| %-10s | %s\n", "Note", note)
		fmt.Printf("| %-10s | %s\n", "Timestamp", timestamp)
		fmt.Printf("| %-10s | %s\n", "Category", category)
		fmt.Printf("| %-10s | %s\n", "URL", url)
		fmt.Printf("%-10s\n", strings.Repeat("-", 50))
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "View notes on active file",
	Long: `Use the list command to retrieve all your notes on your current active file
	Example Usage: ./noted list
	Make sure to use quotation marks for multi word categories!`,

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
	listCmdHelpTemplate := internal.CreateHelpTemplate(
		"~~~~~~~~~~~~~~~~~ List Help ~~~~~~~~~~~~~~~~~",
		"./noted list",
		"Use the 'list' command to list out all created notes for your current active file",
		[]string{},
		[]string{
			"./noted list",
		},
		[]string{
			"This command only lists notes for your active file previously created",
		},
	)
	listCmd.SetHelpTemplate(internal.GenerateHelpMessage(listCmdHelpTemplate))
	rootCmd.AddCommand(listCmd)
}
