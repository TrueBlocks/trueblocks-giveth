package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// roundsCmd represents the rounds command
var roundsCmd = &cobra.Command{
	Use:   "rounds",
	Short: "Print information about the rounds",
	Long: `The rounds tool prints information about the rounds. Rounds begin at 16:00
hour every other Thursday afternoon UTC time. The first round was on December 24, 2021.
A round's number can be used in other commands to as a shorthand to a date range.`,
	RunE: internal.RunRounds,
}

func init() {
	// roundsCmd.PersistentFlags().String("foo", "", "A help for foo")
	// roundsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(roundsCmd)
}
