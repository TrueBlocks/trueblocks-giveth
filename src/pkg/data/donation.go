package data

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Donation struct {
	Type         string  `json:"type"`
	Round        string  `json:"round"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	CreatedAt    string  `json:"createdAt"`
	ValueUsd     float64 `json:"valueUsd"`
	GiverAddress string  `json:"giverAddress"`
	TxHash       string  `json:"txHash"`
	Network      string  `json:"network"`
	Source       string  `json:"source"`
	GiverName    string  `json:"giverName,omitempty"`
	GiverEmail   string  `json:"giverEmail,omitempty"`
	ProjectLink  string  `json:"projectlink"`
}

func (d Donation) String() string {
	b, _ := json.MarshalIndent(d, "", "  ")
	return string(b)
}

type sortField int

const (
	NoSort sortField = iota
	SortByHash
)

func NewDonations(path string, format string, sortBy sortField) (donations []Donation, err error) {
	header := make(map[string]int)
	lines := file.AsciiFileToLines(path)
	for _, line := range lines {
		record := strings.Split(strings.Replace(line, "\"", "", -1), ",")
		if len(header) == 0 {
			for i, r := range record {
				header[r] = i
			}

		} else {
			gem := ""
			if len(record) > header["giverEmail"] {
				gem = record[header["giverEmail"]]
			}
			pl := ""
			if len(record) > header["projectLink"] {
				pl = record[header["projectLink"]]
			}
			donation := Donation{
				Type:         record[header["type"]],
				Round:        record[header["round"]],
				Amount:       parseFloat(record[header["amount"]]),
				Currency:     record[header["currency"]],
				CreatedAt:    record[header["createdAt"]],
				ValueUsd:     parseFloat(record[header["valueUsd"]]),
				GiverAddress: strings.ToLower(record[header["giverAddress"]]),
				TxHash:       record[header["txHash"]],
				Network:      strings.Replace(record[header["network"]], "xDAI", "gnosis", -1),
				Source:       record[header["source"]],
				GiverName:    record[header["giverName"]],
				GiverEmail:   gem,
				ProjectLink:  pl,
			}
			donations = append(donations, donation)
		}
	}

	switch sortBy {
	case SortByHash:
		sort.Slice(donations, func(i, j int) bool {
			return donations[i].TxHash < donations[j].TxHash
		})
	case NoSort:
	default:
	}
	return
}
