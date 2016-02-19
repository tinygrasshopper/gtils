package osutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/osutils"
	"github.com/pivotalservices/gtils/osutils/fakes"
	"golang.org/x/crypto/ssh"
)

var _ = Describe("Given Remote Operations", func() {
	var remoteOperations *RemoteOperations
	var sshConfig command.SshConfig
	BeforeEach(func() {
		sshConfig = command.SshConfig{
			Username: "userId",
			Password: "password",
			Host:     "127.0.0.1",
			Port:     22,
		}

	})
	Describe("given a mockSSHConnection", func() {
		mockSftpclient := &fake.MockSFTPClient{}
		remoteOperations = &RemoteOperations{
			GetSSHConnection: func(config command.SshConfig, clientConfig *ssh.ClientConfig) (sftpclient SFTPClient, err error) {
				sftpclient = mockSftpclient
				return
			},
		}
		Context("calling GetRemoteFile()", func() {
			It("then it should not error", func() {
				_, err := remoteOperations.GetRemoteFile()
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
		Context("calling RemoveRemoteFile()", func() {
			It("then it should not error", func() {
				err := remoteOperations.RemoveRemoteFile()
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
	Describe("given a NewRemoteOperationsWithPath method", func() {

		Context("called on a  with valid configuration and path", func() {

			It("then it should not panic", func() {
				Ω(func() {
					remoteOperations = NewRemoteOperationsWithPath(sshConfig, "/var/temp")
				}).ShouldNot(Panic())
				Ω(remoteOperations).ShouldNot(BeNil())
			})

		})
		Context("called on a  with valid configuration and empty path", func() {

			It("then it should panic", func() {
				Ω(func() {
					remoteOperations = NewRemoteOperationsWithPath(sshConfig, "")
				}).Should(Panic())
			})

		})
	})
	Describe("given a NewRemoteOperations method", func() {

		Context("called on a  with valid configuration and path", func() {

			It("then it should not panic", func() {
				Ω(func() {
					remoteOperations = NewRemoteOperations(sshConfig)
				}).ShouldNot(Panic())
				Ω(remoteOperations).ShouldNot(BeNil())
			})

		})
		Context("called SetPath", func() {

			It("then it should equal the path returned", func() {
				remoteOperations = NewRemoteOperations(sshConfig)
				Ω(remoteOperations).ShouldNot(BeNil())
				remoteOperations.SetPath("blah")
				Ω(remoteOperations.Path()).Should(Equal("blah"))
			})

		})
	})

})
