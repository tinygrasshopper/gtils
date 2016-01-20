package http_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
)

var _ = Describe("given a GetFileSize function", func() {
	Context("when called on a valid file", func() {
		It("then it should return that files size", func() {
			filename := "fixtures/installation.json"
			f, _ := os.Open(filename)
			info, _ := f.Stat()
			Ω(GetFileSize(filename)).Should(Equal(info.Size()))
		})
	})
	Context("when called on a in-valid file", func() {
		It("then it should return -1", func() {
			filename := "this-file-does-not-exist"
			Ω(GetFileSize(filename)).Should(Equal(int64(-1)))
		})
	})
})
