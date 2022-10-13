package chifra

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func (tr *SimpleTransfer) String() string {
	d := colors.BrightBlue + strings.Split(tr.Date, ":")[0] + ":" + strings.Split(tr.Date, ":")[1] + colors.Off
	s, _ := AddressToName(tr.Sender)
	t, _ := AddressToName(tr.Token)
	r, _ := AddressToName(tr.Recipient)
	a := "[" + colors.Yellow + utils.PadLeft(tr.Amount, 25, '.') + colors.Off + "]"
	return fmt.Sprintf(" %s %s %s %s of %s to %s", d, s.Name, colors.BrightWhite+"[{DISPO}]"+colors.Off, a, t.Name, r.Name)
}

type SimpleTransfer struct {
	Hash             string        `json:"hash"`
	BlockNumber      uint64        `json:"blockNumber"`
	TransactionIndex uint64        `json:"transactionIndex"`
	Timestamp        int64         `json:"timestamp"`
	Date             string        `json:"date"`
	Ether            float64       `json:"ether"`
	Value            big.Int       `json:"value"`
	Sender           string        `json:"from"`
	Token            string        `json:"to"`
	Input            string        `json:"input"`
	Encoding         string        `json:"encoding"`
	ArticulatedTx    ArticulatedTx `json:"articulatedTx"`
	CompressedTx     string        `json:"compressedTx"`
	Recipient        string        `json:"recipient"`
	Amount           string        `json:"amount"`
}
