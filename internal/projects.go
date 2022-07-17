package internal

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
	"github.com/spf13/cobra"
)

// RunProjects runs the rounds command
func RunProjects(cmd *cobra.Command, args []string) error {
	projects, remote, categories, format, verbose, err := getProjectsOptions(cmd, args)
	if err != nil {
		return err
	}

	if remote {
		var url string = `curl --location --request POST 'https://mainnet.serve.giveth.io/graphql' --header 'Content-Type: application/json' --data-raw '{"query":"{projectById(id:%d) {id, title, balance, image, slug, slugHistory, creationDate, updatedAt, admin, description, walletAddress, impactLocation, qualityScore, verified, traceCampaignId, listed, givingBlocksId, status { id, symbol, name, description }, categories { name }, reaction { id }, adminUser { id, email, firstName, walletAddress }, organization { name, label, supportCustomTokens }, addresses {address, isRecipient, networkId }, totalReactions, totalDonations, totalTraceDonations}}","variables":{}}' | jq >data/raw/%05d.json; sleep 4`
		for i := 1; i < int(float64(len(projects))*1.3); i++ {
			fmt.Printf(strings.Replace(url, "\n", " ", -1)+"\n", i, i)
		}
		return nil
	}

	if categories {
		cats := data.GetCategories()
		sorted := []data.CategoryCounter{}
		for key, values := range cats {
			sorted = append(sorted, data.NewCategoryCounter(key, len(values)))
		}
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Count == sorted[j].Count {
				return sorted[i].Key < sorted[j].Key
			}
			return sorted[i].Count < sorted[j].Count
		})

		for _, s := range sorted {
			values := cats[s.Key]
			fmt.Println(len(values), properTitle(s.Key))
			if verbose {
				for _, v := range values {
					fmt.Println("\t", v)

				}
			}
		}

		return nil
	}

	if format == "txt" || format == "csv" {
		unused := data.SimpleProject{}
		output.Header(unused, os.Stdout, format)
		defer output.Footer(unused, os.Stdout, format)
		for i, project := range projects {
			if project.Id == "1389" {
				fmt.Println(project)
				fmt.Println()
			}
			output.Line(data.ToSimpleProject(&project), os.Stdout, format, i == 0)
		}

	} else {
		unused := data.Project{}
		output.Header(unused, os.Stdout, format)
		defer output.Footer(unused, os.Stdout, format)
		for i, project := range projects {
			output.Line(project, os.Stdout, format, i == 0)
		}
	}

	return nil
}

func markCore(p *data.Project) {
	cores := make(map[string]bool, 4)
	cores["0x4d9339dd97db55e3b9bcbe65de39ff9c04d1c2cd"] = true
	cores["0x900db999074d9277c5da2a43f252d74366230da0"] = true
	cores["0xecb179ea5910d652eda6988e919c7930f5ffcf11"] = true
	cores["0xf924ff0f192f0c7c073161e0d62ce7635114e74f"] = true
	if cores[p.WalletAddress] {
		id, _ := strconv.Atoi(p.Id)
		p.Id = fmt.Sprintf("core-%04d", id)
	}
}

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

// getProjectsOptions processes command line options for the Rounds command
func getProjectsOptions(cmd *cobra.Command, args []string) (rounds []data.Project, remote bool, categories bool, format string, verbose bool, err error) {
	format, err = cmd.Flags().GetString("fmt")
	if err != nil {
		log.Fatal(err)
	}
	if format == "" {
		format = "txt"
	}

	remote, err = cmd.Flags().GetBool("remote")
	if err != nil {
		log.Fatal(err)
	}
	categories, err = cmd.Flags().GetBool("categories")
	if err != nil {
		log.Fatal(err)
	}

	verbose, err = cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Fatal(err)
	}

	projects := data.GetProjects()
	for i := 0; i < len(projects); i++ {
		markCore(&projects[i])
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].WalletAddress < projects[j].WalletAddress
	})

	return projects, remote, categories, format, verbose, err
}
