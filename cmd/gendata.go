package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
	"testcase/commons"
	"testcase/generatedata"
)

func NewGenDataCommand() *cobra.Command {
	gd := &cobra.Command{
		Use:   "gendata <datadescraption.yml>",
		Short: "Generate basic data through yaml description file",
		Run:   genDataCommandFunc,
	}
	//gd.AddCommand(NewBaseDataCommand())
	return gd
}

func genDataCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrln("Please input test yaml file path")
		return
	}

	for _, v := range args {
		//cmd.Println(v)
		//判断文件是否存在
		if !commons.FileExists(v) {
			cmd.PrintErrf("file %s not exists \n", v)
			continue
		}

		//判断文件格式
		isYaml := strings.HasSuffix(v, ".yml") || strings.HasSuffix(v, ".yaml")
		if !isYaml {
			cmd.PrintErrf("file %s not a yml or yaml file \n", v)
			continue
		}

		file, err := os.Open(v)
		if err != nil {
			cmd.PrintErrln(err)
		}
		dec := yaml.NewDecoder(file)

		for {
			data := generatedata.GenData{}

			err := dec.Decode(&data)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					cmd.PrintErrln(err)
				}
				break
			}

			data.Exec()

		}
	}
}
