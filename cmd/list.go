/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/spf13/cobra"
)

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
		for _, page := range dbQueryResponse.Results{
			fmt.Println(page.URL)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
