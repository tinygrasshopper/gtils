package log

import (
	"flag"
	"io"

	"github.com/pivotal-golang/lager"
)

type LogType uint

type Log struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Message   string `json:"message"`
	LogLevel  string `json:"log_level"`
	Data      Data   `json:"data"`
}

type Data map[string]interface{}

type Logger interface {
	Debug(message string, data ...Data)
	Info(message string, data ...Data)
	Error(message string, err error, data ...Data)
	Fatal(message string, err error, data ...Data)
}

type logger struct {
	lager.Logger
	Name   string
	Writer io.Writer
}

const (
	Lager LogType = iota
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

var (
	log         *logger
	minLogLevel string
)

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&minLogLevel,
		"logLevel",
		string(INFO),
		"log level: debug, info, error or fatal",
	)
}

func SetLogLevel(level string) {
	minLogLevel = level
}

func LogFactory(name string, logType LogType, writer io.Writer) Logger {
	log := &logger{Name: name, Writer: writer}
	if logType == Lager {
		return NewLager(log)
	}
	return NewLager(log)
}
