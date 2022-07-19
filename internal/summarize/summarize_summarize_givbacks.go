package summarize

import (
	"log"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func summarizeGivbacks(cmd *cobra.Command, args []string) (err error) {
	err = summarizeByAddress(cmd, args)
	if err != nil {
		return err
	}

	return summarizeByAddressByRound(cmd, args)
}

func summarizeByAddress(cmd *cobra.Command, args []string) (err error) {
	log.Println("Counts by address...")
	donations, globals, err := getSummarizeOptions(cmd, args)
	if err != nil {
		return err
	}

	theMap := map[string]int{}
	for _, donation := range donations {
		theMap[donation.GiverAddress]++
	}

	results := []data.StringCounter{}
	for addr, cnt := range theMap {
		results = append(results, data.StringCounter{
			Key:   addr,
			Count: cnt,
		})
	}

	return data.WriteSummary("./data/summaries/donation_count_by_address.csv", results, data.Reverse, globals.Format)
}

func summarizeByAddressByRound(cmd *cobra.Command, args []string) (err error) {
	log.Println("Counts by round within address...")
	donations, globals, err := getSummarizeOptions(cmd, args)
	if err != nil {
		return err
	}

	theMap := map[string]map[string]int{}
	for _, donation := range donations {
		if theMap[donation.Round] == nil {
			theMap[donation.Round] = map[string]int{}
		}
		theMap[donation.Round][donation.GiverAddress]++
	}

	results := []data.TwoStringCounter{}
	for round, m := range theMap {
		for addr, cnt := range m {
			results = append(results, data.TwoStringCounter{
				Key1:  addr,
				Key2:  round,
				Count: cnt,
			})
		}
	}

	return data.WriteSummary("./data/summaries/donation_count_by_address_by_round.csv", results, data.Reverse, globals.Format)
}
