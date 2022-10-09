package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func RunSourceOfFunds(cmd *cobra.Command, args []string) error {
	hash, chain, _, err := getSofOptions(cmd, args)
	if err != nil {
		return err
	}

	cmdArgs := []string{"transactions", "--no_header", "--articulate", "--fmt", "txt", "--chain", chain, hash}
	logger.Log(logger.Info, colors.Green, "Calling", "chifra", strings.Join(cmdArgs, " "), colors.Off)
	fmt.Println(utils.PadRight("hash:", 10, ' '), hash)
	vals := Cut(commandToString(cmdArgs), []int{0, 1, 2, 4, 5, 12}, []string{"bn", "tx", "ts", "from", "to", "art"})
	fmt.Println(strings.Join(vals, " "))

	return nil
}

func Cut(str string, fields []int, fns []string) []string {
	var ret []string
	parts := strings.Split(str, "\t")
	for i, f := range fields {
		if f < len(parts) {
			fmt.Println(utils.PadRight(fns[i]+":", 10, ' '), parts[f])
			if validate.IsValidAddress(parts[f]) {
				if name, err := data.GetChifraName(parts[f]); err == nil {
					fmt.Println(utils.PadRight(fns[i]+"-name:", 10, ' '), name.Name)
				}
			}
			ret = append(ret, parts[f])
		}
	}
	return ret
}

type UniqAppearance struct {
	BlockNumber      uint   `json:"bn"`
	TransactionIndex uint   `json:"tx"`
	TraceIndex       string `json:"tc"`
	Address          string `json:"addr"`
	Reason           string `json:"reason"`
}

func stringToAppearance(in string) UniqAppearance {
	in = strings.Replace(strings.Replace(in, "{ \"data\": [", "", -1), "] }", "", -1)
	var v UniqAppearance
	err := json.Unmarshal([]byte(in), &v)
	if err != nil {
		fmt.Println(err)
	}
	return v
}

func commandToString(args []string) string {
	if ret, err := exec.Command("chifra", args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running the command: ", err)
		os.Exit(1)
	} else {
		return string(ret)
	}
	return ""
}

// func extractNeighbors(donor, chain string, id int, last uint64) {
// 	if donor == "0x101a52c72e8b48e9be5bcbb63e7c1fd140861bae" {
// 		return
// 	}
// 	first := "19719093"
// 	if chain == "mainnet" {
// 		first = "13858106"
// 	}
// 	logger.Log(logger.Info, "Getting neighbors for", chain, donor)
// 	output := utils.SystemCall("chifra", []string{
// 		"export", "--chain", chain, "--no_header", "--fmt", "txt", "--first_block", first, "--last_block", fmt.Sprintf("%d", last), "--neighbors", donor,
// 	}, []string{})
// 	neighbors := strings.Split(strings.Replace(strings.Replace(output, "\"", "", -1), "\\n", "\n", -1), "\n")
// 	stringsToFile2(id, chain, "raw", donor, strings.Join(neighbors, "\n"))

// 	raw_senders := []string{}
// 	senders := []string{}
// 	for _, neighbor := range neighbors {
// 		if strings.Contains(neighbor, "from") && !strings.Contains(neighbor, donor) {
// 			raw_senders = Unique(append(raw_senders, neighbor))
// 			parts := strings.Split(neighbor, "\t")
// 			senders = Unique(append(senders, parts[2]))
// 		}
// 	}
// 	// sort.Slice(raw_senders, func(i, j int) bool {
// 	// 	return raw_senders[i] < raw_senders[j]
// 	// })
// 	stringsToFile2(id, chain, "raw_senders", donor, strings.Join(raw_senders, "\n")+"\n")

// 	sort.Slice(senders, func(i, j int) bool {
// 		return senders[i] < senders[j]
// 	})
// 	stringsToFile2(id, chain, "senders", donor, strings.Join(senders, "\n")+"\n")
// }

// func stringsToFile2(id int, s, t, donor, str string) string {
// 	fn := fmt.Sprintf("./Round_%03d/neighbors/%s/%s/%s.txt", id, s, t, donor)
// 	file.StringToAsciiFile(fn, str)
// 	return fn
// }

// func extractDonors(donations []data.Donation, chain string, id int) []string {
// 	donors := []string{}
// 	for _, donation := range donations {
// 		if donation.Network == chain {
// 			donors = append(donors, donation.GiverAddress)
// 		}
// 	}
// 	sort.Slice(donors, func(i, j int) bool {
// 		return donors[i] < donors[j]
// 	})

// 	out := []string{}
// 	for _, donor := range donors {
// 		out = Unique(append(out, donor))
// 	}

// 	logger.Log(logger.Info, "Extracting donor addresses for", chain)
// 	output := strings.Join(out, "\n") + "\n"
// 	stringsToFile(id, "results", chain, "donors", output)

// 	return out
// }

// func extractTransations(donations []data.Donation, chain string, id int) {
// 	logger.Log(logger.Info, "Getting hashes for", chain)
// 	var g []string
// 	for _, item := range donations {
// 		if item.Network == chain {
// 			g = append(g, item.TxHash)
// 		}
// 	}

// 	logger.Log(logger.Info, "Getting transaction details for", chain)
// 	fn := stringsToFile(id, "results", chain, "txhashes", strings.Join(g, "\n")+"\n")
// 	output := utils.SystemCall("chifra", []string{
// 		"transactions", "--chain", chain, "--no_header", "--fmt", "csv", "--file", fn,
// 	}, []string{})
// 	stringsToFile(id, "results", chain, "txs", strings.Replace(strings.Replace(output, "\"", "", -1), "\\n", "\n", -1))
// }

// getSofOptions processes command line options for the Rounds command
func getSofOptions(cmd *cobra.Command, args []string) (hash, chain string, globals Globals, err error) {
	globals, err = GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}

	hash, err = cmd.Flags().GetString("hash")
	if err != nil {
		return
	}
	ok, err := validate.IsValidHex("tx_hash", hash, 32)
	if !ok || err != nil {
		return
	}

	chain, err = cmd.Flags().GetString("chain")
	if err != nil {
		return
	}
	if chain != "mainnet" && chain != "gnosis" {
		err = errors.New("Invalid chain " + chain)
		return
	}

	return
}

// func stringsToFile(id int, f, s, t, str string) string {
// 	fn := fmt.Sprintf("./Round_%03d/%s/%s/%s_%d.csv", id, f, s, t, id)
// 	file.StringToAsciiFile(fn, str)
// 	return fn
// }

// func Unique(in []string) []string {
// 	inResult := make(map[string]bool)
// 	var result []string
// 	for _, i := range in {
// 		if _, ok := inResult[i]; !ok {
// 			inResult[i] = true
// 			result = append(result, i)
// 		}
// 	}
// 	return result
// }

// func CleanFolder(id int) {
// 	folder := fmt.Sprintf("Round_%03d", id)
// 	os.RemoveAll("./" + folder)
// 	file.EstablishFolder(folder + "/results/gnosis/")
// 	file.EstablishFolder(folder + "/results/mainnet/")
// 	file.EstablishFolder(folder + "/neighbors/gnosis/raw/")
// 	file.EstablishFolder(folder + "/neighbors/gnosis/raw_senders/")
// 	file.EstablishFolder(folder + "/neighbors/gnosis/senders/")
// 	file.EstablishFolder(folder + "/neighbors/mainnet/raw/")
// 	file.EstablishFolder(folder + "/neighbors/mainnet/raw_senders/")
// 	file.EstablishFolder(folder + "/neighbors/mainnet/senders/")
// }
