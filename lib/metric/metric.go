package metric

import (
	"time"

	logging "github.com/op/go-logging"
)

func Timed(start time.Time, name string, log *logging.Logger) {
	elapsed := time.Since(start)
	log.Debugf("%s took %s", name, elapsed)
}
