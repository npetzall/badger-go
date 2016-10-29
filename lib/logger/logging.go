package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	`%{color}%{id} ▶ %{time:2006-01-02 15:04:05.000} %{module:.16s} ▶ %{level:.8s} %{color:reset} %{message}`,
)
var errorFormat = logging.MustStringFormatter(
	`%{color}%{id} ▶ %{longfile}%{color:reset}`,
)

func Init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendError := logging.NewLogBackend(os.Stderr, "", 0)
	backendErrorFormatter := logging.NewBackendFormatter(backendError, errorFormat)
	backendErrorLeveled := logging.AddModuleLevel(backendErrorFormatter)
	backendErrorLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backendFormatter, backendErrorLeveled)
}
