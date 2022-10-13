package chifra

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func ChifraTransactions(w *os.File, fields map[string]string) (*SimpleTransfer, error) {
	args := []string{
		"transactions",
		"--no_header",
		"--articulate",
		"--fmt",
		"json",
		"--chain",
		fields["chain"],
		fields["hash"],
	}

	if result, err := commandToRecord[SimpleTransaction](w, args); err != nil {
		return nil, err
	} else {
		return TransactionToTransfer(result)
	}
}

func (tx *SimpleTransaction) String() string {
	d := colors.BrightBlue + strings.Split(tx.Date, ":")[0] + ":" + strings.Split(tx.Date, ":")[1] + colors.Off
	s, _ := AddressToName(tx.Sender)
	t, _ := AddressToName(tx.Token)
	r, _ := AddressToName(tx.Recipient)
	a := "[" + colors.Yellow + utils.PadLeft(tx.Amount, 25, '.') + colors.Off + "]"
	return fmt.Sprintf(" %s %s %s %s of %s to %s", d, s.Name, colors.BrightWhite+"donated"+colors.Off, a, t.Name, r.Name)
}

type SimpleTransaction struct {
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

type ArticulatedTx struct {
	Name    string            `json:"name"`
	Inputs  map[string]string `json:"inputs"`
	Outputs map[string]string `json:"outputs"`
}

func TransactionToTransfer(tx *SimpleTransaction) (*SimpleTransfer, error) {
	var ret SimpleTransfer
	ret.Sender = tx.Sender
	ret.Token = tx.Token
	ret.Date = tx.Date
	ret.BlockNumber = tx.BlockNumber

	if len(tx.Input) >= 10 {
		tx.Encoding = tx.Input[:10]

	} else if len(tx.Input) < 10 {
		ret.Recipient = tx.Token
		ret.Token = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		ret.Amount = fmt.Sprintf("%d", tx.Value.Uint64())
	}

	if tx.Encoding == "0xa9059cbb" {
		removes := []string{
			"{",
			"}",
			"inputs:",
			"outputs:",
		}
		input := tx.CompressedTx
		for _, r := range removes {
			input = strings.Replace(input, r, "", -1)
		}
		parts := strings.Split(input, "|")
		for i := 0; i < len(parts); i++ {
			parts[i] = strings.Split(parts[i], ":")[1]
		}
		ret.Recipient = parts[1]
		ret.Amount = parts[2]
	}

	return &ret, nil
}
