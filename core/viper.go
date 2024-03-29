

package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testcase/commons"
	"testcase/global"
)

func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		//默认config文件查找路径 ./ -> 执行文件路径
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(ConfigEnv); configEnv == "" {
				//获取可执行文件的绝对路径
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

				if commons.FileExists(ConfigFile) {
					config = ConfigFile
				}

				if commons.FileExists(dir + "/" + ConfigFile) {
					config = dir + "/" + ConfigFile
				}

				//fmt.Printf("您正在使用config的默认值,config的路径为%v\n", config)
			} else {
				config = configEnv
				//fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", config)
			}
		} else {
			//fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		}
	} else {
		config = path[0]
		//fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.RSPConfig); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&global.RSPConfig); err != nil {
		fmt.Println(err)
	}

	return v
}
