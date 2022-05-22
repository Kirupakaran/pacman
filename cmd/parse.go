/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/kirupakaran/pacman/app"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses package.json files under all top-level sub-directories in the given directory",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("dir").Changed {
			if app.IsValidDir(cmd.Flag("dirPath").Value.String()) {
				return nil
			}
			return fmt.Errorf("invalid directory: %s", args[0])
		} else if cmd.Flag("repos").Changed {
			if app.IsValidFile(cmd.Flag("repoList").Value.String()) {
				return nil
			}
			return fmt.Errorf("invalid file: %s", args[0])
		} else {
			return fmt.Errorf("either repo list or directory path required")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("dir").Changed {
			app.Parse(cmd.Flag("dirPath").Value.String())
		} else if cmd.Flag("repos").Changed {
			app.ParseByRepo(cmd.Flag("repoList").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	parseCmd.Flags().BoolP("repos", "r", false, "If you want to pass a list of repos; must set GITHUB_PAT env")
	parseCmd.Flags().StringP("repoList", "", "repoList", "Pass a file containing list of repos")

	parseCmd.Flags().BoolP("dir", "d", false, "If you want to pass a directory")
	parseCmd.Flags().StringP("dirPath", "", "", "Pass a directory containing multiple sub-directories of node repos")

	parseCmd.MarkFlagsRequiredTogether("repos", "repoList")
	parseCmd.MarkFlagsRequiredTogether("dir", "dirPath")
	parseCmd.MarkFlagsMutuallyExclusive("repos", "dir")
}
