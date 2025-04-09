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

var versionPersistentFlag1 string
var versionPersistentFlag2 string
var versionPersistentFlag3 string
var versionPersistentFlag4 string

type ViperFlagsVersion struct {
	VersionFlag1           string `mapstructure:"versionflag1"`
	VersionFlag2           string `mapstructure:"versionflag2"`
	VersionFlag3           string `mapstructure:"versionflag3"`
	VersionFlag4           string `mapstructure:"versionflag4"`
	VersionPersistentFlag1 string `mapstructure:"versionpersistentflag1"`
	VersionPersistentFlag2 string `mapstructure:"versionpersistentflag2"`
	VersionPersistentFlag3 string `mapstructure:"versionpersistentflag3"`
	VersionPersistentFlag4 string `mapstructure:"versionpersistentflag4"`
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithField("cobra-cmd", cmd.Use).Infof("version subcommand called")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag1: %s", vprFlgsVersion.VersionFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag2: %s", vprFlgsVersion.VersionFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag3: %s", vprFlgsVersion.VersionFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("versionflag4: %s", vprFlgsVersion.VersionFlag4)

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
	cobra.OnInitialize(initConfigVersion)
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")
	versionCmd.PersistentFlags().StringVar(&versionPersistentFlag1, "versionpersistentflag1", "value from default", "version persistent flag 1")
	versionCmd.PersistentFlags().StringVar(&versionPersistentFlag2, "versionpersistentflag2", "value from default", "version persistent flag 2")
	versionCmd.PersistentFlags().StringVar(&versionPersistentFlag3, "versionpersistentflag3", "value from default", "version persistent flag 3")
	versionCmd.PersistentFlags().StringVar(&versionPersistentFlag4, "versionpersistentflag4", "value from default", "version persistent flag 4")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	versionCmd.Flags().StringVar(&versionFlag1, "versionflag1", "value from default", "version flag 1")
	versionCmd.Flags().StringVar(&versionFlag2, "versionflag2", "value from default", "version flag 2")
	versionCmd.Flags().StringVar(&versionFlag3, "versionflag3", "value from default", "version flag 3")
	versionCmd.Flags().StringVar(&versionFlag4, "versionflag4", "value from default", "version flag 4")
}

func UnmarshalSubMerged(v *viper.Viper, section string, target any) error {
	// 1. Skip if no config file is loaded at all
	if v.ConfigFileUsed() == "" {
		return v.Unmarshal(target) // only env, flags, defaults
	}

	// 2. Extract the subsection of the config file
	sub := v.GetStringMap(section)
	if len(sub) == 0 {
		// No subsection found, fallback to flags/env/default
		return v.Unmarshal(target)
	}

	// 3. Merge section into Viper's config layer (not override!)
	if err := v.MergeConfigMap(sub); err != nil {
		return fmt.Errorf("failed to merge config section '%s': %w", section, err)
	}

	// 4. Now unmarshal with proper priority:
	// flags > env > merged config > defaults
	return v.Unmarshal(target)
}

func initConfigVersion() {
	vprBuf := viper.GetViper()
	// Bind subcommand-specific cobra flags to viper
	err := vprBuf.BindPFlags(versionCmd.Flags())
	if err != nil {
		logrus.WithField("cobra-cmd", versionCmd.Use).Errorf("error binding flags: %v", err)
		os.Exit(1)
	}

	err = UnmarshalSubMerged(vprBuf, versionCmd.Use, &vprFlgsVersion)
	if err != nil {
		logrus.Fatalf("failed to unmarshal version config: %v", err)
	}
}
