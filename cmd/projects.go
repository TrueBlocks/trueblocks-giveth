package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Produces a listing of the projects included in the addresses list",
	Long:  `This tool produces a list and manages the projects list.`,
	RunE:  internal.RunProjects,
}

func init() {
	// projectsCmd.PersistentFlags().String("foo", "", "A help for foo")
	projectsCmd.Flags().BoolP("categories", "a", false, "show only the categories for each project")
	projectsCmd.Flags().SortFlags = false

	rootCmd.AddCommand(projectsCmd)
}
