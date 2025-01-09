/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"github.com/MustafaLo/Noted/internal"
	"github.com/MustafaLo/Noted/models"
	"github.com/spf13/cobra"
)


func authenticate(env map[string]string)(*models.APIService, error){
	if _, exists := env["NOTION_API_KEY"]; !exists{
		return nil, fmt.Errorf("env file needs authentication key")
	}

	client, err := internal.InitService(env["NOTION_API_KEY"]); 
	if err != nil{
		return nil, fmt.Errorf("error authenticating %w", err)
	}
	return client, nil
}

func intialize_db(s *models.APIService, env map[string]string)(string, error){
	if _, exists := env["NOTION_PAGE_ID"]; !exists{
		return "", fmt.Errorf("env file needs shared parent page ID")
	}

	database_id, err := internal.IntializeDatabase(s, env["NOTION_DATABASE_ID"], env["NOTION_PAGE_ID"])
	if err != nil{
		return "", fmt.Errorf("error intializing %w", err)
	}
	return database_id, nil
}


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Noted",
	Short: "A CLI tool that connects to your Notion and allows you to keep a synchronized space for educational notes",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		envMap, err := internal.LoadEnv(); 
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := authenticate(envMap)
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}

		database_id, err := intialize_db(client, envMap)
		if err != nil{
			fmt.Println(err)
			os.Exit(1)
		}

		//Passing in client and envMap to all sub commands
		ctx := context.WithValue(cmd.Context(), "client", client)
		ctx = context.WithValue(ctx, "databaseID", database_id)
		cmd.SetContext(ctx)
	} ,
}

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


