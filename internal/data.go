package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/utils"
	"github.com/spf13/cobra"
)

func RunData(cmd *cobra.Command, args []string) error {
	dataType, globals, err := getDataOptions(cmd, args)
	if err != nil {
		return err
	}

	var queries []data.Query
	if dataType == "purple-list" {
		queries = append(queries, getQuery(dataType, "json", data.Round{}))

	} else {
		for _, round := range globals.Rounds {
			queries = append(queries, getQuery(dataType, "csv", round))
		}
	}

	if globals.Script {
		for _, q := range queries {
			fmt.Fprintln(os.Stdout, "curl", "\""+q.Url+"\"", "--output", q.Fn, "; sleep", int(globals.Sleep))
		}
	} else {
		for _, q := range queries {
			if globals.Update {
				log.Println("Updating: ", q.Fn)
				q.Execute()
				if !globals.Verbose {
					if globals.Sleep == 0 {
						globals.Sleep = 2
					}
					log.Println("Sleeping for", globals.Sleep, "seconds")
					time.Sleep(globals.Sleep * time.Second)
				}
			}
		}
		data.RenderQueries(globals.Format, os.Stdout, queries)
	}

	return nil
}

func getUrl(cmd, format string, round data.Round) string {
	defOpts := "startDate=[SD]&endDate=[ED]&download=[DL]"
	var opts = map[string]string{
		"purple-list":     "[CMD]",
		"eligible":        "[CMD]?[OPTS]&justCountListed=no",
		"not-eligible":    "[CMD]?[OPTS]&justCountListed=no",
		"purple-verified": "[CMD]?[OPTS]",
		"calc-givback":    "[CMD]?[OPTS]&givPrice=[PRICE]&givAvailable=[AVAIL]&givMaxFactor=0.75&relayerAddress=0xd0e81E3EE863318D0121501ff48C6C3e3Fd6cbc7&maxAddressesPerFunctionCall=200",
	}

	url := "https://givback.develop.giveth.io/" + opts[cmd]
	url = strings.Replace(url, "[OPTS]", defOpts, -1)
	url = strings.Replace(url, "[CMD]", data.Cmds[cmd], -1)
	url = strings.Replace(url, "[SD]", utils.GetGivethDate(round.StartDate), -1)
	url = strings.Replace(url, "[ED]", utils.GetGivethDate(round.EndDate), -1)
	if format == "txt" || format == "csv" {
		url = strings.Replace(url, "[DL]", "yes", -1)
	} else {
		url = strings.Replace(url, "[DL]", "no", -1)
	}

	if cmd == "calc-givback" {
		url = strings.Replace(url, "[PRICE]", fmt.Sprintf("%g", round.Price), -1)
		url = strings.Replace(url, "[AVAIL]", fmt.Sprintf("%d", round.Available), -1)
		return url
	}
	return url
}

func getQuery(cmd, format string, round data.Round) data.Query {
	return data.Query{Cmd: cmd, Url: getUrl(cmd, format, round), Fn: data.GetFilename(cmd, format, round)}
}

func DataTypes() []string {
	return []string{
		"purple-list",
		"not-eligible",
		"eligible",
		"calc-givback",
		"purple-verified",
	}
}

func isValidType(dataType string) bool {
	for _, d := range DataTypes() {
		if d == dataType {
			return true
		}
	}
	return false
}

// getDataOptions processes command line options for the Rounds command
func getDataOptions(cmd *cobra.Command, args []string) (dataType string, globals Globals, err error) {
	globals, err = GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}

	dataType, err = cmd.Flags().GetString("data")
	if err != nil {
		return
	}

	if globals.Script {
		if globals.Sleep == 0 {
			// If it's unset and we're going to be calling the API, set it
			globals.Sleep = 3
		}
	}

	if len(dataType) > 0 {
		// we got one, check if it's valid. If yes, we're done
		if isValidType(dataType) {
			err = validateData(dataType, globals)
			return
		}
	}

	for _, arg := range args {
		if isValidType(arg) {
			dataType = arg
			err = validateData(dataType, globals)
			return
		}
	}

	if len(args) > 0 {
		err = fmt.Errorf("invalid option '%s'", args[0])
	} else {
		err = fmt.Errorf("flag needs an argument: --data")
	}

	return
}

func validateData(dataType string, globals Globals) (err error) {
	if dataType == "purple-list" {
		return nil
	}

	return
}
