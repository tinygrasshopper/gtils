package integration_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pborman/uuid"
	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pgCatchCommand string
)

var _ = Describe("PgDump Integration Tests", func() {
	const (
		cfBinDir = "/var/vcap/packages/postgres/bin/"
	)
	var (
		postgresDB   = "console"
		postgresUser = os.Getenv("UAADB_INT_ENV_DB_USER")
		postgresPass = os.Getenv("UAADB_INT_ENV_DB_PASS")
		postgresPort = 2544
		sshUser      = os.Getenv("UAADB_INT_ENV_SSH_USER")
		sshPass      = os.Getenv("UAADB_INT_ENV_SSH_PASS")
		sshHost      = os.Getenv("UAADB_INT_PORT_22_TCP_ADDR")
		pgRemoteDump *PgDump
		sshConfig    = command.SshConfig{
			Username: sshUser,
			Password: sshPass,
			Host:     sshHost,
			Port:     22,
		}
	)

	Describe("Given a Import method", func() {
		Context("When a call to postgres yields an error", func() {
			var (
				inputReader      io.Reader
				sshConnectionErr error
				remoteCommandErr error
			)
			BeforeEach(func() {
				databaseNameGUID := uuid.New()
				PGDmpRestoreBin = strings.TrimPrefix(PGDmpRestoreBin, cfBinDir)
				inputReader, _ = os.Open("pgdump_test.go")
				pgRemoteDump, sshConnectionErr = NewPgRemoteDump(postgresPort, databaseNameGUID, postgresUser, postgresPass, sshConfig)
				remoteCommandErr = pgRemoteDump.Import(inputReader)
			})

			It("Then we should not see any ssh connection errors", func() {
				Ω(sshConnectionErr).ShouldNot(HaveOccurred())
			})

			It("Then we should see an error returned from the remote command call", func() {
				Ω(remoteCommandErr).Should(HaveOccurred())
			})
		})

		Context("When called against a valid postgres instance with a valid dumpfile", func() {
			var (
				inputReader      io.Reader
				sshConnectionErr error
				remoteCommandErr error
			)
			BeforeEach(func() {
				PGDmpRestoreBin = strings.TrimPrefix(PGDmpRestoreBin, cfBinDir)
				inputReader, _ = os.Open("fixtures/tst.dmp")
				pgRemoteDump, sshConnectionErr = NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteCommandErr = pgRemoteDump.Import(inputReader)
			})

			It("Then we should not see any ssh connection errors", func() {
				Ω(sshConnectionErr).ShouldNot(HaveOccurred())
			})

			XIt("Then we should not see any errors from the remote command execution - ignored b/c we need to configure a proper container to test against", func() {
				Ω(remoteCommandErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Given a Dump method", func() {
		Context("When called against a valid postgres instance", func() {
			var (
				controlDumpfileChecksum = "d41d8cd98f00b204e9800998ecf8427e"
				outputWriter            bytes.Buffer
				sshConnectionErr        error
				remoteCommandErr        error
			)
			BeforeEach(func() {
				PGDmpRestoreBin = strings.TrimPrefix(PGDmpRestoreBin, cfBinDir)
				inputReader, _ := os.Open("fixtures/tst.dmp")
				remoteEnvStager, _ := NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteEnvStager.Import(inputReader)
				PGDmpDumpBin = strings.TrimPrefix(PGDmpDumpBin, cfBinDir)
				pgRemoteDump, sshConnectionErr = NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteCommandErr = pgRemoteDump.Dump(&outputWriter)
			})

			It("Then we should not see any ssh connection errors - ignored b/c we need to configure a proper container to test against", func() {
				Ω(sshConnectionErr).ShouldNot(HaveOccurred())
			})

			XIt("Then we should not see any errors from the remote command execution", func() {
				Ω(remoteCommandErr).ShouldNot(HaveOccurred())
			})

			It("Then it should yield a valid dumpfile", func() {
				dumpfilebytes := outputWriter.Bytes()
				hash := fmt.Sprintf("%x", md5.Sum(dumpfilebytes))
				Ω(hash).Should(Equal(controlDumpfileChecksum))
			})
		})
	})
})
