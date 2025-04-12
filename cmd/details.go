/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var detailsFlag1 string
var detailsFlag2 string
var detailsFlag3 string
var detailsFlag4 string

type ViperFlagsDetails struct {
	DetailsFlag1 string `mapstructure:"detailsflag1"`
	DetailsFlag2 string `mapstructure:"detailsflag2"`
	DetailsFlag3 string `mapstructure:"detailsflag3"`
	DetailsFlag4 string `mapstructure:"detailsflag4"`
}

var vprFlgsDetails ViperFlagsDetails

// detailsCmd represents the details command
var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitViperSubCmd(viper.GetViper(), cmd, &vprFlgsDetails)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("details called")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("detailsflag1: %s", vprFlgsDetails.DetailsFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("detailsflag2: %s", vprFlgsDetails.DetailsFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("detailsflag3: %s", vprFlgsDetails.DetailsFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("detailsflag4: %s", vprFlgsDetails.DetailsFlag4)

		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("flags from subcommand version")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag1: %s", vprFlgsVersion.VersionFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag2: %s", vprFlgsVersion.VersionFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag3: %s", vprFlgsVersion.VersionFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag4: %s", vprFlgsVersion.VersionFlag4)

		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("persistent flags from subcommand version")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionpersistentflag1: %s", vprFlgsVersion.VersionPersistentFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionpersistentflag2: %s", vprFlgsVersion.VersionPersistentFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionpersistentflag3: %s", vprFlgsVersion.VersionPersistentFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionpersistentflag4: %s", vprFlgsVersion.VersionPersistentFlag4)

		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("Persistent flags from rootCmd")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag1: %s", vprFlgsRoot.RootPersistentFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag2: %s", vprFlgsRoot.RootPersistentFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag3: %s", vprFlgsRoot.RootPersistentFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag4: %s", vprFlgsRoot.RootPersistentFlag4)
	},
}

func init() {
	versionCmd.AddCommand(detailsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// detailsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// detailsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	detailsCmd.Flags().StringVar(&detailsFlag1, "detailsflag1", "value from default", "version flag 1")
	detailsCmd.Flags().StringVar(&detailsFlag2, "detailsflag2", "value from default", "version flag 2")
	detailsCmd.Flags().StringVar(&detailsFlag3, "detailsflag3", "value from default", "version flag 3")
	detailsCmd.Flags().StringVar(&detailsFlag4, "detailsflag4", "value from default", "version flag 4")
}
