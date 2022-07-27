package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"testcase/cases"
	"testcase/commons"
	"testcase/global"
)

func NewExecCommand() *cobra.Command {
	exec := &cobra.Command{
		Use:   "exec <subcommand>",
		Short: "Execute scenario test",
	}
	exec.AddCommand(NewExecFromFileCommand())
	exec.AddCommand(NewExecFromDirectoryCommand())
	return exec
}

func NewExecFromFileCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "file <test task yaml file path>",
		Short: "execute test from yaml file",
		Run:   execTestCaseFromFileFunc,
	}
	return sc
}

func NewExecFromDirectoryCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "dir <test task dir path>",
		Short: "execute test from dirctory include yml files",
		Run:   execTestCaseFromDirectoryFunc,
	}
	return sc
}

func execTestCaseFromFileFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrln("Please input test yaml file path")
		return
	}

	for _, v := range args {
		//判断文件是否存在
		if !commons.FileExists(v) {
			cmd.PrintErrf("file %s not exists \n", v)
			continue
		}

		//判断文件格式
		yml := strings.HasSuffix(v, ".yml")
		yaml := strings.HasSuffix(v, ".yaml")
		if !yml && !yaml {
			cmd.PrintErrf("file %s not a yml or yaml file \n", v)
			continue
		}

		tc := cases.NewTestCase()
		if err := tc.ParseYamlFile(v); err != nil {
			cmd.PrintErrln(err)
		}
		global.RSPLog.Sugar().Info(tc)

		tc.Exec()
	}
}

func execTestCaseFromDirectoryFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.PrintErrln("Please input directory include yaml files")
		return
	}

	if !commons.IsDir(args[0]) {
		cmd.PrintErrf(" %s not a directory \n", args[0])
		return
	}

	files, err := commons.GetAllFiles(args[0])
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	fmt.Println(files)
	yamlfiles := []string{}
	for _, v := range files {
		//过滤指定格式
		ok := strings.HasSuffix(v, ".yml") || strings.HasSuffix(v, ".yaml")
		if ok {
			yamlfiles = append(yamlfiles, v)
		}
	}

	if len(yamlfiles) == 0 {
		cmd.PrintErrln("No yaml files in the folder!")
		return
	}
	for _, v := range yamlfiles {
		tc := cases.NewTestCase()
		tc.ParseYamlFile(v)
		fmt.Println(tc)
		tc.Exec()
	}

}
