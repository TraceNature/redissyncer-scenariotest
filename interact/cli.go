package interact

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"testcase/commons"
	"testcase/core"
	"testcase/global"

	//"interactioncli/check"
	"testcase/cmd"
)

type CommandFlags struct {
	URL      string
	CAPath   string
	CertPath string
	KeyPath  string
	Help     bool
}

var (
	commandFlags = CommandFlags{}
	cfgFile      string
	//detach          bool
	//syncserver      string
	//Confignotseterr error
	interact bool
	//version  bool
)

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

var query = ""

var readLineCompleter *readline.PrefixCompleter

func init() {
	cobra.EnablePrefixMatching = true
	cobra.OnInitialize(initConfig)
}

func cliRun(cmd *cobra.Command, args []string) {
	banner := "\n               _ _                                  _____ \n  _ __ ___  __| (_)___ ___ _   _ _ __   ___ ___ _ _|_   _|\n | '__/ _ \\/ _` | / __/ __| | | | '_ \\ / __/ _ \\ '__|| |  \n | | |  __/ (_| | \\__ \\__ \\ |_| | | | | (_|  __/ |   | |  \n |_|  \\___|\\__,_|_|___/___/\\__, |_| |_|\\___\\___|_|   |_|  \n                           |___/                          \n"
	if interact {
		//err := check.CheckEnv()
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}

		cmd.Println(banner)
		cmd.Println("Input 'help;' for usage. \nCommand must end with ';'. \n'tab' for command complete.\n^C or exit to quit.")
		loop()
		return
	}

	if len(args) == 0 {
		cmd.Help()
		return
	}

}

func getBasicCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "redissyncer-test",
		Short: "redissyncer command line interface",
		Long:  "",
	}

	rootCmd.PersistentFlags().BoolVarP(&commandFlags.Help, "help", "h", false, "help message")

	rootCmd.AddCommand(
		cmd.NewConfigCommand(),
		cmd.NewExecCommand(),
		cmd.NewGenDataCommand(),
		cmd.NewCaseListCommand(),
	)

	rootCmd.Flags().ParseErrorsWhitelist.UnknownFlags = true
	rootCmd.SilenceErrors = true
	return rootCmd
}

func getInteractCmd(args []string) *cobra.Command {
	rootCmd := getBasicCmd()
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
	}

	rootCmd.SetArgs(args)
	rootCmd.ParseFlags(args)
	rootCmd.SetOut(os.Stdout)
	hiddenFlag(rootCmd)

	return rootCmd
}

func getMainCmd(args []string) *cobra.Command {
	rootCmd := getBasicCmd()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&syncserver, "syncserver", "s", "", "sync server address")
	//rootCmd.Flags().BoolVarP(&detach, "detach", "d", true, "Run pdctl without readline.")
	rootCmd.Flags().BoolVarP(&interact, "interact", "i", false, "Run pdctl with readline.")
	//rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and exit.")

	rootCmd.Run = cliRun

	rootCmd.SetArgs(args)
	rootCmd.ParseFlags(args)
	rootCmd.SetOut(os.Stdout)

	readLineCompleter = readline.NewPrefixCompleter(genCompleter(rootCmd)...)
	return rootCmd
}

// Hide the flags in help and usage messages.
func hiddenFlag(cmd *cobra.Command) {
	cmd.LocalFlags().MarkHidden("pd")
	cmd.LocalFlags().MarkHidden("cacert")
	cmd.LocalFlags().MarkHidden("cert")
	cmd.LocalFlags().MarkHidden("key")
}

// MainStart start main command
func MainStart(args []string) {
	startCmd(getMainCmd, args)
}

// Start start interact command
func Start(args []string) {
	startCmd(getInteractCmd, args)
}

func startCmd(getCmd func([]string) *cobra.Command, args []string) {
	rootCmd := getCmd(args)

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if global.RSPViper != nil {
		return
	}
	if cfgFile != "" {
		if !commons.FileExists(cfgFile) {
			panic("config file not exists")
		}

		global.RSPViper = core.Viper(cfgFile)
		global.RSPLog = core.Zap()
		return
	}
	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

}

func loop() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "redissyncer-test> ",
		HistoryFile:            "/tmp/readline.tmp",
		AutoComplete:           readLineCompleter,
		DisableAutoSaveHistory: true,
		InterruptPrompt:        "^C",
		EOFPrompt:              "^D",
		HistorySearchFold:      true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var cmds []string

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				break
			} else if err == io.EOF {
				break
			}
			continue
		}
		if line == "exit" {
			os.Exit(0)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)

		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt("... ")
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt("redissyncer-test> ")
		rl.SaveHistory(cmd)

		args, err := shellwords.Parse(cmd)
		if err != nil {
			fmt.Printf("parse command err: %v\n", err)
			continue
		}
		Start(args)
	}
}

func genCompleter(cmd *cobra.Command) []readline.PrefixCompleterInterface {
	pc := []readline.PrefixCompleterInterface{}

	for _, v := range cmd.Commands() {
		if v.HasFlags() {
			flagsPc := []readline.PrefixCompleterInterface{}
			flagUsages := strings.Split(strings.Trim(v.Flags().FlagUsages(), " "), "\n")
			for i := 0; i < len(flagUsages)-1; i++ {
				flagsPc = append(flagsPc, readline.PcItem(strings.Split(strings.Trim(flagUsages[i], " "), " ")[0]))
			}
			flagsPc = append(flagsPc, genCompleter(v)...)
			pc = append(pc, readline.PcItem(strings.Split(v.Use, " ")[0], flagsPc...))
		} else {
			pc = append(pc, readline.PcItem(strings.Split(v.Use, " ")[0], genCompleter(v)...))
		}
	}
	return pc
}
