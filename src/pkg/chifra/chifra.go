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

func TraceSourceOfFunds(w *os.File, depth, max_depth int, hash, chain string) error {
	callParams := map[string]string{"chain": chain, "hash": hash}
	if tr, err := ChifraTransactions(w, callParams); err != nil {
		return err

	} else {
		dispo := "donated"
		if depth > 0 {
			dispo = "sent"
		}
		arrow := strings.Repeat("-", (depth*2)) + ">"
		fmt.Fprintln(w, arrow+strings.Replace(tr.String(), "[{DISPO}]", dispo, -1))
		if depth+1 < max_depth {
			if nRecords, err := ChifraListCount(w, chain, tr.Sender); err != nil {
				return err

			} else {
				if nRecords > 20000 {
					if os.Getenv("SHOW_SKIPS") == "true" {
						fmt.Fprintln(os.Stderr, colors.Yellow, "Skipping address", tr.Sender, "too many records:", nRecords, colors.Off)
					}
					return nil
				}

				one := os.Getenv("NEW") != "true"
				if one {
					ChifraExport(w, tr, chain, depth, max_depth, postExportFunc)
				} else {
					tr.Sender = tr.Recipient
					processedMap := map[string]bool{}
					hashes, err := GetAppearanceHashes(w, tr, chain, processedMap)
					if err != nil {
						return err
					}
					for _, h := range hashes {
						if h != hash {
							if err := TraceSourceOfFunds(w, depth+1, max_depth, h, chain); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

const transferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func GetAppearanceHashes(w *os.File, tr *SimpleTransfer, chain string, hashes map[string]bool) ([]string, error) {
	verbose := os.Getenv("SHOW_SKIPS") == "true"
	if tr.Token == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		return nil, nil
	}

	args := []string{
		"export",
		"--logs",
		"--no_header",
		"--cache",
		"--reversed",
		"--articulate",
		"--fmt", "json",
		"--chain", chain,
		"--first_block", firstBlocks[chain],
		"--last_block", fmt.Sprintf("%d", tr.BlockNumber),
		"--emitter", tr.Token,
		transferTopic,
		tr.Sender,
	}

	if logs, err := commandToRecords[SimpleLog](w, args); err != nil {
		return nil, err

	} else {
		newHashes := []string{}
		for _, log := range logs {
			isExchange, which := hasExchange(log.CompressedLog)
			isFriend, who := hasFriend(log.CompressedLog)
			if isExchange {
				if verbose {
					fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping staking contract "+which+colors.Off)
				}

			} else if isFriend {
				if verbose {
					fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping friend "+who+colors.Off)
				}
			} else {
				if strings.Contains(log.CompressedLog, "to:"+tr.Sender) {
					// parts := strings.Split(log.CompressedLog, "|")
					// fmt.Printf("%9d-%5d-%5d: %s %s\n", log.BlockNumber, log.TransactionIndex, log.LogIndex, parts[len(parts)-1], log.TransactionHash)
					hashes[log.TransactionHash] = true
					newHashes = append(newHashes, log.TransactionHash)
				}
			}
		}
		return newHashes, nil
	}
}

type SimpleLog struct {
	BlockNumber      uint64 `json:"blockNumber"`
	TransactionIndex uint64 `json:"transactionIndex"`
	LogIndex         uint64 `json:"logIndex"`
	TransactionHash  string `json:"transactionHash"`
	Timestamp        int64  `json:"timestamp"`
	Address          string `json:"address"`
	Topic0           string `json:"topic0"`
	Topic1           string `json:"topic1"`
	Topic2           string `json:"topic2"`
	Topic3           string `json:"topic3"`
	Data             string `json:"data"`
	CompressedLog    string `json:"compressedLog"`
}

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

func commandToRecord[T SimpleTransaction | SimpleListCount | SimpleLog](w *os.File, args []string) (*T, error) {
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

func commandToRecords[T SimpleTransaction | SimpleListCount | SimpleLog](w *os.File, args []string) ([]T, error) {
	// fmt.Println(colors.Green, "Calling: chifra", strings.Join(args, " "), colors.Off)
	if bytes, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", ("chifra " + strings.Join(args, " ")), "error:", err)
		os.Exit(1)
	} else {
		if isEmpty(bytes) {
			return []T{}, nil
		}
		resp := ChifraResponse[T]{}
		if err := json.Unmarshal(bytes, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	return nil, nil
}
func isEmpty(bytes []byte) bool {
	return false
}

// -----------
type ChifraResponse[T SimpleTransaction | SimpleListCount | SimpleLog] struct {
	Data []T `json:"data"`
}