/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"time"

	"github.com/ezcdlabs/ezcd/cmd/cli/service"
	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/spf13/cobra"
)

// commitStageStartedCmd represents the commitStageStarted command
var commitStageStartedCmd = &cobra.Command{
	Use:   "commit-stage-started",
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
		authorName, err := cmd.Flags().GetString("author-name")
		if err != nil {
			log.Fatalf("Failed to get 'author-name' flag: %v\n", err)
		}
		authorEmail, err := cmd.Flags().GetString("author-email")
		if err != nil {
			log.Fatalf("Failed to get 'author-email' flag: %v\n", err)
		}
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			log.Fatalf("Failed to get 'message' flag: %v\n", err)
		}
		dateString, err := cmd.Flags().GetString("date")
		if err != nil {
			log.Fatalf("Failed to get 'date' flag: %v\n", err)
		}
		var date time.Time
		date, err = time.Parse(time.RFC3339, dateString)
		if err != nil {
			date, err = time.Parse(time.DateTime, dateString)
		}
		if err != nil {
			log.Fatalf("Failed to parse 'date' flag: %v\n", err)
		}

		commitData := ezcd.CommitData{
			Hash:        hash,
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			Message:     message,
			Date:        date,
		}

		ezcdService := service.Get()
		err = ezcdService.CommitStageStarted(project, commitData)

		if err != nil {
			log.Fatalf("Failed to add commit: %v\n", err)
		}

		log.Printf("Commit added: %v\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(commitStageStartedCmd)

	commitStageStartedCmd.Flags().StringP("project", "P", "", "The project slug")
	commitStageStartedCmd.Flags().StringP("hash", "H", "", "The commit hash")
	commitStageStartedCmd.Flags().StringP("author-name", "N", "", "The commit author's name")
	commitStageStartedCmd.Flags().StringP("author-email", "E", "", "The commit author's email")
	commitStageStartedCmd.Flags().StringP("message", "M", "", "The commit message")
	commitStageStartedCmd.Flags().StringP("date", "D", "", "The commit date")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitStageStartedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitStageStartedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
