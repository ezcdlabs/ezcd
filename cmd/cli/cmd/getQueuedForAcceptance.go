/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGetQueuedForAcceptanceCommand(serviceLoader EzcdServiceLoader) *cobra.Command {
	getQueuedForAcceptanceCmd := &cobra.Command{
		Use:   "get-queued-for-acceptance",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")
			if err != nil || project == "" {
				return fmt.Errorf("failed to get 'project' flag: %v", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			project, _ := cmd.Flags().GetString("project")
			cmd.SilenceUsage = true

			ezcdService, err := serviceLoader.Load()
			if err != nil {
				return fmt.Errorf("failed to get ezcd service: %v", err)
			}

			commit, err := ezcdService.GetQueuedForAcceptance(project)

			if err != nil {
				return fmt.Errorf("failed to get commit queued for acceptance: %v", err)
			}

			if commit == nil {
				return fmt.Errorf("no commit queued for acceptance")
			}

			fmt.Fprintln(cmd.OutOrStdout(), commit.Hash)
			return nil
		},
	}

	getQueuedForAcceptanceCmd.Flags().StringP("project", "P", "", "The project slug")

	return getQueuedForAcceptanceCmd
}
