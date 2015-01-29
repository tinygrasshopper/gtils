package bosh_test

import (
	"errors"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/bosh"
	. "github.com/pivotalservices/gtils/http"
	. "github.com/pivotalservices/gtils/http/httptest"
)

var (
	responseMock          http.Response
	httpRequestSuccessful bool
	ip                    string = "10.10.10.10"
	port                  int    = 25555
	username              string = "test"
	password              string = "test"
	capturedEntity        HttpRequestEntity
)

func captureEntity(entity HttpRequestEntity) {
	capturedEntity = entity
}

var fakeGetAdaptor RequestAdaptor = func() (response *http.Response, err error) {
	if !httpRequestSuccessful {
		err = errors.New("Http Request failed")
		return
	}
	response = &responseMock
	return
}

var _ = Describe("Bosh", func() {
	var (
		gateway *MockGateway = &MockGateway{
			FakeGetAdaptor: fakeGetAdaptor,
			Capture:        captureEntity,
		}
		director = NewBoshDirector(ip, username, password, port, gateway)
	)
	Describe("Bosh Deployment", func() {
		Context("Http Request Failure", func() {
			It("Should return error", func() {
				httpRequestSuccessful = false
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("Http Request Succeed", func() {
			BeforeEach(func() {
				httpRequestSuccessful = true
			})
			It("Should return error when http status code is non 200", func() {
				responseMock = readFakeResponse(500, "fixtures/manifest.yml")
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return error when http response is not valid manifest", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest_invalid.yml")
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return nil error when http response is valid manifest", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest.yml")
				_, err := director.GetDeploymentManifest("name")
				Ω(err).Should(BeNil())
			})
			It("Should compose correct request when the http response is valid", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest.yml")
				director.GetDeploymentManifest("name")
				Ω(capturedEntity.Url).Should(Equal("https://10.10.10.10:25555/deployments/name"))
				Ω(capturedEntity.Username).Should(Equal("test"))
				Ω(capturedEntity.Password).Should(Equal("test"))
				Ω(capturedEntity.ContentType).Should(Equal("text/yaml"))
			})
		})
	})

})

func readFakeResponse(statusCode int, file string) (res http.Response) {
	body, _ := os.Open(file)
	response := http.Response{}
	response.Body = body
	response.StatusCode = statusCode
	return response
}
