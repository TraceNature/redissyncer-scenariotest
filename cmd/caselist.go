package cmd

import (
	"github.com/spf13/cobra"
	"testcase/cases"
)

// NewConfigCommand return a config subcommand of rootCmd
func NewCaseListCommand() *cobra.Command {
	cl := &cobra.Command{
		Use:   "caselist ",
		Short: "list all cases",
		Run: caseListCommandFunc,
	}
	return cl
}



func caseListCommandFunc(cmd *cobra.Command, args []string) {
	cases.DisplayCasesList()
}

