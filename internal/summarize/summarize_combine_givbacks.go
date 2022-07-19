package summarize

import (
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
)

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
