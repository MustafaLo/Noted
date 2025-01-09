/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/spf13/cobra"
)



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

func getCodeBlock(filePath string, lines string)(string, error){
	file, err := os.Open(filePath)
	if err != nil{
		return "", fmt.Errorf("error opening active file: %w", err)
	}
	defer file.Close()


	var codeLines []string
	scanner := bufio.NewScanner(file)
	currentLine := 1
	start, end := splitLineParts(lines)

	for scanner.Scan(){
		if currentLine >= start && currentLine <= end{
			codeLines = append(codeLines, scanner.Text())
		}
		if currentLine > end{
			break
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil{
		return "", fmt.Errorf("error reading file: %w", err)
	}

	codeBlock := strings.Join(codeLines, "\n")
	if strings.TrimSpace(codeBlock) == ""{
		return "", fmt.Errorf("code block specified was empty")
	}

	return codeBlock, nil
}


var note string
var category string
var lines string

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Write notes about your code",
	Long: `Use the note command to write notes on highlighted portions of your code.
	Subcommands:
	* --message (-m): Required flag that you should use to specify the note you'd like to write
	    * Example Usage: 
			- ./noted note -m "Example note"

	* --lines (-l): Optional flag you can use to specify the range of your code block to comment on. The Current File Tracker 
					Extension will automatically detect highlighted code blocks on your active file
	   * Example Usage: 
	        - ./noted note -m "Example note" -l 15-20
			- ./noted note -m "Example note" -l 25
		
	* --category (-c): Optional flag you can use to specify an existing or new category that your note falls under
	   * Example Usage:
			- ./noted note -m "Example note" -c "Syntax"
			- ./noted note -m "Example note" -c "Design"`,

	Run: func(cmd *cobra.Command, args []string) {
		activeFileMetaData, err := internal.GetCurrentFileMetadata()
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
			codeBlock, err := getCodeBlock(activeFileMetaData.FilePath, lines)
			if err != nil{
				fmt.Printf("Error: %s", err)
				return
			} 
			err = internal.UpdateDatabaseEntry(client, pageID, codeBlock, activeFileMetaData.Language, note)
			if err != nil{
				fmt.Printf("Error: %s", err)
				return
			}

		}
	},
}


func init() {
	noteCmd.Flags().StringVarP(&note, "message", "m", "", "Message (required)")
	noteCmd.MarkFlagRequired("message")
	noteCmd.Flags().StringVarP(&category, "category", "c", "None", "Category of note")
	noteCmd.Flags().StringVarP(&lines, "lines", "l", "", "Lines to highlight")

	rootCmd.AddCommand(noteCmd)
}
