/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		vprBuf := viper.GetViper()
		// Bind subcommand-specific cobra flags to viper
		err := vprBuf.BindPFlags(cmd.Flags())
		if err != nil {
			logrus.WithField("cobra-cmd", cmd.Use).Errorf("error binding flags: %v", err)
			os.Exit(1)
		}

		err = UnmarshalSubMerged(vprBuf, cmd.Parent().Use+"."+cmd.Use, &vprFlgsDetails)
		if err != nil {
			logrus.Fatalf("failed to unmarshal version config: %v", err)
		}

		return nil
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
		logrus.WithField("cobra-cmd", cmd.Use).Infof("Persistent flags from rootCmd")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag1: %s", vprFlgsRoot.RootFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag2: %s", vprFlgsRoot.RootFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag3: %s", vprFlgsRoot.RootFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag4: %s", vprFlgsRoot.RootFlag4)
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
