/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MustafaLo/Noted/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "View notes on active file",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
		activeFileMetaData, err:= internal.GetCurrentFileMetadata()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
