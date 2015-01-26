package bosh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"gopkg.in/yaml.v2"

	. "github.com/pivotalservices/gtils/http"
)

type API struct {
	Path           string
	ContentType    string
	Method         string
	HandleResponse HandleRespFunc
}

type regexReplace func(string) string

func replaceUrlParams(params map[string]string) regexReplace {
	return func(str string) string {
		var retValue string
		for key, value := range params {
			matchString := fmt.Sprintf("{%s}", key)
			if str == matchString {
				retValue = value
			}
		}
		return retValue
	}
}

func (api *API) execute(ip string, port int, username, password string, body io.Reader, pathParams, queryParams map[string]string) (ret interface{}, err error) {
	endpoint, err := api.GetUrl(ip, port, pathParams, queryParams)
	if err != nil {
		return
	}
	gateway := NewHttpGateway(endpoint, username, password, api.ContentType, api.HandleResponse, body)
	return gateway.Execute(api.Method)
}

func (api *API) GetUrl(ip string, port int, pathParams, queryParams map[string]string) (urlString string, err error) {
	host := fmt.Sprintf("https://%s:%d/%s", ip, port, api.Path)
	urlString = regexp.MustCompile("{(.+?)}").ReplaceAllStringFunc(host, replaceUrlParams(pathParams))
	u, err := url.Parse(urlString)
	if err != nil {
		return
	}
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func retrieveManifest(response *http.Response) (resp interface{}, err error) {
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

var retrieveManifestAPI API = API{
	Path:           "deployments/{deployment}",
	Method:         "GET",
	HandleResponse: retrieveManifest,
}
