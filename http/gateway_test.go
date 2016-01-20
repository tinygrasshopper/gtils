package http_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
)

var (
	requestCatcher   *http.Request
	roundTripSuccess bool
	httpEntity       = HttpRequestEntity{
		Url:         "http://endpoint/test",
		Username:    "username",
		Password:    "password",
		ContentType: "contentType",
	}
)

type MockRoundTripper struct {
}

func (roundTripper *MockRoundTripper) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: 200,
	}
	if !roundTripSuccess {
		err = errors.New("Mock error")
	}
	*requestCatcher = *request
	return
}

var _ = Describe("Rest client", func() {
	var (
		client = NewRestClient()
	)
	BeforeEach(func() {
		roundTripSuccess = true
		requestCatcher = &http.Request{}
		NewRoundTripper = func() http.RoundTripper {
			return &MockRoundTripper{}
		}
	})
	Context("Calling client NewRequest", func() {
		It("return successfully", func() {

			request, apiResponse := client.NewRequest("GET", "https://opsmanager.com/v1/installation", "admin", "admin", nil)

			Expect(apiResponse.IsSuccessful()).To(Equal(true))
			Expect(request.HttpReq.Header.Get("Authorization")).To(Equal("Basic YWRtaW46YWRtaW4="))
			Expect(request.HttpReq.Header.Get("accept")).To(Equal("application/json"))
		})
	})
})
