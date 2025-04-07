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

		// If a config file is not found, log a trace error. Otherwise, read it in.
		if err := vprBuf.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				logrus.WithField("cobra-cmd", cmd.Use).Infof("No config file found; continue with cobra default values")
				if err := vprBuf.Unmarshal(&vprFlgsVersion); err != nil {
					logrus.WithField("cobra-cmd", cmd.Use).WithError(err).Fatal("versionCmd: failed to unmarshal viper config")
				}
			} else {
				// Config file was found but another error occurred
				fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Load the config files values that are bound to version Cobra CLI flags
			if err := vprBuf.Sub("version").Unmarshal(&vprFlgsVersion); err != nil {
				logrus.WithError(err).Fatal("versionCmd: failed to unmarshal viper config")
			}

			//
			if err := vprBuf.Unmarshal(&vprFlgsVersion); err != nil {
				logrus.WithError(err).Fatal("versionCmd: failed to unmarshal viper config")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithField("cobra-cmd", cmd.Use).Infof("version called")
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

	versionCmd.Flags().StringVar(&versionFlag1, "versionflag1", "from default", "version flag 1")
	versionCmd.Flags().StringVar(&versionFlag2, "versionflag2", "from default", "version flag 2")
	versionCmd.Flags().StringVar(&versionFlag3, "versionflag3", "from default", "version flag 3")
	versionCmd.Flags().StringVar(&versionFlag4, "versionflag4", "from default", "version flag 4")
}
