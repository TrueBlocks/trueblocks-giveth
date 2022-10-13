package chifra

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

func ChifraExport(w *os.File, tr *SimpleTransaction, chain string, depth, max_depth int, post postFunc) {
	args := []string{
		"export",
		"--no_header",
		"--fmt",
		"txt",
		"[{ADDRESS}]",
		"--chain",
		"[{CHAIN}]",
		"--logs",
		"--reversed",
		"--first_block",
		"[{FIRST}]",
		"--last_block",
		"[{LAST}]",
		"--cache",
		"--articulate",
		"[{TOPIC}]",
	}

	if tr.Token != "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		args = append(args, "--emitter")
		args = append(args, "[{TOKEN}]")
	} else {
		return
	}

	fields := []string{"address", "chain", "token", "first", "last", "topic"}
	values := []string{tr.Sender, chain, tr.Token, firstBlocks[chain], fmt.Sprintf("%d", tr.BlockNumber), "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}
	for f, field := range fields {
		for i := 0; i < len(args); i++ {
			args[i] = strings.Replace(args[i], "[{"+strings.ToUpper(field)+"}]", values[f], -1)
		}
	}

	filterExport := func(line string) bool {
		return strings.Contains(line, "to:"+tr.Sender)
	}

	commandToFields(w, args, filterExport, post)
}

var firstBlocks = map[string]string{
	"mainnet": "13868853",
	"gnosis":  "19747830",
}

// postExportFunc
var postExportFunc = func(w *os.File, strIn string, filter func(string) bool) {
	verbose := os.Getenv("SHOW_SKIPS") == "true"
	logs := strings.Split(strIn, "\n")
	for _, ll := range logs {
		// logger.Log(logger.Progress, "Processing log ", i, " of ", len(logs))
		if !filter(ll) {
			continue
		}
		strOut := cleanLog(ll)
		fieldLocs := []int{0, 1, 2, 3, 4, 5, 11}
		ln := cut(strOut, fieldLocs)
		lll := ln[len(ln)-1]
		isExchange, which := hasExchange(lll)
		isFriend, who := hasFriend(lll)
		if isExchange {
			if verbose {
				fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping staking contract "+which+colors.Off)
			}

		} else if isFriend {
			if verbose {
				fmt.Fprintln(os.Stderr, colors.Yellow+"    Skipping friend "+who+colors.Off)
			}

		} else {
			lll = strings.Replace(lll, "{name:Transfer|inputs:{_amount:", "", -1)
			lll = strings.Replace(lll, "_from:", "", -1)
			lll = strings.Replace(lll, "_to:", "", -1)
			lll = strings.Replace(lll, "}}", "", -1)
			p := strings.Split(lll, "|")
			ln = ln[:len(ln)-1]
			ln = append(ln, p...)
			fieldLocs = []int{0, 1, 2, 3, 4, 5, 7, 8, 6}
			fieldNames := []string{"blockNum", "txId", "logId", "hash", "date/tx_id", "token", "sender", "recip", "amount"}
			if ln[6] != "0" { // do not show zero value stuff or ether details
				showLine(w, strings.Join(ln, "\t"), fieldLocs, fieldNames, false /* silent */, 1)
			}
		}
	}
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

func showLine(w *os.File, line string, fields []int, fns []string, silent bool, depth int) []string {
	showAddrs := os.Getenv("SHOW_ADDRS") == "true"
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
					if name, err := AddressToName(val, showAddrs /* decorated */); err == nil {
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

func cut(line string, fields []int) []string {
	var ret []string
	parts := strings.Split(line, "\t")
	for _, f := range fields {
		if f < len(parts) {
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

func hasExchange(s string) (bool, string) {
	exs := []string{
		"0xc0dbdca66a0636236fabe1b3c16b1bd4c84bb1e1",
		"0x08ea9f608656a4a775ef73f5b187a2f1ae2ae10e",
		"0x55ff0cef43f0df88226e9d87d09fa036017f5586",
	}

	for _, e := range exs {
		if strings.Contains(s, "from:"+e) {
			return true, e
		}
	}

	return false, ""
}

func hasFriend(s string) (bool, string) {
	fs := []string{
		"0x839395e20bbb182fa440d08f850e6c7a8f6f0780",
	}

	for _, f := range fs {
		if strings.Contains(s, "from:"+f) {
			return true, f
		}
	}

	return false, ""
}
