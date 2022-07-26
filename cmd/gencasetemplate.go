package cmd

import (
	"github.com/spf13/cobra"
)

func NewGenCaseTemplateCmd() *cobra.Command {
	gc := &cobra.Command{
		Use:   "gencase",
		Short: "Generate test case description file by case type",
		Run:   genDataCommandFunc,
	}
	gc.AddCommand(GenCaseSingle2Single())
	gc.AddCommand(GenCaseSingle2SingleWithDBMap())
	gc.AddCommand(GenCaseSingle2Cluster())
	gc.AddCommand(GenCaseCluster2Cluster())
	gc.AddCommand(GenCaseImportRdb2Single())
	gc.AddCommand(GenCaseImportAof2Single())
	gc.AddCommand(GenCaseImportRdb2Cluster())
	gc.AddCommand(GenCaseImportAof2Cluster())
	return gc
}

func GenCaseSingle2Single() *cobra.Command {
	gd := &cobra.Command{
		Use:   "single2single",
		Short: "Generate single2single test case description file",
		Run:   genCaseSingle2SingleFunc,
	}
	return gd
}

func genCaseSingle2SingleFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		//cmd.PrintErrln("Please input test yaml file path")
		//return
	}

	cmd.Println("genCaseSingle2SingleFunc")

}

func GenCaseSingle2SingleWithDBMap() *cobra.Command {
	gd := &cobra.Command{
		Use:   "single2singlewithdbmap",
		Short: "Generate single2singlewithdbmap test case description file",
		Run:   genCaseSingle2SingleWithDBMapFunc,
	}

	return gd
}

func genCaseSingle2SingleWithDBMapFunc(cmd *cobra.Command, args []string) {

}

func GenCaseSingle2Cluster() *cobra.Command {
	gd := &cobra.Command{
		Use:   "single2cluster",
		Short: "Generate single2cluster test case description file",
		Run:   genCaseSingle2ClusterFunc,
	}

	return gd
}

func genCaseSingle2ClusterFunc(cmd *cobra.Command, args []string) {

}

func GenCaseCluster2Cluster() *cobra.Command {
	gd := &cobra.Command{
		Use:   "cluster2cluster",
		Short: "Generate cluster2cluster test case description file",
		Run:   genCaseCluster2ClusterFunc,
	}
	return gd
}

func genCaseCluster2ClusterFunc(cmd *cobra.Command, args []string) {

}
func GenCaseImportRdb2Single() *cobra.Command {
	gd := &cobra.Command{
		Use:   "importrdb2single",
		Short: "Generate single2single test case description file",
		Run:   genCaseImportRdb2SingleFunc,
	}
	return gd
}

func genCaseImportRdb2SingleFunc(cmd *cobra.Command, args []string) {

}

func GenCaseImportAof2Single() *cobra.Command {
	gd := &cobra.Command{
		Use:   "importaof2single",
		Short: "Generate single2single test case description file",
		Run:   genCaseImportAof2SingleFunc,
	}

	return gd
}

func genCaseImportAof2SingleFunc(cmd *cobra.Command, args []string) {

}

func GenCaseImportRdb2Cluster() *cobra.Command {
	gd := &cobra.Command{
		Use:   "importrdb2cluster",
		Short: "Generate single2single test case description file",
		Run:   genCaseImportRdb2ClusterFunc,
	}

	return gd
}
func genCaseImportRdb2ClusterFunc(cmd *cobra.Command, args []string) {

}
func GenCaseImportAof2Cluster() *cobra.Command {
	gd := &cobra.Command{
		Use:   "importaof2cluster",
		Short: "Generate single2single test case description file",
		Run:   genCaseImportAof2ClusterFunc,
	}

	return gd
}

func genCaseImportAof2ClusterFunc(cmd *cobra.Command, args []string) {

}
