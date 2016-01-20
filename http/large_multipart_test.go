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
			httpServerMock *mock.HTTPServer = &mock.HTTPServer{}
			request        *http.Request
			controlFile    []byte
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
			filePath := fmt.Sprintf("fixtures/%s", fileName)
			conn := ConnAuth{
				Url: httpServerMock.Server.URL,
			}
			fileRef, _ := os.Open(filePath)
			defer fileRef.Close()
			res, err = LargeMultiPartUpload(conn, paramName, filePath, getFileSize(filePath), fileRef, controlParams)
			controlFile, _ = ioutil.ReadFile(filePath)
		})

		AfterEach(func() {
			lo.G.Debug("tearing down server")
			multipartForm = nil
			httpServerMock.Teardown()
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

func getFileSize(filename string) (fileSize int64) {
	var (
		fileInfo os.FileInfo
		err      error
		file     *os.File
	)

	if file, err = os.Open(filename); err == nil {
		fileInfo, err = file.Stat()
		fileSize = fileInfo.Size()
	}

	if err != nil {
		fileSize = -1
	}
	return
}
