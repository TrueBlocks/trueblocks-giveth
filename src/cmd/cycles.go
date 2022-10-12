package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// cyclesCmd represents the data command
var cyclesCmd = &cobra.Command{
	Use:   "cycles",
	Short: "Tools to find cycles in donor's transaction history",
	Long: `Various tools to extract data needed to find cycles in a donor's
transaction history which may indicate a voilation of the rules.`,
	RunE: internal.RunCycles,
}

func init() {
	opts := "[txs|neighbors]"
	cyclesCmd.Flags().StringP("data", "d", "", "One of "+opts)
	cyclesCmd.Flags().SortFlags = false
	rootCmd.AddCommand(cyclesCmd)
}
