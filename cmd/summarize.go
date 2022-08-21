package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal/summarize"
	"github.com/spf13/cobra"
)

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarizes and combines data by type and time period (i.e. rounds)",
	Long: `This tool begins the process of combining, summarizing and simplifying
data gathered from both the Giveth API and on-chain data.`,
	RunE: summarize.RunSummarize,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summarizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	summarizeCmd.Flags().SortFlags = false
	rootCmd.AddCommand(summarizeCmd)
}
