/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	// "go/scanner"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/spf13/cobra"
)

func getCurrentFileMetadata() (models.FileMetadata, error) {
    var metadata models.FileMetadata
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

func printFileMetaData(metadata models.FileMetadata) {
    fmt.Printf("File Name: %s\n", metadata.FileName)
    fmt.Printf("File Path: %s\n", metadata.FilePath)
    fmt.Printf("Lines: Start=%d, End=%d\n", metadata.Lines.Start, metadata.Lines.End)
    fmt.Printf("Timestamp: %s\n", metadata.Timestamp)
}


func setLines(highlighted_start int, highlighted_end int)(string, error){
	//Check for no lines set
	if lines == "" && highlighted_start == 0 && highlighted_end == 0{
		return "None", nil
	} else if lines != "" {
		err := isValidLinesFormat(lines)
		if err != nil{
			return "", err
		}
		return lines, nil
	} 
	var highlighted_lines string

	if highlighted_start == highlighted_end{
		highlighted_lines = strconv.Itoa(highlighted_start)
	} else {
		highlighted_lines = strconv.Itoa(highlighted_start) + "-" + strconv.Itoa(highlighted_end)
	}
	
	err := isValidLinesFormat(highlighted_lines)
	if err != nil{
		return "", err
	}

	return highlighted_lines, nil
}

func splitLineParts(line_range string)(int, int){
	parts := strings.Split(line_range, "-")
	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])
	return start, end
}

func isValidLinesFormat(line_range string)(error){
	re := regexp.MustCompile(`^\d+$|^\d+-\d+$`)

	if !re.MatchString(line_range) {
		return fmt.Errorf("invalid format for lines: must be a number or range (e.g., 5 or 5-12)")
	}

	if strings.Contains(line_range, "-"){
		start, end := splitLineParts(line_range)
		if start >= end{
			return fmt.Errorf("invalid range: start must be less than end (e.g., 5-12)")
		}
	}

	return nil
}

func getCodeBlock(fileName string, lines string)(string, error){
	file, err := os.Open(fileName)
	if err != nil{
		return "", fmt.Errorf("error opening active file: %w", err)
	}
	defer file.Close()

	
	scanner := bufio.NewScanner(file)
	currentLine := 1
	var codeLines []string

	for scanner.Scan(){
		if currentLine >= start && currentLine <= end{
			codeLines = append(codeLines, scanner.Text())
		}

		if currentLine > end{
			break
		}

		currentLine ++
	}

	if err := scanner.Err(); err != nil{
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return strings.Join(codeLines, "\n"), nil
}


var note string
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

		client := cmd.Context().Value("client").(*models.APIService)
		databaseID := cmd.Context().Value("databaseID").(string)
		
		lines, err = setLines(activeFileMetaData.Lines.Start, activeFileMetaData.Lines.End)
		if err != nil{
			fmt.Printf("Error %s", err)
			return
		}

		pageID, err := internal.CreateDatabaseEntry(client, databaseID, activeFileMetaData, note, lines, category)
		if err != nil{
			fmt.Printf("Error: %s", err)
			return
		}

		if lines != "None"{
			codeBlock, err := getCodeBlock(activeFileMetaData.FileName, lines)
		}


		printFileMetaData(activeFileMetaData)
		fmt.Println(pageID)
	},
}


func init() {
	noteCmd.Flags().StringVarP(&note, "message", "m", "", "Message (required)")
	noteCmd.MarkFlagRequired("message")
	noteCmd.Flags().StringVarP(&category, "category", "c", "None", "Category of note")

	noteCmd.Flags().StringVarP(&lines, "lines", "l", "", "Lines to highlight")
	rootCmd.AddCommand(noteCmd)
}
