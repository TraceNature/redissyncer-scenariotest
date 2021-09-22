package compare

import (
	"testcase/global"
	//"testcase/globalzap"
)

var zaplogger = global.RSPLog

type CompareResult struct {
	IsEqual        bool
	NotEqualReason map[string]interface{}
	KeyDiffReason  []interface{}
	KeyType        string
	Key            string
	SourceDB       int //源redis DB number
	TargetDB       int //目标redis DB number
}

func NewCompareResult() CompareResult {
	return CompareResult{
		IsEqual: true,
	}
}

type CompareData interface {
	CompareDB()
	CompareKeys(keys []string)
	CompareString(key string) *CompareResult
	CompareList(key string) *CompareResult
	CompareHash(key string) *CompareResult
	CompareSet(key string) *CompareResult
	CompareZset(key string) *CompareResult
}
