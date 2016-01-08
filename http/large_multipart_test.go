package http_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
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
			httpServerMock *mock.HttpServer = &mock.HttpServer{}
			request        *http.Request
			controlFile    []byte
			res            *http.Response
			err            error
			tmpfile        = fmt.Sprintf("/tmp/upload-", uuid.New())
			headerSize     = 272
		)

		BeforeEach(func() {
			httpServerMock.Setup()
			httpServerMock.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				request = r
				f, err := os.Create(tmpfile)
				fmt.Println(err)
				io.Copy(f, r.Body)
				r.Body.Close()
				f.Close()
			})
			filePath := fmt.Sprintf("fixtures/%s", fileName)
			conn := ConnAuth{
				Url: httpServerMock.Server.URL,
			}
			fileRef, _ := os.Open(filePath)
			defer fileRef.Close()
			res, err = LargeMultiPartUpload(conn, paramName, filePath, fileRef, nil)
			controlFile, _ = ioutil.ReadFile(filePath)
		})

		AfterEach(func() {
			lo.G.Debug("tearing down server")
			httpServerMock.Teardown()
			os.Remove(tmpfile)
		})

		It("should send the file to the server", func() {
			f, _ := ioutil.ReadFile(tmpfile)
			Ω(err).Should(BeNil())
			Ω(len(f)).Should(Equal(len(controlFile) + headerSize))
			Ω(res.StatusCode).Should(Equal(200))
			Ω(request.Method).Should(Equal("POST"))
			Ω(res).ShouldNot(BeNil())
			Ω(request.ContentLength).Should(Equal(res.Request.ContentLength))
		})
	})
})
