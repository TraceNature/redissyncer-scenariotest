package cases

import (
	"fmt"
	"testcase/global"
)

func DisplayCasesList() {
	fmt.Println("All Cases:")
	for k, v := range CaseTypeMap {
		fmt.Println(k, v)
	}
	logger.Sugar().Info(global.RSPViper.Get("server"))
}
