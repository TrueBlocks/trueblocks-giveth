package internal

import (
	"fmt"
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
	if err := summarizeGivbacks(); err != nil {
		return err
	}
	return nil
}

func summarizeGivbacks() error {
	donations, _ := data.NewDonations("./data/summaries/all_donations.csv", "csv")

	addrMap := map[string]map[string]int{}
	for _, donation := range donations {
		if addrMap["1"] == nil {
			addrMap["1"] = map[string]int{}
		}
		addrMap["1"][donation.GiverAddress]++
	}

	for round, m := range addrMap {
		for addr, cnt := range m {
			fmt.Println(round, addr, cnt)
		}
	}

	return nil
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
