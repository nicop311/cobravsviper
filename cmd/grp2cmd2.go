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

var grp2cmd2Flag1 string
var grp2cmd2Flag2 string
var grp2cmd2Flag3 string
var grp2cmd2Flag4 string

var grp2cmd2PersistentFlag1 string
var grp2cmd2PersistentFlag2 string
var grp2cmd2PersistentFlag3 string
var grp2cmd2PersistentFlag4 string

type ViperGrp2cmd2 struct {
	Grp2cmd2Flag1 string `mapstructure:"grp2cmd2flag1"`
	Grp2cmd2Flag2 string `mapstructure:"grp2cmd2flag2"`
	Grp2cmd2Flag3 string `mapstructure:"grp2cmd2flag3"`
	Grp2cmd2Flag4 string `mapstructure:"grp2cmd2flag4"`

	Grp2cmd2PersistentFlag1 string `mapstructure:"grp2cmd2persistentflag1"`
	Grp2cmd2PersistentFlag2 string `mapstructure:"grp2cmd2persistentflag2"`
	Grp2cmd2PersistentFlag3 string `mapstructure:"grp2cmd2persistentflag3"`
	Grp2cmd2PersistentFlag4 string `mapstructure:"grp2cmd2persistentflag4"`
}

var vprFlgsGrp2cmd2 ViperGrp2cmd2

// grp2cmd2Cmd represents the grp2cmd2 command
var grp2cmd2Cmd = &cobra.Command{
	Use:     "grp2cmd2",
	Short:   "Test Nested Command",
	GroupID: "group2",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("flags from subcommand grp2cmd2")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2flag1: %s", vprFlgsGrp2cmd2.Grp2cmd2Flag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2flag2: %s", vprFlgsGrp2cmd2.Grp2cmd2Flag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2flag3: %s", vprFlgsGrp2cmd2.Grp2cmd2Flag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2flag4: %s", vprFlgsGrp2cmd2.Grp2cmd2Flag4)

		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("persistent flags from subcommand grp2cmd2")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2persistentflag1: %s", vprFlgsGrp2cmd2.Grp2cmd2PersistentFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2persistentflag2: %s", vprFlgsGrp2cmd2.Grp2cmd2PersistentFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2persistentflag3: %s", vprFlgsGrp2cmd2.Grp2cmd2PersistentFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("grp2cmd2persistentflag4: %s", vprFlgsGrp2cmd2.Grp2cmd2PersistentFlag4)

		fmt.Println("")

		logrus.WithField("cobra-cmd", cmd.Use).Infof("Persistent flags from rootCmd")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag1: %s", vprFlgsRoot.RootPersistentFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag2: %s", vprFlgsRoot.RootPersistentFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag3: %s", vprFlgsRoot.RootPersistentFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag4: %s", vprFlgsRoot.RootPersistentFlag4)
	},
}

func init() {
	cobra.OnInitialize(func() {
		InitViperSubCmd(viper.GetViper(), grp2cmd2Cmd, &vprFlgsGrp2cmd2)
	})
	rootCmd.AddCommand(grp2cmd2Cmd)

	grp2cmd2Cmd.PersistentFlags().StringVar(&grp2cmd2PersistentFlag1, "grp2cmd2persistentflag1", "value from default", "grp2cmd2 Persistent flag 1")
	grp2cmd2Cmd.PersistentFlags().StringVar(&grp2cmd2PersistentFlag2, "grp2cmd2persistentflag2", "value from default", "grp2cmd2 Persistent flag 2")
	grp2cmd2Cmd.PersistentFlags().StringVar(&grp2cmd2PersistentFlag3, "grp2cmd2persistentflag3", "value from default", "grp2cmd2 Persistent flag 3")
	grp2cmd2Cmd.PersistentFlags().StringVar(&grp2cmd2PersistentFlag4, "grp2cmd2persistentflag4", "value from default", "grp2cmd2 Persistent flag 4")

	grp2cmd2Cmd.Flags().StringVar(&grp2cmd2Flag1, "grp2cmd2flag1", "value from default", "grp2cmd2 flag 1")
	grp2cmd2Cmd.Flags().StringVar(&grp2cmd2Flag2, "grp2cmd2flag2", "value from default", "grp2cmd2 flag 2")
	grp2cmd2Cmd.Flags().StringVar(&grp2cmd2Flag3, "grp2cmd2flag3", "value from default", "grp2cmd2 flag 3")
	grp2cmd2Cmd.Flags().StringVar(&grp2cmd2Flag4, "grp2cmd2flag4", "value from default", "grp2cmd2 flag 4")
}
