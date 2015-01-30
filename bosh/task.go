package bosh

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

func retrieveTaskId(resp *http.Response) (taskId int, err error) {
	if resp.StatusCode != 302 {
		err = errors.New("The resp code from toggle request should return 302")
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
