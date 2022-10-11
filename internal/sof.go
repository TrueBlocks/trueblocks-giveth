package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/chifra"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func FindSourceOfFunds(w *os.File, place string, depth int, hash, chain string) {
	// let them know we're here
	fmt.Fprintln(w, "\n", colors.BrightBlack+strings.Repeat("-", 5), place, hash, chain, strings.Repeat("-", 70), colors.Off)

	// get the transaction we're interested in
	callParams := map[string]string{"chain": chain, "hash": hash}
	tx := chifra.TransactionsCommand(w, callParams, noFunc, transFunc)

	// if there's too many records, we bail out
	nRecords := chifra.ListCountCommand(w, chain, tx.Sender, nil, postListFunc)
	if nRecords > 20000 {
		fmt.Fprintln(os.Stderr, colors.Yellow, "Skipping address", tx.Sender, "too many records:", nRecords, colors.Off)

	} else {
		chifra.SourceOfFunds(w, tx, chain, depth, nil, postExportFunc)
	}
}

// noFunc
func noFunc(line string) bool {
	return true
}

// transFunc
var transFunc = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	if !filter(strIn) {
		return
	}

	names, locations := chifra.TransactionFields()
	strOut := strIn + "\t" + cleanLog(Cut(w, strIn, []int{12}, []string{}, true /* silent */, 0)[0])
	out = Cut(w, strOut, locations, names, false /* silent */, 0)
	return
}

// postListFunc
var postListFunc = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	return Cut(w, strIn, []int{1}, []string{"cnt"}, true /* silent */, 1)
}

// postExportFunc
var postExportFunc = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	logs := strings.Split(strIn, "\n")
	for _, ll := range logs {
		// logger.Log(logger.Progress, "Processing log ", i, " of ", len(logs))
		if !filter(ll) {
			continue
		}
		strOut := cleanLog(ll)
		fieldLocs := []int{0, 1, 2, 3, 4, 5, 11}
		fieldNames := []string{"blockNum", "txId", "logId", "date/tx_id", "address", "hash", "log"}
		ln := Cut(w, strOut, fieldLocs, fieldNames, true /* silent */, 1)
		lll := ln[len(ln)-1]
		isExchange, which := hasExchange(lll)
		isFriend, who := hasFriend(lll)
		if isExchange {
			fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping staking contract "+which+colors.Off)

		} else if isFriend {
			fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping friend "+who+colors.Off)

		} else {
			// fmt.Fprintln(w, colors.BrightBlack, "   ", strings.Repeat("-", 120), colors.Off)
			lll = strings.Replace(lll, "{name:Transfer|inputs:{_amount:", "", -1)
			lll = strings.Replace(lll, "_from:", "", -1)
			lll = strings.Replace(lll, "_to:", "", -1)
			lll = strings.Replace(lll, "}}", "", -1)
			p := strings.Split(lll, "|")
			ln = ln[:len(ln)-1]
			// fmt.Fprintln(w, "len:", len(ln), len(p))
			ln = append(ln, p...)
			// fmt.Fprintln(w, "len:", len(ln), ln)
			fieldLocs = []int{0, 1, 2, 3, 4, 5, 7, 8, 6}
			fieldNames = []string{"blockNum", "txId", "logId", "hash", "date/tx_id", "token", "sender", "recip", "amount"}
			bb := Cut(w, strings.Join(ln, "\t"), fieldLocs, fieldNames, false /* silent */, 1)
			out = append(out, strings.Join(bb, "\t")+"\n")
		}
	}
	return
}

func hasExchange(s string) (bool, string) {
	exs := []string{"0xc0dbdca66a0636236fabe1b3c16b1bd4c84bb1e1", "0x08ea9f608656a4a775ef73f5b187a2f1ae2ae10e", "0x55ff0cef43f0df88226e9d87d09fa036017f5586"}
	for _, e := range exs {
		if strings.Contains(s, "from:"+e) {
			return true, e
		}
	}
	return false, ""
}

func hasFriend(s string) (bool, string) {
	fs := []string{"0x839395e20bbb182fa440d08f850e6c7a8f6f0780"}
	for _, f := range fs {
		if strings.Contains(s, "from:"+f) {
			return true, f
		}
	}
	return false, ""
}

type txId struct {
	bn     string
	tx_id  string
	log_id string
	hash   string
}

func (id *txId) Set(fn, val string) {
	switch fn {
	case "blockNum":
		id.bn = val
	case "txId":
		id.tx_id = val
	case "logId":
		id.log_id = val
	case "hash":
		id.hash = val
	}
}

