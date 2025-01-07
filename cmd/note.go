/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MustafaLo/Noted/internal"
	"github.com/spf13/cobra"
)


func getCurrentFileMetadata()(map[string]interface{}, error){
	var metadata map[string]interface{}
	data, err := os.ReadFile("fileMetadata.json")
	if err != nil{
		return nil, fmt.Errorf("failed to open fileMetaData -- make sure to enable File Tracker Extension")
	}

	ok := json.Unmarshal(data, &metadata)
	if ok != nil{
		return nil, fmt.Errorf("failed to parse fileMetaData")
	}
	return metadata, nil
}


func printFileMetaData(data map[string]interface{}){
	for k, v := range data{
		if k == "lines" && v != nil{
			if linesMap, ok := v.(map[string]interface{}); ok {
				fmt.Printf("key: %s, value: {", k)
				for lk, lv := range linesMap {
					fmt.Printf("%s: %0.f,", lk, lv) // Adjust formatting as needed
				}
				fmt.Println("}")
			} else {
				fmt.Printf("Not okay")
			}
		} else{
			fmt.Printf("key: %s, value: %s\n", k, v)
		}
	}
}

func isValidLinesFormat()(error){
	
}

func setLines(highlighted map[string]interface{})(string, error){

}



var note string
var client *internal.APIService
var databaseID string
var category string
var lines string

var noteCmd = &cobra.Command{
	Use:   "note [lines to note on]",
	Short: "Write notes about your code",
	Long: `Use the note command to write notes on highlighted portions of your code
	You can also optionally use the --lines flag and specify the range of your code block to 
	comment on. The CLI will automatically detect your current workding directory
	and file`,

	Run: func(cmd *cobra.Command, args []string) {
		activeFileMetaData, err := getCurrentFileMetadata()
		if err != nil{
			fmt.Printf("Error %s", err)
			return
		}

		client = cmd.Context().Value("client").(*internal.APIService)
		databaseID = cmd.Context().Value("databaseID").(string)
		if lines != nil{

		}
		lines, err = setLines(activeFileMetaData["lines"].(map[string]interface{}))


		fmt.Println(activeFileMetaData)
		if err := internal.CreateDatabaseEntry(client, databaseID, activeFileMetaData, note, category); err != nil{
			fmt.Printf("Error: %s", err)
			return
		}
	},
}


func init() {
	noteCmd.Flags().StringVarP(&note, "message", "m", "", "Message (required)")
	noteCmd.MarkFlagRequired("message")
	noteCmd.Flags().StringVarP(&category, "category", "c", "", "Category of note")
	noteCmd.Flags().StringVarP(&lines, "lines", "l", "", "Lines to highlight")
	rootCmd.AddCommand(noteCmd)
}
