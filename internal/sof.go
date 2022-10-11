package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

func findSourceOfFunds(w *os.File, i int, hash, chain string) {
	if !isValidHash(hash) {
		// ignore header
		return
	}

	fmt.Fprintln(w, "\n", colors.BrightBlack+strings.Repeat("-", 5), fmt.Sprintf("%d.", i), hash, chain, strings.Repeat("-", 70), colors.Off)

	// get the transaction we're interested in
	callParams := map[string]string{"chain": chain, "hash": hash}
	tx := chifra.TransactionsCommand(w, callParams, noFunc, transFunc)

	// if there's too many records, we bail out
	nRecords := chifra.ListCountCommand(w, chain, tx.Sender, nil, postListCount)
	if nRecords > 20000 {
		fmt.Fprintln(os.Stderr, colors.Yellow, "Skipping address", tx.Sender, "too many records:", nRecords, colors.Off)

	} else {
		//--------------------------------------------------------
		cmdArgs := []string{
			"export",
			"--no_header",
			"--fmt",
			"txt",
			"[{ADDRESS}]",
			"--chain",
			"[{CHAIN}]",
			"--logs",
			"--emitter",
			"[{TOKEN}]",
			"--reversed",
			"--first_block",
			"[{FIRST}]",
			"--last_block",
			"[{LAST}]",
			"--cache",
			"--articulate",
			"[{TOPIC}]"}
		fields := []string{"address", "chain", "token", "first", "last", "topic"}
		values := []string{tx.Sender, chain, tx.Token, firstBlocks[chain], tx.BlockNum, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}
		Replace(cmdArgs, fields, values)

		filterExport := func(line string) bool {
			return strings.Contains(line, "to:"+tx.Sender)
		}

		commandToString12(w, cmdArgs, filterExport, postExport)
	}
}

func Replace(inOut []string, fields, values []string) {
	for f, field := range fields {
		for i := 0; i < len(inOut); i++ {
			inOut[i] = strings.Replace(inOut[i], "[{"+strings.ToUpper(field)+"}]", values[f], -1)
		}
	}
}

func noFunc(line string) bool {
	return true
}

var transFunc = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	if !filter(strIn) {
		return
	}

	log := Cut(w, strIn, []int{12}, []string{}, true /* silent */, 0)[0]
	log = strings.Replace(log, "{name:transfer|inputs:{_to:", "", -1)
	log = strings.Replace(log, "|_value:", "\t", -1)
	log = strings.Replace(log, "|value:", "\t", -1)
	log = strings.Replace(log, "}|outputs:{_success:}}", "", -1)
	log = strings.Trim(log, "\n")

	names, locations := chifra.TransactionFields()
	strOut := strIn + "\t" + log
	out = Cut(w, strOut, locations, names, false /* silent */, 0)
	return
}

var postListCount = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	return Cut(w, strIn, []int{1}, []string{"cnt"}, true /* silent */, 1)
}

var postExport = func(w *os.File, strIn string, filter func(string) bool) (out []string) {
	logs := strings.Split(strIn, "\n")
	for _, ll := range logs {
		// logger.Log(logger.Progress, "Processing log ", i, " of ", len(logs))
		if !filter(ll) {
			continue
		}
		ll = strings.Replace(ll, "{name:transfer|inputs:{_to:", "", -1)
		ll = strings.Replace(ll, "|_value:", "\t", -1)
		ll = strings.Replace(ll, "|value:", "\t", -1)
		ll = strings.Replace(ll, "}|outputs:{_success:}}", "", -1)
		ll = strings.Trim(ll, "\n")
		strOut := ll // strIn + "\t" + ll
		// blockNumber	transactionIndex	logIndex	timestamp	address	topic0	topic1	topic2	topic3	data	compressedLog
		fieldLocs := []int{0, 1, 2, 3, 4, 5, 11}
		fieldNames := []string{"blockNum", "txId", "logId", "date/tx_id", "address", "hash", "log"}
		ln := Cut(w, strOut, fieldLocs, fieldNames, true /* silent */, 1)
		// fmt.Fprintln(w, "len:", len(ln))
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

var firstBlocks = map[string]string{
	"mainnet": "13868853",
	"gnosis":  "19747830",
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

type filterFunc func(string) bool
type postFunc func(*os.File, string, func(string) bool) []string

func commandToString12(w *os.File, args []string, filter filterFunc, post postFunc) []string {
	// fmt.Fprintln(w, colors.Green, "chifra", strings.Join(args, " "), colors.Off)
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

func isValidHash(hash string) bool {
	ok, err := validate.IsValidHex("tx_hash", hash, 32)
	return ok && err == nil
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
			if uint64(i) < max_rows {
				w := os.Stdout
				findSourceOfFunds(w, i, donation.TxHash, donation.Network)
			}
		}
	}

	return nil
}
