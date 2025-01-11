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
	var start int
	var end int
	if strings.Contains(line_range, "-"){
		parts := strings.Split(line_range, "-")
		start, _ = strconv.Atoi(parts[0])
		end, _ = strconv.Atoi(parts[1])
	} else {
		start, _ = strconv.Atoi(line_range)
		end = start
	}
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

// func createHelpTemplate()(models.HelpTemplate){
// 	return models.HelpTemplate{
// 		Usage: "./noted note [flags]",
// 		Description: "Use the 'note' command to create notes on your code. Highlight lines or specify a range using the '--lines' flag. All notes will be categorized and stored in Notion.",
// 		Flags: []string {
// 			"-m, --message   (Required) Specify the note you'd like to write.",
// 			"-l, --lines     (Optional) Specify the range of lines to comment on.",
// 			"-c, --category  (Optional) Specify a category for your note.",
// 		},
// 		Examples: []string{
// 			"./noted note -m \"Example note\"",
// 			"./noted note -m \"Fix bug in loop\" -l 10-20",
// 			"./noted note -m \"Refactor suggestion\" -c \"Improvement\"",
// 		},
// 		Notes: []string{
// 			"Use quotation marks (\" \") to wrap your note message.",
// 			"Ensure the Current File Tracker Extension is active to detect highlighted code blocks automatically.",
// 		},
// 	}
// }


var note string
var category string
var lines string
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Write notes about your code",
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
	noteCmd.Flags().StringVarP(&category, "category", "c", "None", "Category of note (Optional)")
	noteCmd.Flags().StringVarP(&lines, "lines", "l", "", "Lines to highlight (Optional)")

	noteCmdHelpTemplate := internal.CreateHelpTemplate(
		"./noted note [flags]",
		"Use the 'note' command to create notes on your code. Highlight lines or specify a range using the '--lines' flag. All notes will be categorized and stored in Notion.",
		[]string{
			"-m, --message   (Required) Specify the note you'd like to write.",
			"-l, --lines     (Optional) Specify the range of lines to comment on.",
			"-c, --category  (Optional) Specify a category for your note.",
		},
		[]string{
			"./noted note -m \"Example note\"",
			"./noted note -m \"Fix bug in loop\" -l 10-20",
			"./noted note -m \"Refactor suggestion\" -c \"Improvement\"",
		},
		[]string{
			"Use quotation marks (\" \") to wrap your note message.",
			"Ensure the Current File Tracker Extension is active to detect highlighted code blocks automatically.",
		},
	)
	noteCmd.SetHelpTemplate(internal.GenerateHelpMessage(noteCmdHelpTemplate))
	rootCmd.AddCommand(noteCmd)
}
