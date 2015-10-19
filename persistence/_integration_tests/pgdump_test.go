package integration_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"

	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pgCatchCommand string
)

var _ = Describe("PgDump Integration Tests", func() {
	var (
		postgresDB   = "postgres"
		postgresUser = os.Getenv("PCFPSQL_ENV_DB_USER")
		postgresPass = os.Getenv("PCFPSQL_ENV_DB_PASS")
		postgresPort = 5432
		sshUser      = os.Getenv("PCFPSQL_ENV_SSH_USER")
		sshPass      = os.Getenv("PCFPSQL_ENV_SSH_PASS")
		sshHost      = os.Getenv("PCFPSQL_PORT_22_TCP_ADDR")
		pgRemoteDump *PgDump
	)

	Describe("Given a Dump method", func() {
		Context("When called against a valid postgres instance", func() {
			var (
				controlDumpfileChecksum = "2f4c406a0a7f12193af1c2966fe56a54"
				outputWriter            bytes.Buffer
				sshConnectionErr        error
				remoteCommandErr        error
			)
			BeforeEach(func() {
				PGDMP_DUMP_BIN = "pg_dump"
				sshConfig := command.SshConfig{
					Username: sshUser,
					Password: sshPass,
					Host:     sshHost,
					Port:     22,
				}
				pgRemoteDump, sshConnectionErr = NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteCommandErr = pgRemoteDump.Dump(&outputWriter)
			})

			It("Then we should not see any ssh connection errors", func() {
				Ω(sshConnectionErr).ShouldNot(HaveOccurred())
			})

			It("Then we should not see any errors from the remote command execution", func() {
				Ω(remoteCommandErr).ShouldNot(HaveOccurred())
			})

			It("Then it should yield a valid dumpfile", func() {
				hash := fmt.Sprintf("%x", md5.Sum(outputWriter.Bytes()))
				Ω(hash).Should(Equal(controlDumpfileChecksum))
			})
		})
	})
})
