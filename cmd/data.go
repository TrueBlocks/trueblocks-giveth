package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Various routines to download and manipulate the data",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: internal.RunData,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	types := internal.DataTypes()
	opts := "["
	for i, t := range types {
		if i != 0 {
			opts += "|"
		}
		opts += " " + t + " "
	}
	opts += "]"
	dataCmd.Flags().StringP("data", "d", "", "One of "+opts)
	dataCmd.Flags().SortFlags = false

	rootCmd.AddCommand(dataCmd)
}
