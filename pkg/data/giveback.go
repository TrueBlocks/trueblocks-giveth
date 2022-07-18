package data

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Givback struct {
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
				GivDistributed:         getFloat(record[header["givDistributed"]]),
				GivFactor:              getFloat(record[header["givFactor"]]),
				GivPrice:               getFloat(record[header["givPrice"]]),
				GivbackUsdValue:        getFloat(record[header["givbackUsdValue"]]),
				GiverAddress:           record[header["giverAddress"]],
				GiverEmail:             record[header["giverEmail"]],
				GiverName:              record[header["giverName"]],
				TotalDonationsUsdValue: getFloat(record[header["totalDonationsUsdValue"]]),
				Givback:                getFloat(record[header["givback"]]),
				Share:                  getFloat(record[header["share"]]),
			}
			givbacks = append(givbacks, gb)
		}
	}
	return givbacks, nil
}

func getFloat(val string) float64 {
	f, _ := strconv.ParseFloat(val, 64)
	return f
}
