package cmd

import (
	"github.com/TrueBlocks/trueblocks-giveth/internal"
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Produces data related to the projects including addresses.tsv",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: internal.RunProjects,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	projectsCmd.Flags().BoolP("categories", "a", false, "show only the categories for each project")
	projectsCmd.Flags().SortFlags = false

	rootCmd.AddCommand(projectsCmd)
}
