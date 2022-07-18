package data

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Query struct {
	Cmd   string
	Round Round
	Url   string
	Fn    string
}

func (q *Query) Execute(w io.Writer) error {
	log.Println("Retrieving: ", q.Url)

	response, err := http.Get(q.Url)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s", string(responseData))

	return nil
}
