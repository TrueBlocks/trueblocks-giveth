package chifra

import (
	"encoding/json"
	"errors"
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

func commandToRecord[T SimpleTransfer](w *os.File, args []string) (T, error) {
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
