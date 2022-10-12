package internal

import (
	"os"

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
		output.Line(round, os.Stdout, globals.Format, i == 0)
	}

	return nil
}

// getRoundsOptions processes command line options for the Rounds command
func getRoundsOptions(cmd *cobra.Command, args []string) (globals Globals, err error) {
	return GetGlobals("txt", cmd, args)
}
