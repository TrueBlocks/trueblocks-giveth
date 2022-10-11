package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/chifra"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func RunSourceOfFunds(cmd *cobra.Command, args []string) error {
	max_rows, max_depth, globals, err := getSofOptions(cmd, args)
	if err != nil {
		return err
	}

	for _, round := range globals.Rounds {
		donations, _ := data.NewDonations(data.GetFilename("eligible", "csv", round), "csv", data.SortByHash)
		for i, donation := range donations {
			if uint64(i) < max_rows && validate.IsValidHash(donation.TxHash) {
				w := os.Stdout
				// let them know we're here
				fmt.Fprintln(w, "\n", colors.BrightBlack+strings.Repeat("-", 5), fmt.Sprintf("%d-%d-%d.", round.Id, i, len(donations)), donation.TxHash, donation.Network, strings.Repeat("-", 70), colors.Off)
				if err := chifra.FindSourceOfFunds(w, 0, int(max_depth), donation.TxHash, donation.Network); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getSofOptions processes command line options for the Rounds command
func getSofOptions(cmd *cobra.Command, args []string) (max_rows, max_depth uint64, globals Globals, err error) {
	globals, err = GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}

	max_rows, _ = strconv.ParseUint(os.Getenv("MAX_ROWS"), 10, 64)
	if max_rows == 0 {
		max_rows = utils.NOPOS
	}

	max_depth, _ = strconv.ParseUint(os.Getenv("MAX_DEPTH"), 10, 64)
	if max_depth == 0 {
		max_depth = utils.NOPOS
	}

	return
}
