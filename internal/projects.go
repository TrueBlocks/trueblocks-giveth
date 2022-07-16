package internal

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
	"github.com/spf13/cobra"
)

// RunProjects runs the rounds command
func RunProjects(cmd *cobra.Command, args []string) error {
	projects, format, err := getProjectsOptions(cmd, args)
	if err != nil {
		return err
	}

	for i := 0; i < len(projects); i++ {
		markCore(&projects[i])
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].WalletAddress < projects[j].WalletAddress
	})

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

// getProjectsOptions processes command line options for the Rounds command
func getProjectsOptions(cmd *cobra.Command, args []string) (rounds []data.Project, format string, err error) {
	format, err = cmd.Flags().GetString("fmt")
	if err != nil {
		log.Fatal(err)
	}

	return data.GetProjects(), format, err
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
