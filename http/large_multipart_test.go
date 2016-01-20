package http_test

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/mock"
	"github.com/xchapter7x/lo"
)

var _ = Describe("Multipart", func() {

	var (
		paramName = "installation[file]"
		fileName  = "installation.json"
	)
	Describe("LargeMultiPartUpload", func() {
		var (
			conn           ConnAuth
			fixtureFileRef *os.File
			filePath                        = fmt.Sprintf("fixtures/%s", fileName)
			httpServerMock *mock.HTTPServer = &mock.HTTPServer{}
			request        *http.Request
			res            *http.Response
			err            error
			multipartForm  *multipart.Form
			controlParams  = map[string]string{
				"password": "test-pass",
			}
			controlValue = map[string][]string{
				"password": []string{"test-pass"},
			}
		)

		BeforeEach(func() {
			httpServerMock.Setup()
			httpServerMock.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				request = r
				request.ParseMultipartForm(1000000)
				multipartForm = request.MultipartForm
				fmt.Println(err)
			})
			conn = ConnAuth{
				Url: httpServerMock.Server.URL,
			}
			fixtureFileRef, _ = os.Open(filePath)
		})

		AfterEach(func() {
			lo.G.Debug("tearing down server")
			multipartForm = nil
			httpServerMock.Teardown()
			fixtureFileRef.Close()
		})

		Context("when called with a -1 filesize value", func() {
			BeforeEach(func() {
				filesize := int64(-1)
				res, err = LargeMultiPartUpload(conn, paramName, filePath, filesize, fixtureFileRef, controlParams)
			})

			It("should interogate the file and insert the actual filesize", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with valid args and actual filesize value", func() {
			BeforeEach(func() {
				res, err = LargeMultiPartUpload(conn, paramName, filePath, GetFileSize(filePath), fixtureFileRef, controlParams)
			})
			It("should pass in the password paramater as part of the request", func() {
				Ω(multipartForm.Value).Should(Equal(controlValue))
			})

			It("should send the file to the server", func() {
				fixturePath := fmt.Sprintf("fixtures/%s", fileName)
				fixtureBytes, _ := ioutil.ReadFile(fixturePath)
				file, _ := multipartForm.File["installation[file]"][0].Open()
				fileBytes, _ := ioutil.ReadAll(file)

				Ω(err).Should(BeNil())
				Ω(fileBytes).Should(Equal(fixtureBytes))
				Ω(res.StatusCode).Should(Equal(200))
				Ω(request.Method).Should(Equal("POST"))
				Ω(res).ShouldNot(BeNil())
				Ω(request.ContentLength).Should(Equal(res.Request.ContentLength))
			})
		})
	})
})
