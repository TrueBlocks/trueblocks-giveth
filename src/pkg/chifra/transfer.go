package chifra

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func (tr *SimpleTransfer) String() string {
	d := colors.BrightBlue + strings.Split(tr.Date, ":")[0] + ":" + strings.Split(tr.Date, ":")[1] + colors.Off
	s, _ := AddressToName(tr.Sender, false)
	t, _ := AddressToName(tr.Token, false)
	r, _ := AddressToName(tr.Recipient, false)
	a := "[" + colors.Yellow + utils.PadLeft(tr.Amount, 25, '.') + colors.Off + "]"
	return fmt.Sprintf(" %s %s %s %s of %s to %s", d, s.Name, colors.BrightWhite+"donated"+colors.Off, a, t.Name, r.Name)
}

type SimpleTransfer struct {
	// TransactionIndex uint64 `json:"transactionIndex"`
	// Timestamp        int64  `json:"timestamp"`
	// Hash             string `json:"hash"`
	BlockNumber uint64 `json:"blockNumber"`
	Date        string `json:"date"`
	Sender      string `json:"sender"`
	Token       string `json:"token"`
	Recipient   string `json:"recipient"`
	Amount      string `json:"amount"`
}
