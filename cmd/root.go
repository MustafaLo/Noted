/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"github.com/MustafaLo/noted/internal"
	"github.com/MustafaLo/noted/models"
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmdHelpTemplate := internal.CreateHelpTemplate(
		"~~~~~~~~~~~~~~~~~ Root Help ~~~~~~~~~~~~~~~~~~~",
		"No default behavior (see below for flags)",
		"Noted: A Notion-integrated CLI for effortless, in-context educational notes during development.",
		[]string{
			"note: Write notes about your code",
			"list: View notes about your code",
			"insights: Generate AI summaries about your notes",
		},
		[]string{
			"Run ./noted note -h OR ./noted note --help for more information",
			"Run ./noted list -h OR ./noted list --help for more information",
			"Run ./noted insights -h OR ./noted insights --help for more information",	
		},
		[]string{
			"Make sure to download the 'Track Current File' extension and run from your command palette ('Current File Tracker')",
			"Follow the directions on 'developers.notion.com/docs/create-a-notion-integration': ",
			"   - Create a Notion API Key (developers.notion.com) and insert into your env file as \"NOTION_API_KEY=\"",
			"   - Share a Notion Page with your newly created integration",
			"   - Grab that page id (developers.notion.com/docs/working-with-page-content) and insert into your env file as \"NOTION_PAGE_ID\"",
		},
	)
	rootCmd.SetHelpTemplate(internal.GenerateHelpMessage(rootCmdHelpTemplate))
}


