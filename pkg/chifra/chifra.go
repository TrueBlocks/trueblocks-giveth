package chifra

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

func commandToLine(w *os.File, args []string, filter filterFunc, post postFunc) []byte {
	if bytes, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", err)
		os.Exit(1)
	} else {
		return bytes
	}
	return []byte{}
}
