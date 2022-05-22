/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/kirupakaran/pacman/app"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates package.json in the given repo",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a directory as an argument")
		}
		if app.IsValidDir(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid directory: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		app.Update(args[0])
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
