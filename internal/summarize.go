package internal

import (
	"log"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
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

func summarizeGivbacks(cmd *cobra.Command, args []string) (err error) {
	err = sumByAddress(cmd, args)
	if err != nil {
		return err
	}

	return sumByAddressByRound(cmd, args)
}

func sumByAddress(cmd *cobra.Command, args []string) (err error) {
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

func sumByAddressByRound(cmd *cobra.Command, args []string) (err error) {
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

func combineGivbacks() error {
	files := data.GetFilesInFolders([]string{"calculate-givback"})

	out, err := os.Create("./data/summaries/all_givbacks.csv")
	if err != nil {
		return err
	}
	defer out.Close()

	out.Write([]byte("\"type\",\"round\",\"givDistributed\",\"givFactor\",\"givPrice\",\"givbackUsdValue\",\"giverAddress\",\"giverEmail\",\"giverName\",\"totalDonationsUsdValue\",\"givback\",\"share\"\n"))
	for _, f := range files {
		if strings.Contains(f, "Round00") {
			continue
		}
		lines := file.AsciiFileToLines(f)
		for _, line := range lines {
			if !strings.Contains(line, "\"givDistributed\"") {
				out.Write([]byte(line + "\n"))
			}
		}
	}

	return nil
}

func combineDonations() error {
	folders := []string{
		"purpleList-donations-to-verifiedProjects",
		"eligible-donations",
		"not-eligible-donations",
	}

	files := data.GetFilesInFolders(folders)

	out, err := os.Create("./data/summaries/all_donations.csv")
	if err != nil {
		return err
	}
	defer out.Close()

	out.Write([]byte("\"type\",\"round\",\"amount\",\"currency\",\"createdAt\",\"valueUsd\",\"giverAddress\",\"txHash\",\"network\",\"source\",\"giverName\",\"giverEmail\",\"projectLink\"\n"))
	for _, f := range files {
		if strings.Contains(f, "Round00") {
			continue
		}
		lines := file.AsciiFileToLines(f)
		for _, line := range lines {
			if !strings.Contains(line, "\"amount\"") {
				out.Write([]byte(line + "\n"))
			}
		}
	}

	return nil
}

func getSummarizeOptions(cmd *cobra.Command, args []string) (donations []data.Donation, globals Globals, err error) {
	globals, err = getGlobals("csv", cmd, args)
	if err != nil {
		return
	}
	donations, err = data.NewDonations("./data/summaries/all_donations.csv", "csv")
	return
}
