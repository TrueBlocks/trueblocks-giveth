package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Various routines to download and manipulate the data",
	Long: `Downloads to the local hard drive one of the five data
structures (purple-list | not-eligible | eligible | calc-givback |
purple-verified) which are [documented here](./data/QUESTIONS.md)`,
	RunE: internal.RunData,
}

func init() {
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")
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
