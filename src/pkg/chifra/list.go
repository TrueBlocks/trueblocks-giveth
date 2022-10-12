package chifra

import (
	"fmt"
	"os"
	"strings"
)

func ChifraList(w *os.File, fields map[string]string) (*SimpleListCount, error) {
	cmdArgs := []string{}
	var cmdString = []string{
		"list",
		"--no_header",
		"--count",
		"--fmt",
		"json",
		"--chain",
		"[{CHAIN}]",
		"[{ADDRESS}]",
	}

	silent := len(fields["--silent"]) > 0
	for _, f := range cmdString {
		f = strings.Replace(f, "[{CHAIN}]", fields["chain"], -1)
		f = strings.Replace(f, "[{HASH}]", fields["hash"], -1)
		cmdArgs = append(cmdArgs, f)
	}

	if result, err := commandToRecord[SimpleListCount](w, cmdArgs); err != nil {
		return nil, err
	} else {
		if !silent {
			fmt.Fprintln(w, result.String())
		}
		return &result, nil
	}
}

func (lc *SimpleListCount) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%d", lc.Address, lc.NRecords, lc.FileSize, lc.LastScanned)
}

type SimpleListCount struct {
	Address     string `json:"address"`
	NRecords    uint64 `json:"nRecords"`
	FileSize    uint64 `json:"fileSize"`
	LastScanned uint64 `json:"lastScanned"`
}
