package chifra

import (
	"fmt"
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
		if len(result.Input) > 10 && len(result.Encoding) == 0 {
			result.Encoding = result.Input[:10]
		} else {
			result.Recipient = result.Token
			result.Token = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
			result.Amount = fmt.Sprintf("%f", result.Ether)
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

func TraceSourceForTx(w *os.File, depth, max_depth int, hash, chain string) error {
	callParams := map[string]string{"chain": chain, "hash": hash}
	if tx, err := ChifraTransactions(w, callParams); err != nil {
		return err
	} else {
		SourceOfFunds(w, *tx, chain, depth, nil, postExportFunc)
	}
	return nil
}

/*
{
  "data": [
    {
      "hash": "0x6df1c8ccaf5f29c955c5359a575195a62daa8444352b4ff9fd01dab4c8c702ad",
      "blockHash": "0xa4a3e66273fb001b23c187db25e9185f3396a0121ac7519a21b7b8f6fef63f70",
      "blockNumber": 20552778,
      "transactionIndex": 1,
      "timestamp": 1644427220,
      "from": "0xeb2865c3324c0839ef657fc080128fcf440b9a91",
      "to": "0x4f4f9b8d5b4d0dc10506e5551b0513b61fd59e75",
      "value": 0,
      "gas": 45146,
      "gasPrice": 1317860964,
      "maxFeePerGas": 1317860964,
      "maxPriorityFeePerGas": 1317860964,
      "input": "0xa9059cbb0000000000000000000000004d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd000000000000000000000000000000000000000000000000002386f26fc10000",
      "isError": 0,
      "hasToken": 0,
      "receipt": {
        "contractAddress": "0x0",
        "gasUsed": 42655,
        "effectiveGasPrice": 1317860964,
        "logs": [
          {
            "address": "0x4f4f9b8d5b4d0dc10506e5551b0513b61fd59e75",
            "logIndex": 2,
            "topics": [
              "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
              "0x000000000000000000000000eb2865c3324c0839ef657fc080128fcf440b9a91",
              "0x0000000000000000000000004d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd"
            ],
            "data": "0x000000000000000000000000000000000000000000000000002386f26fc10000",
            "articulatedLog": {
              "name": "Transfer",
              "inputs": {
                "_amount": "10000000000000000",
                "_from": "0xeb2865c3324c0839ef657fc080128fcf440b9a91",
                "_to": "0x4d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd"
              }
            },
            "compressedLog": "{name:Transfer|inputs:{_amount:10000000000000000|_from:0xeb2865c3324c0839ef657fc080128fcf440b9a91|_to:0x4d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd}}"
          }
        ],
        "status": 1
      },
      "articulatedTx": {
        "name": "transfer",
        "stateMutability": "nonpayable",
        "inputs": {
          "_to": "0x4d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd",
          "_value": "10000000000000000"
        },
        "outputs": {
          "_success": ""
        }
      },
      "compressedTx": "{name:transfer|inputs:{_to:0x4d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd|_value:10000000000000000}|outputs:{_success:}}",
      "gasCost": 56213359419420,
      "gasUsed": 42655,
      "date": "2022-02-09 17:20:20 UTC",
      "ether": 0
    }
  ]
}
*/
