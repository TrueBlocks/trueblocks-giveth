package chifra

import (
	"os"
	"strings"
)

func SourceOfFunds(w *os.File, tx SimpleTransaction, chain string, depth int, filter filterFunc, post postFunc) {
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
	replace(cmdArgs, fields, values)

	filterExport := func(line string) bool {
		return strings.Contains(line, "to:"+tx.Sender)
	}

	commandToStrings(w, cmdArgs, filterExport, post)
}

var firstBlocks = map[string]string{
	"mainnet": "13868853",
	"gnosis":  "19747830",
}
