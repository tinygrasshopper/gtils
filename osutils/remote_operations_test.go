package osutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/osutils"
)

var _ = Describe("Given Remote Operations", func() {
	var remoteOperations *RemoteOperations
	Describe("given a NewRemoteOperationsWithPath method", func() {
		var sshConfig command.SshConfig
		BeforeEach(func() {
			sshConfig = command.SshConfig{
				Username: "userId",
				Password: "password",
				Host:     "127.0.0.1",
				Port:     22,
			}

		})
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

})
