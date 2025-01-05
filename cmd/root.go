/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/MustafaLo/Noted/internal"
	"github.com/spf13/cobra"
)


func authenticate()(error){
	if err := internal.InitService(); err != nil{
		return fmt.Errorf("Error authenticating %w", err)
	}
	return nil
}

func intialize_db()(error){
	if err := internal.IntializeDatabase(); err != nil{
		return fmt.Errorf("Error intializing %w", err)
	}
	return nil
}


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Noted",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := authenticate(); err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
		if err := intialize_db(); err != nil{
			fmt.Println(err)
			os.Exit(1)
		}
	} ,


	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Noted.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


