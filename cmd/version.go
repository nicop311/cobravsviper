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

var versionFlag1 string
var versionFlag2 string
var versionFlag3 string
var versionFlag4 string

type ViperFlagsVersion struct {
	VersionFlag1 string `mapstructure:"versionflag1"`
	VersionFlag2 string `mapstructure:"versionflag2"`
	VersionFlag3 string `mapstructure:"versionflag3"`
	VersionFlag4 string `mapstructure:"versionflag4"`
}

var vprFlgsVersion ViperFlagsVersion

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A SUBcommand",
	Long: `A Cobra Subcommand

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

		err = UnmarshalSubMerged(vprBuf, "version", &vprFlgsVersion)
		if err != nil {
			logrus.Fatalf("failed to unmarshal version config: %v", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithField("cobra-cmd", cmd.Use).Infof("version subcommand called")
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
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	versionCmd.Flags().StringVar(&versionFlag1, "versionflag1", "value from default", "version flag 1")
	versionCmd.Flags().StringVar(&versionFlag2, "versionflag2", "value from default", "version flag 2")
	versionCmd.Flags().StringVar(&versionFlag3, "versionflag3", "value from default", "version flag 3")
	versionCmd.Flags().StringVar(&versionFlag4, "versionflag4", "value from default", "version flag 4")
}

func UnmarshalSubMerged(v *viper.Viper, section string, target any) error {
	if v == nil {
		return fmt.Errorf("viper instance is nil")
	}
	sub := v.GetStringMap(section)
	for k, val := range sub {
		v.Set(k, val) // inject sub keys into main viper
	}
	return v.Unmarshal(target)
}
