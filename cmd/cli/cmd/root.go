/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/spf13/cobra"
)

// define interface that loads the ezcd service
type EzcdServiceLoader interface {
	Load() (ezcd.Ezcd, error)
}

func NewRootCmd(version string, serviceLoader EzcdServiceLoader) *cobra.Command {

	// RootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "ezcd-cli",
		Short: "A CLI tool for reporting events from your CI/CD pipeline to ezCD",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			// if the version flag is set, print the version and exit
			if cmd.Flag("version").Changed {
				cmd.Println(version)
			} else {
				cmd.Help()
			}
		},
	}

	// add flags
	rootCmd.Flags().BoolP("version", "v", false, "Print the version")

	// add sub commands
	rootCmd.AddCommand(NewCreateProjectCommand(serviceLoader))
	rootCmd.AddCommand(NewCommitStageStartedCommand(serviceLoader))
	rootCmd.AddCommand(NewCommitStagePassedCommand(serviceLoader))
	rootCmd.AddCommand(NewAcceptanceStageStartedCommand(serviceLoader))

	return rootCmd
}
