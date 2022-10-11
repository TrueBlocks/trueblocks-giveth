package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// cyclesCmd represents the data command
var sourceOfFundsCmd = &cobra.Command{
	Use:   "sof",
	Short: "Trace the source of funds for the transaction",
	Long: `Given a transaction, traces its source of funds all the
way back to the given start block on the given chain.`,
	RunE: internal.RunSourceOfFunds,
}

func init() {
	rootCmd.AddCommand(sourceOfFundsCmd)
}
