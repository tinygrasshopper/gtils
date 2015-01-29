package bosh

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

func retrieveManifest(response *http.Response) (manifest io.Reader, err error) {
	if response.StatusCode != 200 {
		err = errors.New("The retriveing bosh manifest API response code is not equal to 200")
		return
	}
	m := make(map[string]interface{})
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(body, &m)
	if err != nil {
		return
	}
	data, err := yaml.Marshal(m["manifest"])
	if err != nil {
		return
	}
	return bytes.NewReader(data), nil
}
