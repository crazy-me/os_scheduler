package mapper

import (
	"encoding/json"
	"errors"
	"github.com/crazy-me/os_scheduler/common/entity"
	"reflect"
	"strconv"
	"time"
)

type Network struct {
	Base
}

func (network *Network) Push(r *entity.JobExecuteResult) (err error) {
	var (
		networkDTO *entity.Network
		agentPush  entity.AgentPush
	)
	err = json.Unmarshal(r.Output, &networkDTO)
	if err != nil {
		return
	}

	// 调度结果是否成功
	if networkDTO.Code != 1 {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand +
			"code: " + strconv.Itoa(networkDTO.Code) + "msg: " + networkDTO.Msg)
		return
	}
	pushData := make([]entity.AgentPush, 0)

	t := reflect.TypeOf(networkDTO.Data)
	v := reflect.ValueOf(networkDTO.Data)
	for i := 0; i < t.NumField(); i++ {
		agentPush.Ident = r.ExecStatus.Job.JobIdent
		agentPush.Alias = t.Field(i).Tag.Get("json")
		agentPush.Metric = t.Field(i).Tag.Get("json")
		agentPush.Time = time.Now().Unix()
		if val, ok := v.Field(i).Interface().(string); ok {
			agentPush.Value = val
		}
		pushData = append(pushData, agentPush)
	}
	body, err := json.Marshal(pushData)
	if err != nil {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand + "PushNetwork err: " + err.Error())
		return err
	}
	err = network.HttpRequestByPost(body)
	return
}
