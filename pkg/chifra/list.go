package chifra

import (
	"os"
	"strconv"
	"strings"
)

func ListCountCommand(w *os.File, chain, address string, filter filterFunc, post postFunc) int64 {
	var cmdArgs = []string{
		"list",
		"--no_header",
		"--count",
		"--fmt",
		"txt",
		"--chain",
		"[{CHAIN}]",
		"[{ADDRESS}]"}
	fields := []string{"address", "chain"}
	values := []string{address, chain}
	replace(cmdArgs, fields, values)

	ret := commandToFields(w, cmdArgs, filter, post)
	r, _ := strconv.ParseInt(ret[0], 10, 64)
	return r
}

func replace(inOut []string, fields, values []string) {
	for f, field := range fields {
		for i := 0; i < len(inOut); i++ {
			inOut[i] = strings.Replace(inOut[i], "[{"+strings.ToUpper(field)+"}]", values[f], -1)
		}
	}
}
