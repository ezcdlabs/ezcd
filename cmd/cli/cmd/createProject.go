/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCreateProjectCommand(serviceLoader EzcdServiceLoader) *cobra.Command {
	return &cobra.Command{
		Use:   "create-project",
		Args:  cobra.ExactArgs(1),
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ezcdService, err := serviceLoader.Load()
			if err != nil {
				// tell cobra there was an error:
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to create project: %w", err)
			}

			err = ezcdService.CreateProject(args[0])
			if err != nil {
				// tell cobra there was an error:
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to create project: %w", err)
			}

			fmt.Printf("Project %s created\n", args[0])

			return nil
		},
	}
}
