package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/utils"
	"github.com/spf13/cobra"
)

func RunCycles(cmd *cobra.Command, args []string) error {
	_, globals, err := getCycleOptions(cmd, args)
	if err != nil {
		return err
	}

	CleanFolder(globals.Rounds[0].Id)

	donations, _ := data.NewDonations(data.GetFilename("eligible", "csv", globals.Rounds[0]), "csv", data.NoSort)

	extractTransations(donations, "gnosis", globals.Rounds[0].Id)
	donors := extractDonors(donations, "gnosis", globals.Rounds[0].Id)
	for _, donor := range donors {
		extractNeighbors(donor, "gnosis", globals.Rounds[0].Id, globals.Rounds[0].GnosisRange.Last)
	}

	extractTransations(donations, "mainnet", globals.Rounds[0].Id)
	donors = extractDonors(donations, "mainnet", globals.Rounds[0].Id)
	for _, donor := range donors {
		extractNeighbors(donor, "mainnet", globals.Rounds[0].Id, globals.Rounds[0].MainnetRange.Last)
	}

	return nil
}

func extractNeighbors(donor, chain string, id int, last uint64) {
	if donor == "0x101a52c72e8b48e9be5bcbb63e7c1fd140861bae" {
		return
	}
	first := "19719093"
	if chain == "mainnet" {
		first = "13858106"
	}
	logger.Log(logger.Info, "Getting neighbors for", chain, donor)
	output := utils.SystemCall("chifra", []string{
		"export", "--chain", chain, "--no_header", "--fmt", "txt", "--first_block", first, "--last_block", fmt.Sprintf("%d", last), "--neighbors", donor,
	}, []string{})
	neighbors := strings.Split(strings.Replace(strings.Replace(output, "\"", "", -1), "\\n", "\n", -1), "\n")
	stringsToFile2(id, chain, "raw", donor, strings.Join(neighbors, "\n"))

	raw_senders := []string{}
	senders := []string{}
	for _, neighbor := range neighbors {
		if strings.Contains(neighbor, "from") && !strings.Contains(neighbor, donor) {
			raw_senders = Unique(append(raw_senders, neighbor))
			parts := strings.Split(neighbor, "\t")
			senders = Unique(append(senders, parts[2]))
		}
	}
	// sort.Slice(raw_senders, func(i, j int) bool {
	// 	return raw_senders[i] < raw_senders[j]
	// })
	stringsToFile2(id, chain, "raw_senders", donor, strings.Join(raw_senders, "\n")+"\n")

	sort.Slice(senders, func(i, j int) bool {
		return senders[i] < senders[j]
	})
	stringsToFile2(id, chain, "senders", donor, strings.Join(senders, "\n")+"\n")
}

func stringsToFile2(id int, s, t, donor, str string) string {
	fn := fmt.Sprintf("./Round_%03d/neighbors/%s/%s/%s.txt", id, s, t, donor)
	file.StringToAsciiFile(fn, str)
	return fn
}

func extractDonors(donations []data.Donation, chain string, id int) []string {
	donors := []string{}
	for _, donation := range donations {
		if donation.Network == chain {
			donors = append(donors, donation.GiverAddress)
		}
	}
	sort.Slice(donors, func(i, j int) bool {
		return donors[i] < donors[j]
	})

	out := []string{}
	for _, donor := range donors {
		out = Unique(append(out, donor))
	}

	logger.Log(logger.Info, "Extracting donor addresses for", chain)
	output := strings.Join(out, "\n") + "\n"
	stringsToFile(id, "results", chain, "donors", output)

	return out
}

func extractTransations(donations []data.Donation, chain string, id int) {
	logger.Log(logger.Info, "Getting hashes for", chain)
	var g []string
	for _, item := range donations {
		if item.Network == chain {
			g = append(g, item.TxHash)
		}
	}

	logger.Log(logger.Info, "Getting transaction details for", chain)
	fn := stringsToFile(id, "results", chain, "txhashes", strings.Join(g, "\n")+"\n")
	output := utils.SystemCall("chifra", []string{
		"transactions", "--chain", chain, "--no_header", "--fmt", "csv", "--file", fn,
	}, []string{})
	stringsToFile(id, "results", chain, "txs", strings.Replace(strings.Replace(output, "\"", "", -1), "\\n", "\n", -1))
}

// getCycleOptions processes command line options for the Rounds command
func getCycleOptions(cmd *cobra.Command, args []string) (dataType string, globals Globals, err error) {
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

	// if len(dataType) > 0 {
	// 	// we got one, check if it's valid. If yes, we're done
	// 	if isValidType(dataType) {
	// 		err = validate(dataType, globals)
	// 		return
	// 	}
	// }
	// for _, arg := range args {
	// 	if isValidType(arg) {
	// 		dataType = arg
	// 		err = validate(dataType, globals)
	// 		return
	// 	}
	// }
	// if len(args) > 0 {
	// 	err = fmt.Errorf("invalid option '%s'", args[0])
	// } else {
	// 	err = fmt.Errorf("flag needs an argument: --data")
	// }

	// globals.Rounds = []data.Round{}

	if len(globals.Rounds) == 0 || len(globals.Rounds) > 1 {
		err = fmt.Errorf("exactly one round is required")
	}

	return
}

func stringsToFile(id int, f, s, t, str string) string {
	fn := fmt.Sprintf("./Round_%03d/%s/%s/%s_%d.csv", id, f, s, t, id)
	file.StringToAsciiFile(fn, str)
	return fn
}

func Unique(in []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, i := range in {
		if _, ok := inResult[i]; !ok {
			inResult[i] = true
			result = append(result, i)
		}
	}
	return result
}

func CleanFolder(id int) {
	folder := fmt.Sprintf("Round_%03d", id)
	os.RemoveAll("./" + folder)
	file.EstablishFolder(folder + "/results/gnosis/")
	file.EstablishFolder(folder + "/results/mainnet/")
	file.EstablishFolder(folder + "/neighbors/gnosis/raw/")
	file.EstablishFolder(folder + "/neighbors/gnosis/raw_senders/")
	file.EstablishFolder(folder + "/neighbors/gnosis/senders/")
	file.EstablishFolder(folder + "/neighbors/mainnet/raw/")
	file.EstablishFolder(folder + "/neighbors/mainnet/raw_senders/")
	file.EstablishFolder(folder + "/neighbors/mainnet/senders/")
}
