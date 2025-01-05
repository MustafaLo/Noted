/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
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






// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note [lines to note on]",
	Short: "Write notes about your code",
	Long: `Use the note command to write notes on highlighted portions of your code
	You can also optionally use the --lines flag and specify the range of your code block to 
	comment on. The CLI will automatically detect your current workding directory
	and file`,


	//Use the Flag StringVarP command to directly insert line range into variable
	//as opposed to parsing through the argument string array. Cobra will automatically
	//handle flag parsing

	Run: func(cmd *cobra.Command, args []string) {
		activeFileMetaData, err := getCurrentFileMetadata()
		if err != nil{
			fmt.Printf("Error %s", err)
			return
		}
		printFileMetaData(activeFileMetaData)
	},
}


func init() {
	rootCmd.AddCommand(noteCmd)
}
