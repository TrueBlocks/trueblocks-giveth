package data

import (
	"io/ioutil"
	"net/http"
	"os"
)

type Query struct {
	Cmd   string
	Round Round
	Url   string
	Fn    string
}

func (q *Query) Execute() error {
	response, err := http.Get(q.Url)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(q.Fn)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(responseData)

	return nil
}
