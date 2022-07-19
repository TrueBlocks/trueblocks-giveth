package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Givback struct {
	Type                   string  `json:"type"`
	Round                  string  `json:"round"`
	GivDistributed         float64 `json:"givDistributed"`
	GivFactor              float64 `json:"givFactor"`
	GivPrice               float64 `json:"givPrice"`
	GivbackUsdValue        float64 `json:"givbackUsdValue"`
	GiverAddress           string  `json:"giverAddress"`
	GiverEmail             string  `json:"giverEmail"`
	GiverName              string  `json:"giverName"`
	TotalDonationsUsdValue float64 `json:"totalDonationsUsdValue"`
	Givback                float64 `json:"givback"`
	Share                  float64 `json:"share"`
}

func (g Givback) String() string {
	b, _ := json.MarshalIndent(g, "", "  ")
	return string(b)
}

func NewGivback(path string, format string) (givbacks []Givback, err error) {
	_, typ, _, round := ExplodeFilename(path)

	header := make(map[string]int)
	lines := file.AsciiFileToLines(path)
	for _, line := range lines {
		record := strings.Split(strings.Replace(line, "\"", "", -1), ",")
		if len(header) == 0 {
			for i, r := range record {
				header[r] = i
			}

		} else {
			gb := Givback{
				Type:                   typ,
				Round:                  fmt.Sprintf("Round%02d", round),
				GivDistributed:         parseFloat(record[header["givDistributed"]]),
				GivFactor:              parseFloat(record[header["givFactor"]]),
				GivPrice:               parseFloat(record[header["givPrice"]]),
				GivbackUsdValue:        parseFloat(record[header["givbackUsdValue"]]),
				GiverAddress:           record[header["giverAddress"]],
				GiverEmail:             record[header["giverEmail"]],
				GiverName:              record[header["giverName"]],
				TotalDonationsUsdValue: parseFloat(record[header["totalDonationsUsdValue"]]),
				Givback:                parseFloat(record[header["givback"]]),
				Share:                  parseFloat(record[header["share"]]),
			}
			givbacks = append(givbacks, gb)
		}
	}
	return givbacks, nil
}
