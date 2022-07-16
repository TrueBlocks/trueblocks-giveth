package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
	"github.com/bykof/gostradamus"
	"github.com/spf13/cobra"
)

// RunRounds runs the rounds command
func RunRounds(cmd *cobra.Command, args []string) error {
	rounds, oneRound, format, err := getRoundsOptions(cmd, args)
	if err != nil {
		return err
	}

	unused := data.Round{}
	output.Header(unused, os.Stdout, format)
	defer output.Footer(unused, os.Stdout, format)

	for i, round := range rounds {
		single := uint64(round.Id) == oneRound
		timed := (oneRound == 0) && round.StartDate.Time().Before(gostradamus.Now().Time())
		if single || timed {
			if format == "" {
				ff := "Round %d %s --> %s %d %s\n"
				dateFmt := "YYYY-MM-DDTHH:mm:ss"
				fmt.Printf(ff, round.Id, round.StartDate.Format(dateFmt), round.EndDate.Format(dateFmt), round.Available, round.Price)
			} else {
				output.Line(round, os.Stdout, format, i == 0)
			}
		}
	}
	return nil
}

// getRoundsOptions processes command line options for the Rounds command
func getRoundsOptions(cmd *cobra.Command, args []string) (rounds []data.Round, oneRound uint64, format string, err error) {
	format, err = cmd.Flags().GetString("fmt")
	if err != nil {
		log.Fatal(err)
	}

	oneRound, err = cmd.Flags().GetUint64("round")
	if err != nil {
		log.Fatal(err)
	}

	return data.GetRounds(), oneRound, format, err
}
