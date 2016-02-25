package uaa_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotalservices/gtils/uaa"
)

var _ = Describe("given a GetToken() function", func() {
	Context("when called with valid uaa target information", func() {
		var token string
		var err error
		var controlToken = "12345937635"
		var server *ghttp.Server

		BeforeEach(func() {
			server = NewTestServer(ghttp.NewTLSServer(), controlToken)
			token, err = uaa.GetToken(server.URL()+"/uaa", "fakeuser", "fakepass", "opsman", "")
		})

		AfterEach(func() {
			server.Close()
		})

		It("Then it should return a valid token", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeEmpty())
			Ω(token).Should(Equal(controlToken))
		})
	})
})

func NewTestServer(server *ghttp.Server, token string) *ghttp.Server {
	tokenJson := getFakeToken("./fixtures/token_response.json", token, "", "")
	successTokenHandler := ghttp.RespondWith(http.StatusOK, tokenJson)
	server.AppendHandlers(
		successTokenHandler,
	)
	return server
}

func getFakeToken(fixturePath, token, refresh, jti string) string {
	b, _ := ioutil.ReadFile(fixturePath)
	return fmt.Sprintf(string(b), token, refresh, jti)
}
