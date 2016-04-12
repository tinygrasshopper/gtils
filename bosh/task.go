package bosh

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

//Task -
type Task struct {
	Id          int    `json:"id"`
	State       string `json:"state"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

const (
	//ERROR -
	ERROR int = 1
	//PROCESSING -
	PROCESSING int = 2
	//DONE -
	DONE int = 3
	//QUEUED -
	QUEUED int = 4
)

//TASKRESULT -
var TASKRESULT = map[string]int{"error": ERROR, "processing": PROCESSING, "done": DONE, "queued": QUEUED}

func retrieveTaskID(resp *http.Response) (taskID int, err error) {
	if resp.StatusCode != 302 {
		err = ErrorTaskRedirectStatusCode
		return
	}
	redirectUrls := resp.Header["Location"]
	if redirectUrls == nil || len(redirectUrls) < 1 {
		err = errors.New("Could not find redirect url for bosh tasks")
		return
	}
	regex := regexp.MustCompile(`^.*tasks/`)
	idString := regex.ReplaceAllString(redirectUrls[0], "")
	return strconv.Atoi(idString)
}

func retrieveTaskStatus(resp *http.Response) (task *Task, err error) {
	if resp.StatusCode != 200 {
		err = ErrorTaskStatusCode
		return
	}
	task = &Task{}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, task)
	if err != nil {
		return
	}
	return

}
