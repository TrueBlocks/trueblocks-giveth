package chifra

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

type filterFunc func(string) bool
type postFunc func(*os.File, string, func(string) bool)

func commandToFields(w *os.File, args []string, filter filterFunc, post postFunc) {
	if ret, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", ("chifra " + strings.Join(args, " ")), "error:", err)
		os.Exit(1)
	} else {
		if post == nil {
			log.Fatal("You must provide a post processing function")
		}
		r := string(ret)
		s := strings.Split(r, "\t")
		if len(s) > 0 {
			post(w, r, filter)
		}
	}
}

func commandToRecord[T SimpleTransaction | SimpleListCount](w *os.File, args []string) (*T, error) {
	if bytes, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", ("chifra " + strings.Join(args, " ")), "error:", err)
		os.Exit(1)

	} else {
		resp := ChifraResponse[T]{}
		if err := json.Unmarshal(bytes, &resp); err != nil {
			return nil, err
		} else if len(resp.Data) == 0 {
			return nil, errors.New("no data returned from chifra " + strings.Join(args, " "))
		}
		return &resp.Data[0], nil
	}
	return nil, nil
}

// func commandToRecords[T SimpleTransaction | SimpleListCount | SimpleLog](w *os.File, args []string) ([]T, error) {
// 	if bytes, err := exec.Command("chifra", args...).Output(); err != nil {
// 		fmt.Fprintln(os.Stderr, "There was an error running the command: ", ("chifra " + strings.Join(args, " ")), "error:", err)
// 		os.Exit(1)
// 	} else {
// 		if isEmpty(bytes) {
// 			return []T{}, nil
// 		}
// 		fmt.Fprintln(os.Stderr, string(bytes))
// 		// bytes = append(bytes, 0x7d)
// 		// bytes = append(bytes, 0x5d)
// 		return []T{}, nil
// 		// os.Exit(1)
// 		// resp := ChifraResponse[T]{}
// 		// if err := json.Unmarshal(bytes, &resp); err != nil {
// 		// 	return nil, err
// 		// }
// 		// for _, i := range resp.Data {
// 		// 	fmt.Println(i)
// 		// 	fmt.Println()
// 		// }
// 		// return resp.Data, nil
// 	}
// 	return nil, nil
// }
// func isEmpty(bytes []byte) bool {
// 	return false
// }

// -----------
type ChifraResponse[T SimpleTransaction | SimpleListCount] struct {
	Data []T `json:"data"`
}

func TraceSourceOfFunds(w *os.File, depth, max_depth int, hash, chain string) error {
	if depth == max_depth {
		return nil
	}

	callParams := map[string]string{"chain": chain, "hash": hash}
	if tr, err := ChifraTransactions(w, callParams); err != nil {
		return err

	} else {
		fmt.Fprintln(w, tr.String())

		if nRecords, err := ChifraListCount(w, chain, tr.Sender); err != nil {
			return err

		} else {
			if nRecords > 20000 {
				if os.Getenv("SHOW_SKIPS") == "true" {
					fmt.Fprintln(os.Stderr, colors.Yellow, "Skipping address", tr.Sender, "too many records:", nRecords, colors.Off)
				}
				return nil
			}

			ChifraExport(w, tr, chain, depth, max_depth, postExportFunc)
			// hashes, err := GetAppearanceHashes(w, tr, chain)
			// if err != nil {
			// 	return err
			// }
			// for _, h := range hashes {
			// 	h = hash
			// 	if err := TraceSourceOfFunds(w, depth+1, max_depth, h, chain); err != nil {
			// 		return err
			// 	}
			// }
		}
	}

	return nil
}

// const transferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// func GetAppearanceHashes(w *os.File, tr *SimpleTransfer, chain string) ([]string, error) {
// 	var args = []string{
// 		"export",
// 		"--no_header",
// 		"--fmt",
// 		"json",
// 		tr.Sender,
// 		"--chain",
// 		chain,
// 		"--logs",
// 		"--emitter",
// 		tr.Token,
// 		"--reversed",
// 		"--first_block",
// 		firstBlocks[chain],
// 		"--last_block",
// 		fmt.Sprintf("%d", tr.BlockNumber),
// 		"--cache",
// 		"--articulate",
// 		transferTopic,
// 	}
//
// 	if logs, err := commandToRecords[SimpleLog](w, args); err != nil {
// 		return []string{}, err
// 	} else {
// 		hashes := []string{}
// 		for _, log := range logs {
// 			hashes = append(hashes, log.TransactionHash)
// 		}
// 		return hashes, nil
// 	}
// }
//
// type SimpleLog struct {
// 	BlockNumber      uint64 `json:"blockNumber"`
// 	TransactionIndex uint64 `json:"transactionIndex"`
// 	LogIndex         uint64 `json:"logIndex"`
// 	TransactionHash  string `json:"transactionHash"`
// 	Timestamp        int64  `json:"timestamp"`
// 	Address          string `json:"address"`
// 	Topic0           string `json:"topic0"`
// 	Topic1           string `json:"topic1"`
// 	Topic2           string `json:"topic2"`
// 	Topic3           string `json:"topic3"`
// 	Data             string `json:"data"`
// 	CompressedLog    string `json:"compressedLog"`
// }
