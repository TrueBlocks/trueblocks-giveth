package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/bykof/gostradamus"
	"github.com/spf13/cobra"
)

func RunCombine(cmd *cobra.Command, args []string) error {
	if err := allDonations(); err != nil {
		return err
	}
	return nil
}

func allDonations() error {
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
		_, fn, round, _, _ := fnParts(f)
		if strings.Contains(f, "Round00") {
			continue
		}
		lines := file.AsciiFileToLines(f)
		for _, line := range lines {
			if !strings.Contains(line, "\"amount\"") {
				fields := strings.Split(line, ",")
				line = strings.Join(fields[0:11], ",")
				l := fmt.Sprintf("\"%s\",\"Round%02d\",%s\n", typeMap[fn], round, line)
				out.Write([]byte(l))
			}
		}
	}

	return nil
}

func fnParts(path string) (dir, fn string, round int64, sd, ed gostradamus.DateTime) {
	dir, fn = filepath.Split(path)
	fn = strings.Replace(strings.Replace(fn, "-202", "|202", -1), "-Round", "|Round", -1)
	parts := strings.Split(fn, "|")
	fn = parts[0]
	x := strings.Replace(strings.Replace(parts[3], "Round", "", -1), ".csv", "", -1)
	round, _ = strconv.ParseInt(x, 10, 64)
	sd = ParseDate(parts[1])
	ed = ParseDate(parts[2])
	return
}

func ParseDate(s string) (ret gostradamus.DateTime) {
	s = strings.Replace(s, ".000000Z", "", -1)
	s = strings.Replace(strings.Replace(strings.Replace(s, "T", "-", -1), ":", "-", -1), "_", "-", -1)
	parts := strings.Split(s, "-")
	ret = gostradamus.NewUTCDateTime(
		getInt(parts[0]),
		getInt(parts[1]),
		getInt(parts[2]),
		getInt(parts[3]),
		getInt(parts[4]),
		getInt(parts[5]),
		0,
	)
	return
}

func getInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
