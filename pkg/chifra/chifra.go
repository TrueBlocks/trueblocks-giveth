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

func commandToStrings(w *os.File, args []string, filter filterFunc, post postFunc) []string {
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
