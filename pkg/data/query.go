package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
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
	f.Write(responseData)
	f.Close()

	return q.Enhance()
}

func (q *Query) Enhance() error {
	_, typ, _, round := ExplodeFilename(q.Fn)
	if typ == "eligible" || typ == "not-eligible" || typ == "purple-verified" || typ == "calc-givback" {
		lines := file.AsciiFileToLines(q.Fn)

		out, err := os.Create(q.Fn)
		if err != nil {
			return err
		}
		defer out.Close()

		for i, line := range lines {
			if typ == "eligible" || typ == "not-eligible" || typ == "purple-verified" {
				parts := strings.Split(line, ",")
				if len(parts) >= 12 {
					line = strings.Join(parts[0:11], ",")
				}
			}
			var l string
			if i == 0 {
				l = fmt.Sprintf("\"%s\",\"%s\",%s", "type", "round", line)
			} else {
				l = fmt.Sprintf("\"%s\",\"%s\",%s", typ, fmt.Sprintf("Round%02d", round), line)
			}
			l = strings.Trim(l, ",")
			for {
				if !strings.Contains(l, ",,") {
					break
				}
				l = strings.Replace(l, ",,", ",\"\",", -1)
			}
			for {
				if strings.Count(l, ",") >= 12 {
					break
				}
				l = l + ",\"\""
			}
			l += "\n"
			out.Write([]byte(l))
		}
	}
	return nil
}
