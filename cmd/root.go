/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debug     bool
	logFormat string
	logLevel  string
)

var rootFlag1 string
var rootFlag2 string
var rootFlag3 string
var rootFlag4 string

var rootPersistentFlag1 string
var rootPersistentFlag2 string
var rootPersistentFlag3 string
var rootPersistentFlag4 string
var cfgFile string

type ViperFlagsRoot struct {
	CfgFile             string
	RootFlag1           string `mapstructure:"rootflag1"`
	RootFlag2           string `mapstructure:"rootflag2"`
	RootFlag3           string `mapstructure:"rootflag3"`
	RootFlag4           string `mapstructure:"rootflag4"`
	RootPersistentFlag1 string `mapstructure:"rootpersistentflag1"`
	RootPersistentFlag2 string `mapstructure:"rootpersistentflag2"`
	RootPersistentFlag3 string `mapstructure:"rootpersistentflag3"`
	RootPersistentFlag4 string `mapstructure:"rootpersistentflag4"`

	Debug     bool   `mapstructure:"debug"`
	LogFormat string `mapstructure:"log-format"`
	LogLevel  string `mapstructure:"log-level"`
}

// Initialize the ViperConfig struct with all the root CLI flags bound to Viper env vars
var vprFlgsRoot ViperFlagsRoot

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobravsviper",
	Short: "The ROOT command",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set logs format
		switch vprFlgsRoot.LogFormat {
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{
				PrettyPrint: false,
			})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{
				ForceColors:      true,
				DisableTimestamp: true,
			})
		default:
			return errors.New("logrus unknown output format")
		}
		logrus.Debugf("logrus output format is set to: %s", vprFlgsRoot.LogFormat)

		// Initialize logrus log level and log format for all cobra commands and subcommands.
		debugFlagIsUsed := cmd.Flags().Lookup("debug").Changed

		switch {
		case debugFlagIsUsed:
			// harcode that the --debug flags set logrus level to debug
			logrus.SetLevel(logrus.DebugLevel)
		default:
			// get the log level from viper which is bind to the cobra flag --log-level
			level, err := logrus.ParseLevel(vprFlgsRoot.LogLevel)
			if err != nil {
				return err
			}
			logrus.SetLevel(level)
		}
		logrus.Debugf("logrus log-level is set to: %s", logrus.GetLevel())

		// PersistentPreRunE returns an error or nil
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		logrus.WithField("cobra-cmd", cmd.Use).Infof("Root command called")

		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag1: %s", vprFlgsRoot.RootFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag2: %s", vprFlgsRoot.RootFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag3: %s", vprFlgsRoot.RootFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootflag4: %s", vprFlgsRoot.RootFlag4)

		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag1: %s", vprFlgsRoot.RootPersistentFlag1)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag2: %s", vprFlgsRoot.RootPersistentFlag2)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag3: %s", vprFlgsRoot.RootPersistentFlag3)
		logrus.WithField("cobra-cmd", cmd.Use).Infof("rootpersistentflag4: %s", vprFlgsRoot.RootPersistentFlag4)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Ensure initConfig runs before anything else
	cobra.OnInitialize(initConfig)
	// Define command groups
	group1 := &cobra.Group{
		ID:    "group1",
		Title: "Group 1 Commands:",
	}
	group2 := &cobra.Group{
		ID:    "group2",
		Title: "Group 2 Commands:",
	}

	// Add groups to the root command
	rootCmd.AddGroup(group1)
	rootCmd.AddGroup(group2)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration File")

	// logging level
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Set logrus.SetLevel to \"debug\". This is equivalent to using --log-level=debug. Flags --log-level and --debug flag are mutually exclusive. Corresponding environment variable: K8S_KMS_PLUGIN_DEBUG.")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set logrus.SetLevel. Possible values: trace, debug, info, warning, error, fatal and panic. Flags --log-level and --debug flag are mutually exclusive. Corresponding environment variable: K8S_KMS_PLUGIN_LOG_LEVEL.")
	rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "Logrus log output format. Possible values: text, json. Corresponding environment variable: K8S_KMS_PLUGIN_LOG_FORMAT")
	rootCmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.MarkFlagsMutuallyExclusive("log-level", "debug")

	rootCmd.PersistentFlags().StringVar(&rootPersistentFlag1, "rootpersistentflag1", "value from default", "persistent root flag 1")
	rootCmd.PersistentFlags().StringVar(&rootPersistentFlag2, "rootpersistentflag2", "value from default", "persistent root flag 2")
	rootCmd.PersistentFlags().StringVar(&rootPersistentFlag3, "rootpersistentflag3", "value from default", "persistent root flag 3")
	rootCmd.PersistentFlags().StringVar(&rootPersistentFlag4, "rootpersistentflag4", "value from default", "persistent root flag 4")

	rootCmd.Flags().StringVar(&rootFlag1, "rootflag1", "value from default", "root flag 1")
	rootCmd.Flags().StringVar(&rootFlag2, "rootflag2", "value from default", "root flag 2")
	rootCmd.Flags().StringVar(&rootFlag3, "rootflag3", "value from default", "root flag 3")
	rootCmd.Flags().StringVar(&rootFlag4, "rootflag4", "value from default", "root flag 4")
}

func initConfig() {
	// use a configuration file parsed by viper
	if rootCmd.Flags().Lookup("config").Changed && cfgFile != "" {
		logrus.Tracef("Case config file from the flag: %s", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else if envVar, ok := os.LookupEnv("COBRAVSVIPER_CONFIG"); ok {
		logrus.Tracef("Case config file from the environment variable: %s", envVar)
		viper.SetConfigFile(envVar)
	} else {
		logrus.Infof("Case config file from default location")
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("cobravsviper.conf") // name of config file (viper needs no file extension)
		// TODO: consider using the rootCmd.Flags().Lookup("config").DefValue for viper.SetConfigName
		// logrus.Infof("default config filename %s", rootCmd.Flags().Lookup("config").DefValue)
		logrus.Infof("Search config in .config directory %s with name cobravsviper.conf.yaml (without extension).", home)
		viper.AddConfigPath(home)
		logrus.Infof("Search config in home directory %s with name cobravsviper.conf.yaml (without extension).", filepath.Join(home, ".config/cobravsviper"))
		viper.AddConfigPath(filepath.Join(home, ".config/cobravsviper"))
	}

	// If a config file is not found, log a trace error. Otherwise, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Trace("No config file found; continue with cobra default values")
		} else {
			// Config file was found but another error occurred
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	InitViperSubCmd(viper.GetViper(), rootCmd, &vprFlgsRoot)

	// Debugging: Show all loaded settings
	logrus.Tracef("Viper settings: %+v", viper.AllSettings())

}
