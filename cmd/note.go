/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"path/filepath"
	"os"
)


func getCurrentFileName()(string){
	_, filename, _, ok := runtime.Caller(0)
	if !ok{
		panic("Failed to get current file")
	}
	basefilename := filepath.Base(filename)
	return basefilename
}

func getFileContent(filename string)(){
	content, err := os.ReadFile(filename)
	if err != nil{
		panic("Error!")
	}
	fmt.Println(string(content))
}



// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note [lines to note on]",
	Short: "Write notes about your code",
	Long: `Use the note command to write notes on line ranges of your code.
	You must use the --lines flag and specify the range of your code block to 
	comment on. The CLI will automatically detect your current workding directory
	and file`,

	//Use the Flag StringVarP command to directly insert line range into variable
	//as opposed to parsing through the argument string array. Cobra will automatically
	//handle flag parsing

	Run: func(cmd *cobra.Command, args []string) {
		file := getCurrentFileName()
		getFileContent((file))
		fmt.Println("Noted!")
	},
}



func init() {
	rootCmd.AddCommand(noteCmd)
}
