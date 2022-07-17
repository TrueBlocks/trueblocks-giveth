package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/utils"
	"github.com/bykof/gostradamus"
	"github.com/spf13/cobra"
)

type Query struct {
	Round data.Round
	Url   string
	Fn    string
	Sleep int
}

func RunDownload(cmd *cobra.Command, args []string) error {
	rounds, dataType, format, err := getDownloadOptions(cmd, args)
	if err != nil {
		return err
	}

	var queries []Query
	if dataType == "purple-list" {
		url, fn := getParams("purpleList", format, data.Round{})
		queries = append(queries, Query{Url: url, Fn: fn, Sleep: 3})

	} else {
		for _, round := range rounds {
			if round.StartDate.Time().Before(gostradamus.Now().Time()) {
				switch dataType {
				case "not-eligible":
					url, fn := getParams("not-eligible-donations", format, round)
					queries = append(queries, Query{Round: round, Url: url, Fn: fn, Sleep: 3})

				case "eligible":
					url, fn := getParams("eligible-donations", format, round)
					queries = append(queries, Query{Round: round, Url: url, Fn: fn, Sleep: 3})

				case "calc-givback":
					url, fn := getParams("calculate-givback", format, round)
					queries = append(queries, Query{Round: round, Url: url, Fn: fn, Sleep: 3})

				case "purple-verified":
					url, fn := getParams("purpleList-donations-to-verifiedProjects", format, round)
					queries = append(queries, Query{Round: round, Url: url, Fn: fn, Sleep: 3})

				}
			}
		}
	}

	for _, q := range queries {
		fmt.Println("curl", "\""+q.Url+"\"", "--output", q.Fn, "; sleep", q.Sleep)
	}

	return nil
}

func getUrl(cmd, format string, round data.Round) string {
	baseUrl := "https://givback.develop.giveth.io"
	if round.Id == 0 {
		return baseUrl + "/" + cmd
	}
	sd := utils.GetGivethDate(round.StartDate)
	ed := utils.GetGivethDate(round.EndDate)
	m := make(map[string]string)
	m["eligible-donations"] = "%s/%s?startDate=%s&endDate=%s%s&justCountListed=no"
	m["not-eligible-donations"] = "%s/%s?startDate=%s&endDate=%s%s&justCountListed=no"
	m["purpleList-donations-to-verifiedProjects"] = "%s/%s?startDate=%s&endDate=%s%s"
	m["calculate-givback"] = "%s/%s?startDate=%s&endDate=%s%s&givPrice=%s&givAvailable=%d&givMaxFactor=0.75&relayerAddress=0xd0e81E3EE863318D0121501ff48C6C3e3Fd6cbc7&maxAddressesPerFunctionCall=200"

	download := ""
	if format == "txt" || format == "csv" {
		download = "&download=yes"
	}
	if cmd == "calculate-givback" {
		return fmt.Sprintf(m[cmd], baseUrl, cmd, sd, ed, download, round.Price, round.Available)
	}

	return fmt.Sprintf(m[cmd], baseUrl, cmd, sd, ed, download)
}

func getGivethFilename(cmd, format string, round data.Round) string {
	if round.Id == 0 {
		return "data/" + cmd + "/" + cmd + "." + format
	}
	sd := utils.GetGivethDate(round.StartDate)
	ed := utils.GetGivethDate(round.EndDate)
	fn := "data/" + cmd + "/" + cmd + "-" + sd + "-" + ed
	fn = fmt.Sprintf("%s-Round%02d.%s", fn, round.Id, format)
	return strings.Replace(strings.Replace(fn, "%2F", "_", -1), "%3A", "_", -1)
}

func getParams(cmd, format string, round data.Round) (string, string) {
	return getUrl(cmd, round), getGivethFilename(cmd, format, round)
}

func isValidType(dataType string) bool {
	types := DataTypes()
	for _, d := range types {
		if d == dataType {
			return true
		}
	}
	return false
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

// getDownloadOptions processes command line options for the Rounds command
func getDownloadOptions(cmd *cobra.Command, args []string) (rounds []data.Round, dataType string, format string, err error) {
	format, err = cmd.Flags().GetString("fmt")
	if err != nil {
		log.Fatal(err)
	}

	dataType, err = cmd.Flags().GetString("data")
	if err != nil {
		log.Fatal(err)
	}

	if len(dataType) > 0 {
		// we got one, check if it's valid. If yes, we're done
		if isValidType(dataType) {
			return data.GetRounds(), dataType, format, nil
		}
	}

	for _, arg := range args {
		if isValidType(arg) {
			return data.GetRounds(), arg, format, nil
		}
	}

	if len(args) > 0 {
		return data.GetRounds(), "", format, fmt.Errorf("invalid option '%s'", args[0])
	}
	return data.GetRounds(), "", format, fmt.Errorf("flag needs an argument: --data")
}
