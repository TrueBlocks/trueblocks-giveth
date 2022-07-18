package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/spf13/cobra"
)

func RunCombine(cmd *cobra.Command, args []string) error {
	if err := combineAll(); err != nil {
		return err
	}
	return nil
}

func combineAll() error {
	folders := []string{
		"purpleList-donations-to-verifiedProjects",
		"eligible-donations",
		"not-eligible-donations",
	}

	var typeMap = map[string]string{
		"purpleList-donations-to-verifiedProjects": "purple-verified",
		"eligible-donations":                       "eligible",
		"not-eligible-donations":                   "not-eligible",
	}

	// allLines := []string{}

	files := []string{}
	for _, folder := range folders {
		log.Println("Walking", folder)
		filepath.Walk("./data/"+folder, func(path string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(path, ".csv") {
				files = append(files, path)
			}
			return nil
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	out, err := os.Create("./data/combined/all_donations.csv")
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
				fields := strings.Split(line, ",")
				if len(fields) > 11 {
					line = ""
					for i := 0; i < 11; i++ {
						if i != 0 {
							line += ","
						}
						line += fields[i]
					}
				}
				dir, _ := filepath.Split(f)
				parts := strings.Split(strings.Trim(dir, "/"), "/")
				l := fmt.Sprintf("\"%s\",\"%s\",%s\n", typeMap[parts[len(parts)-1]], getRound(f), line)
				out.Write([]byte(l))
			}
		}
	}

	return nil
}

func getRound(path string) string {
	_, fn := filepath.Split(path)
	parts := strings.Split(fn, "-")
	return strings.Replace(parts[len(parts)-1], ".csv", "", -1)
}
