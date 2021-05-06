/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	//"encoding/json"
	//"github.com/tidwall/gjson"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"testcase/core"
	"testcase/global"
	"testcase/globalzap"
	"testcase/synctaskhandle"

	//"io/ioutil"
	"os"
	//"testcase/synctaskhandle"
)

var cfgFile string
var logger *zap.Logger
var Viper *viper.Viper

func init() {
	logger = globalzap.GetLogger()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "testcase",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		execfile := "./tasks/listtasks.json"

		jsonFile, err := os.Open(execfile)
		defer jsonFile.Close()

		if err != nil {
			log.Println(err)

			os.Exit(1)
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)

		req := &synctaskhandle.Request{
			Server: viper.GetViper().GetString("syncserver"),
			Api:    synctaskhandle.ListTasksPath,
			Body:   string(byteValue),
		}
		resp := req.ExecRequest()
		fmt.Println(viper.GetViper().GetString("syncserver"))
		fmt.Println(resp)

		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yml)")
	rootCmd.MarkFlagRequired("config")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		//core.RSPViper.SetConfigFile(cfgFile)
		global.RSPViper = core.Viper(cfgFile)
	} else {
		// Find home directory.
		//core.RSPViper.AddConfigPath(".")
		//core.RSPViper.SetConfigName("config.yml")
		global.RSPViper = core.Viper()
	}

	global.RSPViper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := global.RSPViper.ReadInConfig(); err == nil {
		logger.Sugar().Info("Using config file:", viper.ConfigFileUsed())
		if global.RSPViper.GetBool("logjsonformat") {
			logger.Sugar().Info(&logrus.JSONFormatter{})
		}
	} else {
		logger.Sugar().Error(err)
		os.Exit(1)
	}

}
