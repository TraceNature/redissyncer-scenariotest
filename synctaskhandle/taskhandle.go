package synctaskhandle

import (
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testcase/global"
	//"testcase/globalzap"
	"testcase/model/response"
	"time"
)

const CreateTaskPath = "/api/v2/createtask"
const StartTaskPath = "/api/v2/starttask"
const StopTaskPath = "/api/v2/stoptask"
const RemoveTaskPath = "/api/v2/removetask"
const ListTasksPath = "/api/v2/listtasks"

//const TaskListByIDs = "/api/task/listbyids"
//const TaskListByName = "/api/task/listbynames"
const LastKeyAcross = "/api/task/lastkeyacross"
const ImportFilePath = "/api/v2/file/createtask"

type Request struct {
	Server string
	Api    string
	Body   string
}

func (r Request) ExecRequest() (result string) {

	client := &http.Client{}
	//req, err := http.NewRequest("POST", r.Server+r.Api, strings.NewReader(r.Body))
	//
	//if err != nil {
	//	logger.Sugar().Error(err)
	//	os.Exit(1)
	//}
	//
	//req.Header.Set("Content-Type", "application/json")
	//resp, err := client.Do(req)

	resp, err := client.Post(r.Server+r.Api, "application/json", strings.NewReader(r.Body))
	global.RSPLog.Info(r.Server + r.Api)

	if err != nil {
		global.RSPLog.Sugar().Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//logger.Sugar().Error(err)
		global.RSPLog.Sugar().Error(err)
		os.Exit(1)
	}

	//var dat map[string]interface{}
	//json.Unmarshal(body, &dat)
	//bodystr, jsonerr := json.MarshalIndent(dat, "", " ")
	//if jsonerr != nil {
	//	logger.Sugar().Error(err)
	//}
	//return string(bodystr)
	return string(body)
}

//创建导入文件任务
func Import(syncserver string, createjson string) []string {
	importreq := &Request{
		Server: syncserver,
		Api:    ImportFilePath,
		Body:   createjson,
	}

	resp := importreq.ExecRequest()

	taskids := gjson.Get(resp, "data").Array()

	if len(taskids) == 0 {
		global.RSPLog.Error("task create faile", zap.Any("response_info", resp))
		os.Exit(1)
	}
	taskidsstrarray := []string{}
	for _, v := range taskids {
		//fmt.Println(gjson.Get(v.String(), "taskId").String())
		taskidsstrarray = append(taskidsstrarray, gjson.Get(v.String(), "taskId").String())
	}

	return taskidsstrarray

}

//创建同步任务
func CreateTask(syncserver string, createjson string) []string {
	createreq := &Request{
		Server: syncserver,
		Api:    CreateTaskPath,
		Body:   createjson,
	}

	fmt.Println(createjson)
	resp := createreq.ExecRequest()
	taskids := gjson.Get(resp, "data").Array()
	fmt.Println(resp)
	if len(taskids) == 0 {
		global.RSPLog.Sugar().Error(errors.New("task create faile \n"), resp)
		//logger.Sugar().Info(resp)
		os.Exit(1)
	}
	taskidsstrarray := []string{}
	for _, v := range taskids {
		//fmt.Println(gjson.Get(v.String(), "taskId").String())
		taskidsstrarray = append(taskidsstrarray, gjson.Get(v.String(), "taskId").String())
	}

	return taskidsstrarray

}

//Start task
func StartTask(syncserver string, taskid string) {
	jsonmap := make(map[string]interface{})
	jsonmap["taskid"] = taskid
	startjson, err := json.Marshal(jsonmap)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		os.Exit(1)
	}
	startreq := &Request{
		Server: syncserver,
		Api:    StartTaskPath,
		Body:   string(startjson),
	}
	startreq.ExecRequest()

}

//Stop task by task ids
func StopTaskByIds(syncserver string, taskId string) {
	jsonmap := make(map[string]interface{})

	jsonmap["taskId"] = taskId
	stopjsonStr, err := json.Marshal(jsonmap)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		os.Exit(1)
	}
	stopreq := &Request{
		Server: syncserver,
		Api:    StopTaskPath,
		Body:   string(stopjsonStr),
	}
	response := stopreq.ExecRequest()
	fmt.Println(response)

}

