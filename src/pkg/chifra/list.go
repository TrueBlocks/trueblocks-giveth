package chifra

import (
	"fmt"
	"os"
)

func ChifraList(w *os.File, fields map[string]string) (*SimpleListCount, error) {
	args := []string{
		"list",
		"--no_header",
		fields["count"],
		"--fmt",
		"json",
		"--chain",
		fields["chain"],
		fields["address"],
	}

	return commandToRecord[SimpleListCount](w, args)
}

func ChifraListCount(w *os.File, chain, address string) (uint64, error) {
	callParams := map[string]string{
		"chain":   chain,
		"address": address,
		"count":   "--count",
	}

	result, err := ChifraList(w, callParams)
	if err != nil {
		return 0, err
	}

	return result.NRecords, nil
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
