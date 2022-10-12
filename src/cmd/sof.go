package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// sourceOfFundsCmd represents the data command
var sourceOfFundsCmd = &cobra.Command{
	Use:   "sof",
	Short: "Trace the source of funds for the transaction",
	Long: `Given a transaction, traces its source of funds all the
way back to the given start block on the given chain.`,
	RunE: internal.RunSourceOfFunds,
}

func init() {
	sourceOfFundsCmd.Flags().StringP("hash", "a", "", "trace source of funds for a single hash")
	sourceOfFundsCmd.Flags().Uint64P("levels", "l", 3, "the maximum number of levels deep to dig")
	rootCmd.AddCommand(sourceOfFundsCmd)
}
