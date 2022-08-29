package summarize

import (
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
)

func combineDonations() error {
	folders := []string{
		"purpleList-donations-to-verifiedProjects",
		"eligible-donations",
		"not-eligible-donations",
	}

	files := data.GetFilesInFolders(folders)

	out, err := os.Create(data.DataFolder() + "summaries/all_donations.csv")
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
