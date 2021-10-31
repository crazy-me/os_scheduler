package mapper

import (
	"encoding/json"
	"errors"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Base struct {
}

func (b *Base) HttpRequestByPost(body []byte) (err error) {
	resp, err := http.Post(conf.C.AgentEndpoint,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(body)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	if string(result) != "success" {
		err = errors.New(string(result))
	}
	return
}

// ParserJob 解析调度结果
func (b *Base) ParserJob(original []byte, jobRes *entity.JobResultDto) error {
	// 解析调度结果
	if err := json.Unmarshal(original, jobRes); err != nil {
		return err
	}

	// 调度结果是否成功
	if 1 != jobRes.Code {
		return errors.New("task scheduler failed")
	}
	return nil

}

// FilterIndex 过滤不需要的指标
func (b *Base) FilterIndex(m map[string]string, s []string) {
	for _, v := range s {
		delete(m, v)
	}
}

// ConvertToGb 将Byte转换为GB
func (b *Base) ConvertToGb(byteString string) string {
	byteValue, err := strconv.Atoi(byteString)
	if err != nil || 0 == byteValue {
		return "0"
	}

	return strconv.Itoa(byteValue / 1024 / 1024 / 1024)
}

// ConvertToMb 将Byte转换为MB
func (b *Base) ConvertToMb(byteString string) string {
	byteValue, err := strconv.Atoi(byteString)
	if err != nil || 0 == byteValue {
		return "0"
	}

	return strconv.Itoa(byteValue / 1024 / 1024)
}

// NoCalculationTsRequired 过滤不需要计算的指标
func (b *Base) NoCalculationTsRequired(queue string) bool {
	for _, v := range conf.C.NoCalculation {
		if queue == v { // 如果存在测返回false
			return false
		}
	}
	return true
}
