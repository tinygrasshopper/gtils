package http_test

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/mock"
)

var _ = Describe("Multipart", func() {

	var (
		paramName = "installation[file]"
		fileName  = "installation.json"
	)
	Describe("LargeMultiPartUpload", func() {
		var (
			httpServerMock *mock.HttpServer = &mock.HttpServer{}
			request        *http.Request
		)

		BeforeEach(func() {
			httpServerMock.Setup()
			httpServerMock.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				request = r
			})
		})

		AfterEach(func() {
			httpServerMock.Teardown()
		})

		It("should send the file to the server", func() {
			filePath := fmt.Sprintf("fixtures/%s", fileName)
			conn := ConnAuth{
				Url: httpServerMock.Server.URL,
			}
			fileRef, _ := os.Open(filePath)
			res, err := LargeMultiPartUpload(conn, paramName, filePath, fileRef, nil)
			Ω(err).Should(BeNil())
			Ω(res.StatusCode).Should(Equal(200))
			Ω(request.Method).Should(Equal("POST"))
			Ω(res).ShouldNot(BeNil())
			Ω(request.ContentLength).Should(Equal(res.Request.ContentLength))
		})
	})
})
