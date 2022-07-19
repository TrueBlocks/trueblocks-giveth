package internal

import (
	"fmt"
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

	obj := toRoundInterface(data.Round{}, globals.Format)
	output.Header(obj, os.Stdout, globals.Format)
	defer output.Footer(obj, os.Stdout, globals.Format)

	for i, round := range globals.Rounds {
		if globals.Format == "" {
			ff := "Round %d %s --> %s %d %g\n"
			dateFmt := "YYYY-MM-DDTHH:mm:ss"
			fmt.Printf(ff, round.Id, round.StartDate.Format(dateFmt), round.EndDate.Format(dateFmt), round.Available, round.Price)
		} else {
			obj := toRoundInterface(round, globals.Format)
			output.Line(obj, os.Stdout, globals.Format, i == 0)
		}
	}
	return nil
}

// getRoundsOptions processes command line options for the Rounds command
func getRoundsOptions(cmd *cobra.Command, args []string) (globals Globals, err error) {
	return getGlobals("", cmd, args)
}

func toRoundInterface(round data.Round, format string) interface{} {
	return interface{}(round)
}
