package utils

import (
	"encoding/json"
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func UnpackJob(value []byte) (ret *entity.Job, err error) {
	var job entity.Job
	if err = json.Unmarshal(value, &job); err != nil {
		return
	}
	ret = &job
	return
}

func ExtractJobKey(jobKey string) string {
	return strings.TrimPrefix(jobKey, constants.JOB_SAVE_DIR)
}
