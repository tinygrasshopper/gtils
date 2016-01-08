package storage_test

import (
	"bytes"
	"crypto/rand"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/storage"
)

var _ = Describe("S3bucket", func() {

	testBin := make([]byte, 1<<8)
	rand.Read(testBin)
	s3, err := SafeCreateS3Bucket("s3.amazonaws.com", "pcfbackup-files", "", "")

	PIt("should write", func() {
		Ω(err).To(BeNil())
		w, err := s3.NewWriter("remove")
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	PIt("should read", func() {
		r, err := s3.NewReader("remove")
		Ω(err).To(BeNil())
		Ω(r).ToNot(BeNil())
		out1 := new(bytes.Buffer)
		l, err := io.Copy(out1, r)
		Ω(err).To(BeNil())
		r.Close()
		Ω(len(testBin) == int(l)).To(BeTrue())
		Ω(bytes.Equal(testBin, out1.Bytes())).To(BeTrue())
	})

	PIt("should delete", func() {
		Ω(s3.Delete("remove")).To(BeNil())

		_, err := s3.NewReader("remove")
		Ω(err).ToNot(BeNil())
	})

})
