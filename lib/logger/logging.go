package logger

import (
	"flag"
	"os"

	logging "github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	`%{color}%{id} ▶ %{time:2006-01-02 15:04:05.000} %{module:.16s} ▶ %{level:.8s} %{color:reset} %{message}`,
)
var errorFormat = logging.MustStringFormatter(
	`%{color}%{id} ▶ %{longfile}%{color:reset}`,
)

var flagLogLevel = flag.String("loglevel", "INFO", "Set loglevel for default logger")

func Init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(getDefaultLogLevel(), "")

	backendError := logging.NewLogBackend(os.Stderr, "", 0)
	backendErrorFormatter := logging.NewBackendFormatter(backendError, errorFormat)
	backendErrorLeveled := logging.AddModuleLevel(backendErrorFormatter)
	backendErrorLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backendLeveled, backendErrorLeveled)
}

func getDefaultLogLevel() logging.Level {
	if level, err := logging.LogLevel(*flagLogLevel); err == nil {
		return level
	} else {
		return logging.INFO
	}
}
