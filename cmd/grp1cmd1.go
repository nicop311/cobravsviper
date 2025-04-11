/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// grp1cmd1Cmd represents the grp1cmd1 command
var grp1cmd1Cmd = &cobra.Command{
	Use:     "grp1cmd1",
	Short:   "A brief description of your command",
	GroupID: "group1",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grp1cmd1 called")
	},
}

func init() {
	rootCmd.AddCommand(grp1cmd1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grp1cmd1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grp1cmd1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