func (id *txId) Get() string {
	ret := "[" + id.bn + "." + id.tx_id
	if len(id.log_id) > 0 {
		ret += "." + id.log_id
	}
	ret += " " + id.hash + "]"
	return ret
}

type transfer struct {
	token     string
	sender    string
	recipient string
	date      string
}

func (id *transfer) Set(fn, val string) {
	switch fn {
	case "token":
		fallthrough
	case "token-name":
		id.token = val
	case "sender":
		fallthrough
	case "sender-name":
		id.sender = val
	case "recip":
		fallthrough
	case "recip-name":
		id.recipient = val
	case "date":
		parts := strings.Split(val, ":")
		id.date = parts[0] + ":" + parts[1]
	}
}

func (id *transfer) Get(w string) string {
	return colors.BrightBlue + id.date + colors.Off + " " + id.sender + " " + w + " " + "[{AMT}] of " + id.token + " to " + id.recipient
}

func Cut(w *os.File, line string, fields []int, fns []string, silent bool, depth int) []string {
	showAddrs := os.Getenv("ADDRS") == "true"
	theId := txId{}
	theTransfer := transfer{}
	var ret []string
	parts := strings.Split(line, "\t")
	for i, f := range fields {
		if f < len(parts) {
			indent := strings.Repeat("   ", depth)
			hidden := false
			fn := fmt.Sprintf("%sfield_%d", indent, i)
			if len(fns) > i {
				fn = fns[i]
				if fn == "blockNum" || fn == "txId" || fn == "logId" || fn == "hash" {
					hidden = true
				} else if fn == "token" || fn == "sender" || fn == "recip" {
					hidden = true
				}
			}

			if !silent {
				val := parts[f]
				if fn == "amount" {
					w := colors.BrightWhite + "donated" + colors.Off
					if depth > 0 {
						w = "sent"
					}
					msg := theTransfer.Get(w)
					val = strings.Replace(msg, "[{AMT}]", "["+colors.Yellow+utils.PadLeft(val, 25, '.')+colors.Off+"]", -1)

				} else if fn == "date/tx_id" {
					if !strings.Contains(val, "UTC") {
						ts, _ := strconv.ParseUint(val, 10, 64)
						d, _ := tslib.FromTsToDate(uint64(ts))
						val = d.Format("YYYY-MM-DD HH:mm:ss UTC")
					}
					theTransfer.Set("date", val)
					val += " " + theId.Get()
					hidden = true
				}

				if validate.IsValidAddress(val) || val == "0x0" {
					if name, err := chifra.AddressToName(val, showAddrs /* decorated */); err == nil {
						fn = fn + "-name"
						val = name.Name
					}
				}

				theId.Set(fn, val)
				theTransfer.Set(fn, val)
				if !hidden {
					// fmt.Fprintln(w, indent, utils.PadRight(fn+":", 12, ' '), val)
					fmt.Fprintln(w, indent, val)
				}
			}
			ret = append(ret, parts[f])
		}
	}
	return ret
}

func cleanLog(strIn string) string {
	ret := strIn
	ret = strings.Replace(ret, "{name:transfer|inputs:{_to:", "", -1)
	ret = strings.Replace(ret, "|_value:", "\t", -1)
	ret = strings.Replace(ret, "|value:", "\t", -1)
	ret = strings.Replace(ret, "}|outputs:{_success:}}", "", -1)
	ret = strings.Trim(ret, "\n")
	return ret
}

// getSofOptions processes command line options for the Rounds command
func getSofOptions(cmd *cobra.Command, args []string) (max_rows uint64, globals Globals, err error) {
	globals, err = GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}

	max_rows, _ = strconv.ParseUint(os.Getenv("MAX_ROWS"), 10, 64)
	if max_rows == 0 {
		max_rows = utils.NOPOS
	}

	return
}

func RunSourceOfFunds(cmd *cobra.Command, args []string) error {
	max_rows, globals, err := getSofOptions(cmd, args)
	if err != nil {
		return err
	}

	for _, round := range globals.Rounds {
		donations, _ := data.NewDonations(data.GetFilename("eligible", "csv", round), "csv", data.SortByHash)
		for i, donation := range donations {
			if uint64(i) < max_rows && validate.IsValidHash(donation.TxHash) {
				w := os.Stdout
				FindSourceOfFunds(w, fmt.Sprintf("%d-%d-%d.", round.Id, i, len(donations)), 0, donation.TxHash, donation.Network)
			}
		}
	}

	return nil
}
