package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func RunSourceOfFunds(cmd *cobra.Command, args []string) error {
	hash, levels, max_rows, globals, err := getSofOptions(cmd, args)
	if err != nil {
		return err
	}

	w := os.Stdout

	if len(hash) > 0 {
		Header(w, 0, levels, hash, globals.Chain)
		fmt.Println(hash)

	} else {
		for _, round := range globals.Rounds {
			donations, _ := data.NewDonations(data.GetFilename("eligible", "csv", round), "csv", data.SortByHash)
			for _, donation := range donations {
				fmt.Println(donation)
			}
		}
	}

	return nil
}

func Header(w *os.File, i, l int, hash, chain string) {
	fmt.Fprintln(w, strings.Repeat(" ", 120), "\n", colors.BrightBlack+strings.Repeat("-", 5), fmt.Sprintf("%d-%d.", i, l), hash, chain, strings.Repeat("-", 70), colors.Off)
}

// getSofOptions processes command line options for the Rounds command
func getSofOptions(cmd *cobra.Command, args []string) (string, int, int, Globals, error) {
	globals, err := GetGlobals("csv", cmd, args)
	if err != nil {
		return "", 0, 0, globals, err
	}

	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return hash, 0, 0, globals, err
	}

	levels, err := cmd.Flags().GetUint64("levels")
	if err != nil {
		return hash, 0, 0, globals, err
	}
	if levels == 0 {
		levels = 1
	}

	max_rows, _ := strconv.ParseUint(os.Getenv("MAX_ROWS"), 10, 64)
	if max_rows == 0 {
		max_rows = 4000000000
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		os.Setenv("SHOW_SKIPS", "true")
	}

	return hash, int(levels), int(max_rows), globals, nil
}
