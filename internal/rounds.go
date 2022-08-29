package internal

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
	"github.com/spf13/cobra"
)

// RunRounds runs the rounds command
func RunRounds(cmd *cobra.Command, args []string) error {
	globals, err := getRoundsOptions(cmd, args)
	if err != nil {
		return err
	}

	obj := data.Round{}
	output.Header(obj, os.Stdout, globals.Format)
	defer output.Footer(obj, os.Stdout, globals.Format)

	for i, round := range globals.Rounds {
		if globals.Format == "" {
			ff := "Round %d %s --> %s %d %g %d %d\n"
			dateFmt := "YYYY-MM-DDTHH:mm:ss"
			ts := uint64(round.EndDate.UnixTimestamp())
			bnGnosis, _ := tslib.FromTsToBn("gnosis", ts)
			bnMainnet, _ := tslib.FromTsToBn("mainnet", ts)
			fmt.Printf(ff, round.Id, round.StartDate.Format(dateFmt), round.EndDate.Format(dateFmt), round.Available, round.Price, bnGnosis, bnMainnet)
		} else {
			output.Line(round, os.Stdout, globals.Format, i == 0)
		}
	}
	return nil
}

// getRoundsOptions processes command line options for the Rounds command
func getRoundsOptions(cmd *cobra.Command, args []string) (globals Globals, err error) {
	return GetGlobals("", cmd, args)
}
