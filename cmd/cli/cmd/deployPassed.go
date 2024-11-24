/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func NewDeployPassedCommand(serviceLoader EzcdServiceLoader) *cobra.Command {
	deployPassedCmd := &cobra.Command{
		Use:   "deploy-passed",
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
			hash, err := cmd.Flags().GetString("hash")
			if err != nil || hash == "" {
				return fmt.Errorf("failed to get 'hash' flag: %v", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			project, _ := cmd.Flags().GetString("project")
			hash, _ := cmd.Flags().GetString("hash")
			cmd.SilenceUsage = true

			ezcdService, err := serviceLoader.Load()
			if err != nil {
				return fmt.Errorf("failed to get ezcd service: %v", err)
			}

			err = ezcdService.DeployPassed(project, hash) //DeployPassed(project, hash)

			if err != nil {
				return fmt.Errorf("failed to set deploy as passed: %v   %v", err, project)
			}

			log.Printf("deploy set as passed: %v\n", hash)
			return nil
		},
	}

	deployPassedCmd.Flags().StringP("project", "P", "", "The project slug")
	deployPassedCmd.Flags().StringP("hash", "H", "", "The commit hash")

	return deployPassedCmd
}
