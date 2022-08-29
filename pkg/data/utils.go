package data

import (
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/bykof/gostradamus"
)

func parseFloat(val string) float64 {
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

func parseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func ParseDate(s string) (ret gostradamus.DateTime) {
	s = strings.Replace(s, ".000000Z", "", -1)
	s = strings.Replace(strings.Replace(strings.Replace(s, "T", "-", -1), ":", "-", -1), "_", "-", -1)
	parts := strings.Split(s, "-")
	ret = gostradamus.NewUTCDateTime(
		parseInt(parts[0]),
		parseInt(parts[1]),
		parseInt(parts[2]),
		parseInt(parts[3]),
		parseInt(parts[4]),
		parseInt(parts[5]),
		0,
	)
	return
}

func GetFilesInFolders(folders []string) (files []string) {
	for _, folder := range folders {
		log.Println("Walking", folder)
		filepath.Walk(DataFolder()+folder, func(path string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(path, ".csv") {
				files = append(files, path)
			}
			return nil
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	return
}

// var cmdFolderMap = map[string]string{
// 	"purple-list":     "purpleList",
// 	"eligible":        "eligible-donations",
// 	"not-eligible":    "not-eligible-donations",
// 	"purple-verified": "purpleList-donations-to-verifiedProjects",
// 	"calc-givback":    "calculate-givback",
// }

var folderCmdMap = map[string]string{
	"purpleList":                               "purple-list",
	"eligible-donations":                       "eligible",
	"not-eligible-donations":                   "not-eligible",
	"purpleList-donations-to-verifiedProjects": "purple-verified",
	"calculate-givback":                        "calc-givback",
}

// func ExplodeFilename(path string) (dir, fn, typ string, round int64, sd, ed gostradamus.DateTime) {
func ExplodeFilename(path string) (fn, typ, fmt string, round int64) {
	// dir, fn = filepath.Split(path)
	_, fn = filepath.Split(path)
	fn = strings.Replace(strings.Replace(fn, "-202", "|202", -1), "-Round", "|Round", -1)
	parts := strings.Split(fn, "|")
	fn = parts[0]
	var xx []string
	if len(parts) > 3 {
		xx = strings.Split(parts[3], ".")
		fmt = xx[1]
	} else {
		xx = strings.Split(parts[0], ".")
		fmt = xx[1]
	}
	round, _ = strconv.ParseInt(strings.Replace(xx[0], "Round", "", -1), 10, 64)
	// sd = ParseDate(parts[1])
	// ed = ParseDate(parts[2])
	typ = folderCmdMap[fn]
	return
}

func DataFolder() string {
	if file.FolderExists("./data") {
		return "./data/"
	}
	return "/Users/jrush/Development/trueblocks-giveth/data/"
}
