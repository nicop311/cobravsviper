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

var zuLuSub221Flag1 string
var zuLuSub221Flag2 string
var zuLuSub221Flag3 string
var zuLuSub221Flag4 string

type ViperZuLuSub221 struct {
	ZuLuSub221Flag1 string `mapstructure:"zu-lu-sub221flag1"`
	ZuLuSub221Flag2 string `mapstructure:"zu-lu-sub221flag2"`
	ZuLuSub221Flag3 string `mapstructure:"zu-lu-sub221flag3"`
	ZuLuSub221Flag4 string `mapstructure:"zu-lu-sub221flag4"`
}

var vprFlgsZuLuSub221 ViperZuLuSub221

// zuLuSub221Cmd represents the zuLuSub221 command
var zuLuSub221Cmd = &cobra.Command{
	Use:   "zu-lu-sub221",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Manually call parent’s PersistentPreRunE
		if cmd.Parent() != nil && cmd.Parent().PersistentPreRunE != nil {
			if err := cmd.Parent().PersistentPreRunE(cmd.Parent(), args); err != nil {
				return err
			}
		}

		InitViperSubCmdE(viper.GetViper(), cmd, &vprFlgsZuLuSub221)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("flags from subcommand sub221")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("zu-lu-sub221flag1: %s", vprFlgsZuLuSub221.ZuLuSub221Flag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("zu-lu-sub221flag2: %s", vprFlgsZuLuSub221.ZuLuSub221Flag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("zu-lu-sub221flag3: %s", vprFlgsZuLuSub221.ZuLuSub221Flag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("zu-lu-sub221flag4: %s", vprFlgsZuLuSub221.ZuLuSub221Flag4)

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
	grp2cmd2Cmd.AddCommand(zuLuSub221Cmd)

	zuLuSub221Cmd.Flags().StringVar(&zuLuSub221Flag1, "zu-lu-sub221flag1", "value from default", "zu-lu-sub221 flag 1")
	zuLuSub221Cmd.Flags().StringVar(&zuLuSub221Flag2, "zu-lu-sub221flag2", "value from default", "zu-lu-sub221 flag 2")
	zuLuSub221Cmd.Flags().StringVar(&zuLuSub221Flag3, "zu-lu-sub221flag3", "value from default", "zu-lu-sub221 flag 3")
	zuLuSub221Cmd.Flags().StringVar(&zuLuSub221Flag4, "zu-lu-sub221flag4", "value from default", "zu-lu-sub221 flag 4")
}
