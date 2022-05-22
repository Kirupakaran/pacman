/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/kirupakaran/pacman/app"
	"github.com/spf13/cobra"
)

// unifyCmd represents the unify command
var unifyCmd = &cobra.Command{
	Use:   "unify [OPTIONS]",
	Short: "Unify package versions in package.json of base|repos. If no option is passed, unifies packages in base to a common patch version.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Unify(cmd.Flag("minor").Changed)
	},
}

func init() {
	rootCmd.AddCommand(unifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	unifyCmd.Flags().BoolP("minor", "", false, "Unifies packages based on minor version")
}
