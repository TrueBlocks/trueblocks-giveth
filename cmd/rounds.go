package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// roundsCmd represents the rounds command
var roundsCmd = &cobra.Command{
	Use:   "rounds",
	Short: "Print information about the rounds",
	Long: `Rounds begin at 16:00 hour every other Thursday afternoon UTC time. The first
round was on December 24, 2021. A round's number can be used in other commands to
as a shorthand to a date range.`,
	RunE: internal.RunRounds,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// roundsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// roundsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(roundsCmd)
}
