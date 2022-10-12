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
type postFunc func(*os.File, string, func(string) bool) []string

func commandToFields(w *os.File, args []string, filter filterFunc, post postFunc) []string {
	if ret, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", err)
		os.Exit(1)
	} else {
		if post == nil {
			log.Fatal("You must provide a post processing function")
		}
		r := string(ret)
		s := strings.Split(r, "\t")
		if len(s) > 0 {
			return post(w, r, filter)
		}
	}
	return []string{}
}

func commandToRecord[T SimpleTransfer | SimpleListCount](w *os.File, args []string) (T, error) {
	if bytes, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", err)
		os.Exit(1)
	} else {
		resp := ChifraResponse[T]{}
		if err := json.Unmarshal(bytes, &resp); err != nil {
			return T{}, err
		} else if len(resp.Data) == 0 {
			return T{}, errors.New("transaction not found " + strings.Join(args, " "))
		}
		return resp.Data[0], nil
	}
	return T{}, nil
}

type ChifraResponse[T SimpleTransfer] struct {
	Data []T `json:"data"`
}

func TraceSourceOfFunds(w *os.File, depth, max_depth int, hash, chain string) error {
	callParams := map[string]string{"chain": chain, "hash": hash}
	if tx, err := ChifraTransactions(w, callParams); err != nil {
		return err

	} else {
		callParams = map[string]string{"chain": chain, "sender": tx.Sender}
		nRecords := ChifraList(w, callParams, postListFunc)
		if nRecords > 20000 {
			verbose := os.Getenv("SHOW_SKIPS") == "true"
			if verbose {
				fmt.Fprintln(os.Stderr, colors.Yellow, "Skipping address", tx.Sender, "too many records:", nRecords, colors.Off)
			}
			return nil
		}

		ChifraExport(w, *tx, chain, depth, nil, postExportFunc)
	}

	return nil
}
