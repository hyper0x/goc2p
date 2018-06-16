package loadgen

import (
	"logging"
)

var logger logging.Logger

func init() {
	logger = logging.NewSimpleLogger()
}
