package chifra

import (
	"os"
	"strconv"
)

func ChifraList(w *os.File, fields map[string]string, post postFunc) int64 {
	chain := fields["chain"]
	address := fields["sender"]
	var cmdArgs = []string{
		"list",
		"--no_header",
		"--count",
		"--fmt",
		"txt",
		"--chain",
		"[{CHAIN}]",
		"[{ADDRESS}]"}
	fieldList := []string{"address", "chain"}
	valueList := []string{address, chain}
	replace(cmdArgs, fieldList, valueList)

	ret := commandToFields(w, cmdArgs, nil, post)
	r, _ := strconv.ParseInt(ret[0], 10, 64)
	return r
}
