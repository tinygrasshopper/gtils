package bosh

import (
	"fmt"
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
	gateway  http.HttpGateway
}

func NewBoshDirector(ip, username, password string, port int, gateway http.HttpGateway) *BoshDirector {
	return &BoshDirector{
		ip:       ip,
		port:     port,
		username: username,
		password: password,
		gateway:  gateway,
	}
}

func (director *BoshDirector) GetDeploymentManifest(deploymentName string) (manifest io.Reader, err error) {
	endpoint := fmt.Sprintf("https://%s:%d/deployments/%s", director.ip, director.port, deploymentName)
	if err != nil {
		return
	}
	httpEntity := http.HttpRequestEntity{
		Url:         endpoint,
		Username:    director.username,
		Password:    director.password,
		ContentType: "text/yaml",
	}
	request := director.gateway.Get(httpEntity)
	resp, err := request()
	if err != nil {
		return
	}
	manifest, err = retrieveManifest(resp)
	return
}
