package integration_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"code.google.com/p/go-uuid/uuid"
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
		postgresDB   = "postgres"
		postgresUser = os.Getenv("PCFPSQL_ENV_DB_USER")
		postgresPass = os.Getenv("PCFPSQL_ENV_DB_PASS")
		postgresPort = 5432
		sshUser      = os.Getenv("PCFPSQL_ENV_SSH_USER")
		sshPass      = os.Getenv("PCFPSQL_ENV_SSH_PASS")
		sshHost      = os.Getenv("PCFPSQL_PORT_22_TCP_ADDR")
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
				PGDMP_SQL_BIN = strings.TrimPrefix(PGDMP_SQL_BIN, cfBinDir)
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
				PGDMP_SQL_BIN = strings.TrimPrefix(PGDMP_SQL_BIN, cfBinDir)
				inputReader, _ = os.Open("fixtures/postgres_dump.txt")
				pgRemoteDump, sshConnectionErr = NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteCommandErr = pgRemoteDump.Import(inputReader)
			})

			It("Then we should not see any ssh connection errors", func() {
				Ω(sshConnectionErr).ShouldNot(HaveOccurred())
			})

			It("Then we should not see any errors from the remote command execution", func() {
				Ω(remoteCommandErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Given a Dump method", func() {
		Context("When called against a valid postgres instance", func() {
			var (
				controlDumpfileChecksum = "7fbcef4c4fd53f25847b08a7dd49cb72"
				outputWriter            bytes.Buffer
				sshConnectionErr        error
				remoteCommandErr        error
			)
			BeforeEach(func() {
				PGDMP_SQL_BIN = "psql"
				inputReader, _ := os.Open("fixtures/postgres_dump.txt")
				remoteEnvStager, _ := NewPgRemoteDump(postgresPort, postgresDB, postgresUser, postgresPass, sshConfig)
				remoteEnvStager.Import(inputReader)
				PGDMP_DUMP_BIN = strings.TrimPrefix(PGDMP_DUMP_BIN, cfBinDir)
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
				dumpfilebytes := outputWriter.Bytes()
				hash := fmt.Sprintf("%x", md5.Sum(dumpfilebytes))
				Ω(hash).Should(Equal(controlDumpfileChecksum))
			})
		})
	})
})
