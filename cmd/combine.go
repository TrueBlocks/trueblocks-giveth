package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// combineCmd represents the combine command
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "Combines data by type and time period (i.e. rounds)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: internal.RunCombine,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// combineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	combineCmd.Flags().SortFlags = false
	rootCmd.AddCommand(combineCmd)
}
