package data

import (
	"encoding/json"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

type PurpleList struct {
	List []string `json:"purpleList"`
}

func (pl PurpleList) String() string {
	b, _ := json.MarshalIndent(pl.List, "", "  ")
	return string(b)
}

func NewPurpleList(path string) (pl PurpleList, err error) {
	data := file.AsciiFileToString(path)
	if err = json.Unmarshal([]byte(data), &pl); err != nil {
		return
	}
	sort.Slice(pl.List, func(i, j int) bool {
		return pl.List[i] < pl.List[j]
	})
	return
}
