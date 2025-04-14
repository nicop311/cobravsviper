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

var sub221Flag1 string
var sub221Flag2 string
var sub221Flag3 string
var sub221Flag4 string

type ViperSub221 struct {
	Sub221Flag1 string `mapstructure:"sub221flag1"`
	Sub221Flag2 string `mapstructure:"sub221flag2"`
	Sub221Flag3 string `mapstructure:"sub221flag3"`
	Sub221Flag4 string `mapstructure:"sub221flag4"`
}

var vprFlgsSub221 ViperSub221

// sub221Cmd represents the sub221 command
var sub221Cmd = &cobra.Command{
	Use:   "sub221",
	Short: "Test Nested Command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("flags from subcommand sub221")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("sub221flag1: %s", vprFlgsSub221.Sub221Flag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("sub221flag2: %s", vprFlgsSub221.Sub221Flag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("sub221flag3: %s", vprFlgsSub221.Sub221Flag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("sub221flag4: %s", vprFlgsSub221.Sub221Flag4)

		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("flags from subcommand grp2cmd2")
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
		InitViperSubCmd(viper.GetViper(), sub221Cmd, &vprFlgsSub221)
	})

	grp2cmd2Cmd.AddCommand(sub221Cmd)

	sub221Cmd.Flags().StringVar(&sub221Flag1, "sub221flag1", "value from default", "sub221 flag 1")
	sub221Cmd.Flags().StringVar(&sub221Flag2, "sub221flag2", "value from default", "sub221 flag 2")
	sub221Cmd.Flags().StringVar(&sub221Flag3, "sub221flag3", "value from default", "sub221 flag 3")
	sub221Cmd.Flags().StringVar(&sub221Flag4, "sub221flag4", "value from default", "sub221 flag 4")
}
