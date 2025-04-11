/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Support ENV variables with prefix with viper bound to cobra
	// Example: A CLI flag like --some-flag becomes COBRAVSVIPER_SOME_FLAG in environment variables.
	viper.SetEnvPrefix("COBRAVSVIPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // Converts flags to ENV format
	viper.AutomaticEnv()                                   // Enables automatic binding

	// Ensure all Cobra flags are bound to Viper after initializing them
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		logrus.Errorf("Error binding flags: %v", err)
		os.Exit(1)
	}
	vprBuf := viper.GetViper()
	logrus.Tracef("vprBuf := viper.GetViper(): %+v", vprBuf)

	// Initialize and Load the ViperConfig that are bound to root Cobra CLI flags
	if err := viper.Unmarshal(&vprFlgsRoot); err != nil {
		logrus.Fatalf("Failed to load viper config: %v", err)
	}

	// Debugging: Show all loaded settings
	logrus.Tracef("Viper settings: %+v", viper.AllSettings())

}
