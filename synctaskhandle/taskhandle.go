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
const HealthCheck = "/health"

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
	resp, err := client.Post(r.Server+r.Api, "application/json", strings.NewReader(r.Body))

	if err != nil {
		global.RSPLog.Sugar().Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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

// redissyncer-server 健康检查
func SyncerServerAlive(syncserver string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, _ := client.Get(syncserver + HealthCheck)
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

//清理 redissyncer server 上的所有任务，避免测试数据错误
func SyncerServerClean(syncserver string) error {
	jsonmap := make(map[string]interface{})
	jsonmap["regulation"] = "all"
	listtaskjson, err := json.Marshal(jsonmap)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}
	// 获取 server上所有的任务id
	alltasksreq := &Request{
		Server: syncserver,
		Api:    ListTasksPath,
		Body:   string(listtaskjson),
	}
	resp := alltasksreq.ExecRequest()
	tasks := gjson.Get(resp, "data").Array()
	if len(tasks) == 0 {
		return nil
	}
	taskids := []string{}
	for _, v := range tasks {
		taskid := gjson.Get(v.String(), "taskId").String()
		taskids = append(taskids, taskid)
	}

	// 停止 所有任务
	stopjsonmap := make(map[string]interface{})
	stopjsonmap["taskids"] = taskids
	stopjson, err := json.Marshal(stopjsonmap)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}
	stoptasksreq := &Request{
		Server: syncserver,
		Api:    StopTaskPath,
		Body:   string(stopjson),
	}
	stoptasksreq.ExecRequest()
	// 删除所有任务

	removejsonmap := make(map[string]interface{})
	removejsonmap["taskids"] = taskids
	removejson, err := json.Marshal(stopjsonmap)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}
	removetasksreq := &Request{
		Server: syncserver,
		Api:    RemoveTaskPath,
		Body:   string(removejson),
	}

	resp = removetasksreq.ExecRequest()

	global.RSPLog.Sugar().Debug(resp)

	return nil
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
		os.Exit(1)
	}
	taskidsstrarray := []string{}
	for _, v := range taskids {
		taskidsstrarray = append(taskidsstrarray, gjson.Get(v.String(), "taskId").String())
	}

	return taskidsstrarray

}

//Start task
func StartTask(syncserver string, taskid string) string {
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
	r := startreq.ExecRequest()
	return r

}

//Stop task by task ids
func StopTaskByIds(syncserver string, taskId string) string {
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
	return response
}

//Remove task by name
func RemoveTaskByName(syncserver string, taskname string) {

	jsonmap := make(map[string]interface{})

	taskids, err := GetSameTaskNameIDs(syncserver, taskname)

	global.RSPLog.Sugar().Info("taskids:", taskids)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		os.Exit(1)
	}

	if len(taskids) == 0 {
		global.RSPLog.Sugar().Error("Taskids array is empty")
		return
	}

	for _, id := range taskids {
		jsonmap["taskids"] = []string{id}
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
		global.RSPLog.Sugar().Info("stop task ", id)
		stopResult := stopreq.ExecRequest()
		global.RSPLog.Sugar().Info(stopResult)

		removereq := &Request{
			Server: syncserver,
			Api:    RemoveTaskPath,
			Body:   string(stopjsonStr),
		}

		times := 0
		for {
			if times >= 3 {
				break
			}
			removeResult := removereq.ExecRequest()
			code := gjson.Get(removeResult, "code").String()
			global.RSPLog.Sugar().Info(code)
			if "2000" == code {
				break
			}
			if code == "101" {
				time.Sleep(10 * time.Second)
			}
			times++
		}
	}
}

//获取同步任务状态
func GetTaskStatus(syncserver string, ids []string) (map[string]string, error) {
	jsonmap := make(map[string]interface{})

	jsonmap["regulation"] = "byids"
	jsonmap["taskids"] = ids

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

	taskarray := gjson.Get(listresp, "data").Array()

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
	listjsonmap["tasknames"] = strings.Split(taskname, ",")
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

	tasklist := gjson.Get(listresp, "data").Array()

	if len(tasklist) > 0 {
		for _, v := range tasklist {
			existstaskids = append(existstaskids, gjson.Get(v.String(), "taskId").String())
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
