/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/spf13/cobra"
)

func NewCommitStageStartedCommand(serviceLoader EzcdServiceLoader) *cobra.Command {

	// CommitStageStartedCmd represents the commitStageStarted command
	commitPhaseStartedCmd := &cobra.Command{
		Use:   "commit-stage-started",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			project, err := cmd.Flags().GetString("project")
			if err != nil || project == "" {
				return fmt.Errorf("failed to get 'project' flag: %v", err)
			}
			hash, err := cmd.Flags().GetString("hash")
			if err != nil || hash == "" {
				// why is this line not covered? I can see this error in the test output
				return fmt.Errorf("failed to get 'hash' flag: %v", err)
			}
			authorName, err := cmd.Flags().GetString("author-name")
			if err != nil || authorName == "" {
				return fmt.Errorf("failed to get 'author-name' flag: %v", err)
			}
			authorEmail, err := cmd.Flags().GetString("author-email")
			if err != nil || authorEmail == "" {
				return fmt.Errorf("failed to get 'author-email' flag: %v", err)
			}
			message, err := cmd.Flags().GetString("message")
			if err != nil || message == "" {
				return fmt.Errorf("failed to get 'message' flag: %v", err)
			}
			dateString, err := cmd.Flags().GetString("date")
			if err != nil || dateString == "" {
				return fmt.Errorf("failed to get 'date' flag: %v", err)
			}
			var date time.Time
			date, err = time.Parse(time.RFC3339, dateString)
			if err != nil {
				date, err = time.Parse(time.DateTime, dateString)
			}
			if err != nil {
				return fmt.Errorf("failed to parse 'date' flag: %v", err)
			}

			commitData := ezcd.CommitData{
				Hash:        hash,
				AuthorName:  authorName,
				AuthorEmail: authorEmail,
				Message:     message,
				Date:        date,
			}

			ezcdService, err := serviceLoader.Load()
			if err != nil {
				return fmt.Errorf("failed to get ezcd service: %v", err)
			}
			err = ezcdService.CommitStageStarted(project, commitData)
			if err != nil {
				return fmt.Errorf("failed to add commit: %v", err)
			}

			log.Printf("Commit added: %v\n", hash)
			return nil
		},
	}

	commitPhaseStartedCmd.Flags().StringP("project", "P", "", "The project slug")
	commitPhaseStartedCmd.Flags().StringP("hash", "H", "", "The commit hash")
	commitPhaseStartedCmd.Flags().StringP("author-name", "N", "", "The commit author's name")
	commitPhaseStartedCmd.Flags().StringP("author-email", "E", "", "The commit author's email")
	commitPhaseStartedCmd.Flags().StringP("message", "M", "", "The commit message")
	commitPhaseStartedCmd.Flags().StringP("date", "D", "", "The commit date")

	return commitPhaseStartedCmd
}