//Remove task by name
func RemoveTaskByName(syncserver string, taskname string) {

	jsonmap := make(map[string]interface{})

	taskids, err := GetSameTaskNameIDs(syncserver, taskname)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		os.Exit(1)
	}

	if len(taskids) == 0 {
		return
	}

	for _, id := range taskids {
		jsonmap["taskId"] = id
		stopjsonStr, err := json.Marshal(jsonmap)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			os.Exit(1)
		}
		stopreq := &Request{
			Server: syncserver,
			Api:    StopTaskPath,
			Body:   string(stopjsonStr),
		}
		stopResult := stopreq.ExecRequest()
		fmt.Println(stopResult)

		time.Sleep(10 * time.Second)

		removereq := &Request{
			Server: syncserver,
			Api:    RemoveTaskPath,
			Body:   string(stopjsonStr),
		}

		removeResult := removereq.ExecRequest()
		fmt.Println(removeResult)
	}

}

//获取同步任务状态
func GetTaskStatus(syncserver string, ids []string) (map[string]string, error) {
	jsonmap := make(map[string]interface{})

	jsonmap["regulation"] = "byids"
	//jsonmap["taskIDs"] = ids

	listtaskjsonStr, err := json.Marshal(jsonmap)
	if err != nil {
		return nil, err
	}
	listreq := &Request{
		Server: syncserver,
		Api:    ListTasksPath,
		Body:   string(listtaskjsonStr),
	}
	listresp := listreq.ExecRequest()

	taskarray := gjson.Get(listresp, "result").Array()

	if len(taskarray) == 0 {
		return nil, errors.New("No status return")
	}

	statusmap := make(map[string]string)

	for _, v := range taskarray {
		id := gjson.Get(v.String(), "taskId").String()
		status := v.String()
		statusmap[id] = status
	}

	return statusmap, nil
}

// @title    GetSameTaskNameIDs
// @description   获取同名任务列表
// @auth      Jsw             时间（2020/7/1   10:57 ）
// @param     syncserver        string         "redissyncer ip:port"
// @param    taskname        string         "任务名称"
// @return    taskids        []string         "任务id数组"
func GetSameTaskNameIDs(syncserver string, taskname string) ([]string, error) {

	existstaskids := []string{}
	listjsonmap := make(map[string]interface{})
	listjsonmap["regulation"] = "bynames"
	listjsonmap["taskNames"] = strings.Split(taskname, ",")
	listjsonStr, err := json.Marshal(listjsonmap)

	if err != nil {
		global.RSPLog.Sugar().Info(err.Error())
		return nil, err
	}
	listtaskreq := &Request{
		Server: syncserver,
		Api:    ListTasksPath,
		Body:   string(listjsonStr),
	}

	listresp := listtaskreq.ExecRequest()

	tasklist := gjson.Get(listresp, "result").Array()

	if len(tasklist) > 0 {
		for _, v := range tasklist {
			existstaskids = append(existstaskids, gjson.Get(v.String(), "taskStatus.taskId").String())
		}
	}

	return existstaskids, nil
}

func GetLastKeyAcross(syncserver string, taskID string) (response.LastKeyAcrossResult, error) {
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	var result response.LastKeyAcrossResult
	reqJson := make(map[string]interface{})
	reqJson["taskId"] = taskID
	jsonStr, err := json.Marshal(reqJson)
	if err != nil {
		global.RSPLog.Info(err.Error())
		return result, err
	}
	req := &Request{
		Server: syncserver,
		Api:    LastKeyAcross,
		Body:   string(jsonStr),
	}

	resp := req.ExecRequest()

	//jsonstr, err := jsoniter.Marshal(resp)
	//if err != nil {
	//	return result, nil
	//}

	if err := jsoniter.Unmarshal([]byte(resp), &result); err != nil {
		return result, err
	}

	return result, nil
}
