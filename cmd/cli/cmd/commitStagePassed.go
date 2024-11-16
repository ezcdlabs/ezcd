/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ezcdlabs/ezcd/cmd/cli/service"
	"github.com/spf13/cobra"
)

// commitStagePassedCmd represents the commitStagePassed command
var commitStagePassedCmd = &cobra.Command{
	Use:   "commit-stage-passed",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Fatalf("Failed to get 'project' flag: %v\n", err)
		}
		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			log.Fatalf("Failed to get 'hash' flag: %v\n", err)
		}

		ezcdService := service.Get()
		err = ezcdService.CommitStagePassed(project, hash)

		if err != nil {
			log.Fatalf("Failed to set commit stage as passed: %v\n", err)
		}

		log.Printf("Commit stage set as passed: %v\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(commitStagePassedCmd)

	commitStagePassedCmd.Flags().StringP("project", "P", "", "The project slug")
	commitStagePassedCmd.Flags().StringP("hash", "H", "", "The commit hash")
}
