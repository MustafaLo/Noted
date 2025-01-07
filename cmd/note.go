/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	// "regexp"
	// "strconv"
	// "strings"

	"github.com/MustafaLo/Noted/internal"
	"github.com/spf13/cobra"
)

func getCurrentFileMetadata() (FileMetadata, error) {
    var metadata FileMetadata
    data, err := os.ReadFile("fileMetadata.json")
    if err != nil {
        return metadata, fmt.Errorf("failed to open fileMetaData -- make sure to enable File Tracker Extension")
    }

    err = json.Unmarshal(data, &metadata)
    if err != nil {
        return metadata, fmt.Errorf("failed to parse fileMetaData: %w", err)
    }

    return metadata, nil
}

func printFileMetaData(metadata FileMetadata) {
    fmt.Printf("File Name: %s\n", metadata.FileName)
    fmt.Printf("File Path: %s\n", metadata.FilePath)
    fmt.Printf("Lines: Start=%d, End=%d\n", metadata.Lines.Start, metadata.Lines.End)
    fmt.Printf("Timestamp: %s\n", metadata.Timestamp)
}



// func getCurrentFileMetadata()(map[string]interface{}, error){
// 	var metadata map[string]interface{}
// 	data, err := os.ReadFile("fileMetadata.json")
// 	if err != nil{
// 		return nil, fmt.Errorf("failed to open fileMetaData -- make sure to enable File Tracker Extension")
// 	}

// 	ok := json.Unmarshal(data, &metadata)
// 	if ok != nil{
// 		return nil, fmt.Errorf("failed to parse fileMetaData")
// 	}
// 	return metadata, nil
// }


// func printFileMetaData(data map[string]interface{}){
// 	for k, v := range data{
// 		if k == "lines" && v != nil{
// 			if linesMap, ok := v.(map[string]interface{}); ok {
// 				fmt.Printf("key: %s, value: {", k)
// 				for lk, lv := range linesMap {
// 					fmt.Printf("%s: %0.f,", lk, lv) // Adjust formatting as needed
// 				}
// 				fmt.Println("}")
// 			} else {
// 				fmt.Printf("Not okay")
// 			}
// 		} else{
// 			fmt.Printf("key: %s, value: %s\n", k, v)
// 		}
// 	}
// }

// func setLines(fileMetaData map[string]interface{})(string, error){
// 	//Check for no lines set
// 	if lines == "" && fileMetaData["lines"] == nil{
// 		return "None", nil
// 	} else if lines != "" {
// 		err := isValidLinesFormat(lines)
// 		if err != nil{
// 			return "", err
// 		}
// 		return lines, nil
// 	} 
// 	var highlighted_lines string
	
// 	lines_struct, ok := fileMetaData["lines"].(map[string]interface{})
// 	if !ok {
// 		return "", fmt.Errorf("highlighted lines is not in a proper format -- check active file extension")
// 	}

// 	start, start_exists := lines_struct["start"]
// 	end, end_exists := lines_struct["end"]

// 	if !start_exists || !end_exists {
// 		return "", fmt.Errorf("highlighted lines is not in a proper format -- check active file extension")
// 	}
// 	if start == end{
// 		highlighted_lines = 
// 	}


// }

// func isValidLinesFormat(line_range string)(error){
// 	re := regexp.MustCompile(`^\d+$|^\d+-\d+$`)

// 	if !re.MatchString(line_range) {
// 		return fmt.Errorf("invalid format for --lines: must be a number or range (e.g., 5 or 5-12)")
// 	}

// 	if strings.Contains(line_range, "-"){
// 		parts := strings.Split(line_range, "-")
// 		start, _ := strconv.Atoi(parts[0])
// 		end, _ := strconv.Atoi(parts[1])
// 		if start >= end{
// 			return fmt.Errorf("invalid range: start must be less than end (e.g., 5-12)")
// 		}
// 	}

// 	return nil
// }


var note string
var client *internal.APIService
var databaseID string
var category string
var lines string

type FileMetadata struct {
    FileName  string `json:"fileName"`
    FilePath  string `json:"filePath"`
    Lines     struct {
        Start int `json:"start"`
        End   int `json:"end"`
    } `json:"lines"`
    Timestamp string `json:"timestamp"`
}


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

		// if lines != ""{
		// 	err := isValidLinesFormat()
		// 	if err != nil{
		// 		fmt.Printf("Error %s", err)
		// 		return
		// 	}
		// } else{
		// 	lines_struct, ok := activeFileMetaData["lines"].(map[string]interface{})
		// 	if !ok{
		// 		fmt.Printf("lines is nil or not in a proper format -- check extension")
		// 		return
		// 	}
		// 	lines, err = setLines(lines_struct)
		// 	if err != nil{
		// 		fmt.Printf("Error %s", err)
		// 		return
		// 	}
		// }


		printFileMetaData(activeFileMetaData)
		// if err := internal.CreateDatabaseEntry(client, databaseID, activeFileMetaData, note, category); err != nil{
		// 	fmt.Printf("Error: %s", err)
		// 	return
		// }
	},
}


func init() {
	noteCmd.Flags().StringVarP(&note, "message", "m", "", "Message (required)")
	noteCmd.MarkFlagRequired("message")
	noteCmd.Flags().StringVarP(&category, "category", "c", "", "Category of note")

	noteCmd.Flags().StringVarP(&lines, "lines", "l", "", "Lines to highlight")
	rootCmd.AddCommand(noteCmd)
}
