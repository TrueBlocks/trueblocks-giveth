package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "giveth",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.giveth.yaml)")
	rootCmd.PersistentFlags().MarkHidden("config")

	rootCmd.PersistentFlags().Uint64P("round", "r", 0, "Limits the list of rounds to a single round")
	rootCmd.PersistentFlags().BoolP("update", "u", false, "If present, data commands pull data from Giveth's APIs")
	rootCmd.PersistentFlags().BoolP("script", "c", false, "If present, data commands generate bash script to query Giveth's APIs")
	rootCmd.PersistentFlags().Uint64P("sleep", "s", 0, "Instructs the tool how long to sleep between invocations")
	rootCmd.PersistentFlags().StringP("fmt", "x", "", "One of [json|csv|txt]")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "If present, certain commands will display extra data")
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".giveth" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".giveth")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
