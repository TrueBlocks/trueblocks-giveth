package data

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type Donation struct {
	Amount       uint64  `json:"amount"`
	Currency     string  `json:"currency"`
	CreatedAt    string  `json:"createdAt"`
	ValueUsd     float64 `json:"valueUsd"`
	GiverAddress string  `json:"giverAddress"`
	TxHash       string  `json:"txHash"`
	Network      string  `json:"xDAI"`
	Source       string  `json:"source"`
	GiverName    string  `json:"giverName"`
	GiverEmail   string  `json:"giverEmail"`
	ProjectLink  string  `json:"projectlink"`
}

func (d Donation) String() string {
	b, _ := json.MarshalIndent(d, "", "  ")
	return string(b)
}

func NewDonations(path string, format string) (donations []Donation, err error) {
	header := make(map[string]int)
	lines := file.AsciiFileToLines(path)
	for _, line := range lines {
		record := strings.Split(strings.Replace(line, "\"", "", -1), ",")
		if len(header) == 0 {
			for i, r := range record {
				header[r] = i
			}

		} else {
			amount, _ := strconv.ParseUint(record[header["amount"]], 10, 64)
			valueUsd, _ := strconv.ParseFloat(record[header["valueUsd"]], 64)
			donations = append(donations, Donation{
				Amount:       amount,
				Currency:     record[header["currency"]],
				CreatedAt:    record[header["createdAt"]],
				ValueUsd:     valueUsd,
				GiverAddress: strings.ToLower(record[header["giverAddress"]]),
				TxHash:       record[header["txHash"]],
				Network:      record[header["network"]],
				Source:       record[header["source"]],
				GiverName:    record[header["giverName"]],
				GiverEmail:   record[header["giverEmail"]],
				ProjectLink:  record[header["projectLink"]],
			})
		}
	}

	return
}
