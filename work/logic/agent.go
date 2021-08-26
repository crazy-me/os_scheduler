package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Agent struct {
	JobInfo *entity.JobExecuteResult
	Func    func() (err error)
}

func (agent *Agent) PushNetwork() (err error) {
	var networkData entity.Network
	err = json.Unmarshal(agent.JobInfo.Output, &networkData)
	if err != nil {
		return
	}

	if networkData.Code != 1 {
		err = errors.New("command " + agent.JobInfo.ExecStatus.Job.JobCommand +
			"code: " + strconv.Itoa(networkData.Code) + "msg: " + networkData.Msg)
		return
	}
	var agentPush entity.AgentPush
	pushData := make([]entity.AgentPush, 0)
	networkType := reflect.TypeOf(networkData.Data)
	networkValue := reflect.ValueOf(networkData.Data)
	for i := 0; i < networkType.NumField(); i++ {
		//fmt.Printf("字段:%s -- 值:%v \n", networkType.Field(i).Tag.Get("json"), networkValue.Field(i).Interface())
		agentPush.Ident = "192.168.31.13"
		agentPush.Alias = networkType.Field(i).Tag.Get("json")
		agentPush.Metric = networkType.Field(i).Tag.Get("json")
		agentPush.Time = time.Now().Unix()
		if v, ok := networkValue.Field(i).Interface().(string); ok {
			agentPush.Value = v
		}
		pushData = append(pushData, agentPush)
	}
	body, err := json.Marshal(pushData)
	if err != nil {
		err = errors.New("command " + agent.JobInfo.ExecStatus.Job.JobCommand + "PushNetwork err: " + err.Error())
		return
	}
	err = agent.HttpRequestByPost(body)
	return
}

func (agent *Agent) PushServer() (err error) {
	return
}

func (agent *Agent) PushMysql() (err error) {
	return
}

func (agent *Agent) PushApply() (err error) {
	return
}

func (agent *Agent) HttpRequestByPost(body []byte) (err error) {
	resp, err := http.Post(conf.C.AgentEndpoint,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(body)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))
	if string(result) != "success" {
		err = errors.New(string(result))
	}
	return
}

func (agent *Agent) PushEmpty() (err error) {
	return
}
