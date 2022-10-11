package chifra

import (
	"os"
	"strings"
)

var transactionCmd = []string{
	"transactions",
	"--no_header",
	"--articulate",
	"--fmt",
	"txt",
	"--chain",
	"[{CHAIN}]",
	"[{HASH}]",
}

var transFieldLocs = map[string][]int{
	"blockNum":   {0, 0},
	"txId":       {1, 1},
	"hash":       {2, 9},
	"date/tx_id": {3, 2},
	"sender":     {4, 4},
	"token":      {5, 5},
	"recip":      {6, 13},
	"amount":     {7, 14},
}
var transFields = []string{
	"blockNum", "txId",
	"hash", "date/tx_id",
	"sender", "token",
	"recip", "amount",
}

func TransactionFields() ([]string, []int) {
	return transFields, transactionFieldLocs()
}

func transactionFieldLocs() []int {
	var ret []int
	for _, o := range transFields {
		ret = append(ret, transFieldLocs[o][1])
	}
	return ret
}

func getTransValue(vals []string, fn string) string {
	if transFieldLocs[fn] == nil {
		return ""
	}
	return vals[transFieldLocs[fn][1]]
}

func TransactionsCommand(w *os.File, fields map[string]string, filter filterFunc, post postFunc) SimpleTransaction {
	cmdArgs := []string{}
	cmdCopy := transactionCmd
	for _, f := range cmdCopy {
		f = strings.Replace(f, "[{CHAIN}]", fields["chain"], -1)
		f = strings.Replace(f, "[{HASH}]", fields["hash"], -1)
		cmdArgs = append(cmdArgs, f)
	}
	vals := commandToStrings(w, cmdArgs, filter, post)
	return SimpleTransaction{
		BlockNum: getTransValue(vals, "blockNum"),
		Sender:   getTransValue(vals, "sender"),
		Token:    getTransValue(vals, "token"),
	}
}

type SimpleTransaction struct {
	Sender   string `json:"sender"`
	Token    string `json:"token"`
	BlockNum string `json:"blockNum"`
}
