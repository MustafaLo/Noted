/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"github.com/MustafaLo/Noted/internal"
	"github.com/spf13/cobra"
)


func authenticate(env map[string]string)(*internal.APIService, error){
	if _, exists := env["NOTION_API_KEY"]; !exists{
		return nil, fmt.Errorf("env file needs authentication key")
	}

	client, err := internal.InitService(env["NOTION_API_KEY"]); 
	if err != nil{
		return nil, fmt.Errorf("error authenticating %w", err)
	}
	return client, nil
}

func intialize_db(s *internal.APIService, env map[string]string)(error){
	if _, exists := env["NOTION_PAGE_ID"]; !exists{
		return fmt.Errorf("env file needs shared parent page ID")
	}

	if err := internal.IntializeDatabase(s, env["NOTION_DATABASE_ID"], env["NOTION_PAGE_ID"]); err != nil{
		return fmt.Errorf("error intializing %w", err)
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

		if err := intialize_db(client, envMap); err != nil{
			fmt.Println(err)
			os.Exit(1)
		}

		//Passing in client and envMap to all sub commands
		ctx := context.WithValue(cmd.Context(), "client", client)
		ctx = context.WithValue(ctx, "databaseID", envMap["NOTION_DATABASE_ID"])
		cmd.SetContext(ctx)
	} ,


	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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


