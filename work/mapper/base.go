package mapper

import (
	"errors"
	"github.com/crazy-me/os_scheduler/work/conf"
	"io/ioutil"
	"net/http"
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
