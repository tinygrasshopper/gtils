package bosh

import (
	"io"

	"github.com/pivotalservices/gtils/http"
)

type Bosh interface {
	GetDeploymentManifest(deploymentName string) (io.Reader, error)
}

type BoshDirector struct {
	ip       string
	port     int
	username string
	password string
}

func NewBoshDirector(ip, username, password string, port int) *BoshDirector {
	return &BoshDirector{
		ip:       ip,
		port:     port,
		username: username,
		password: password,
	}
}

var NewBoshGateway = func(director *BoshDirector, api *API, body io.Reader, pathParams, queryParams map[string]string) (gateway http.HttpGateway, err error) {
	endpoint, err := api.GetUrl(director.ip, director.port, pathParams, queryParams)
	if err != nil {
		return
	}
	return http.NewHttpGateway(endpoint, director.username, director.password, api.ContentType, api.HandleResponse, body), nil
}

func (director *BoshDirector) GetDeploymentManifest(deploymentName string) (manifest io.Reader, err error) {
	pathParams := map[string]string{"deploymentName": deploymentName}
	m, err := retrieveManifestAPI.execute(director.ip, director.port, director.username, director.password, nil, pathParams, nil)
	if err != nil {
		return
	}
	return m.(io.Reader), nil
}
