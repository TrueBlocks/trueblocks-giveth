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
	cmdArgs := []string{}
	var cmdString = []string{
		"transactions",
		"--no_header",
		"--articulate",
		"--fmt",
		"json",
		"--chain",
		"[{CHAIN}]",
		"[{HASH}]",
	}

	for _, f := range cmdString {
		f = strings.Replace(f, "[{CHAIN}]", fields["chain"], -1)
		f = strings.Replace(f, "[{HASH}]", fields["hash"], -1)
		cmdArgs = append(cmdArgs, f)
	}

	if result, err := commandToRecord[SimpleTransfer](w, cmdArgs); err != nil {
		return nil, err
	} else {
		if len(result.Input) >= 10 && len(result.Encoding) == 0 {
			result.Encoding = result.Input[:10]
		} else if len(result.Input) < 10 {
			result.Recipient = result.Token
			result.Token = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
			// result.Amount = fmt.Sprintf("%f", result.Ether)
			result.Amount = fmt.Sprintf("%d", result.Value.Uint64())
		}
		if result.Encoding == "0xa9059cbb" {
			// compressedTx format is {name:transfer|inputs:{to:0xb2645970941b45f508b5333c1d628ad619adde20|value:200000000000000000000}|outputs:{_success:}}
			removes := []string{
				"{",
				"}",
				"inputs:",
				"outputs:",
			}
			input := result.CompressedTx
			for _, r := range removes {
				input = strings.Replace(input, r, "", -1)
			}
			parts := strings.Split(input, "|")
			for i := 0; i < len(parts); i++ {
				parts[i] = strings.Split(parts[i], ":")[1]
			}
			result.Recipient = parts[1]
			result.Amount = parts[2]
		}

		fmt.Fprintln(w, result.String())
		return &result, nil
	}
}

func (tx *SimpleTransfer) String() string {
	d := colors.BrightBlue + strings.Split(tx.Date, ":")[0] + ":" + strings.Split(tx.Date, ":")[1] + colors.Off
	s, _ := AddressToName(tx.Sender, false)
	t, _ := AddressToName(tx.Token, false)
	r, _ := AddressToName(tx.Recipient, false)
	a := "[" + colors.Yellow + utils.PadLeft(tx.Amount, 25, '.') + colors.Off + "]"
	return fmt.Sprintf(" %s %s %s %s of %s to %s", d, s.Name, colors.BrightWhite+"donated"+colors.Off, a, t.Name, r.Name)
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
	Recipient        string        `json:"recipient"`
	Amount           string        `json:"amount"`
	Input            string        `json:"input"`
	Encoding         string        `json:"encoding"`
	ArticulatedTx    ArticulatedTx `json:"articulatedTx"`
	CompressedTx     string        `json:"compressedTx"`
}

type ArticulatedTx struct {
	Name    string            `json:"name"`
	Inputs  map[string]string `json:"inputs"`
	Outputs map[string]string `json:"outputs"`
}
