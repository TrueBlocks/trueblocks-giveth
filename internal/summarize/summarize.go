package summarize

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func RunSummarize(cmd *cobra.Command, args []string) error {
	if err := combineDonations(); err != nil {
		return err
	}
	if err := combineGivbacks(); err != nil {
		return err
	}
	if err := summarizeGivbacks(cmd, args); err != nil {
		return err
	}
	return nil
}

func getSummarizeOptions(cmd *cobra.Command, args []string) (donations []data.Donation, globals internal.Globals, err error) {
	globals, err = internal.GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}
	donations, err = data.NewDonations("./data/summaries/all_donations.csv", "csv")
	return
}
