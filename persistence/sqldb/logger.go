package sqldb

import (
	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
)

const (
	// printfArg is number of printf args
	printfArg = 3
	// perfMonitorQuery is performance monitoring query
	perfMonitorQuery = "select count(id) from INFORMATION_SCHEMA.PROCESSLIST WHERE time > 10 AND " +
		"command IN ('Execute', 'Query') ORDER BY time DESC"
)

// GorpLogger is gorp logger
type GorpLogger struct {
}

// NewGorpLogger is a constructor of gorp logger
func NewGorpLogger() gorp.GorpLogger {
	return &GorpLogger{}
}

// Printf is printf method
func (g *GorpLogger) Printf(format string, args ...interface{}) {
	m := logging.Fields{}

	argString := []string{"", "query_string", "args", "execution_time"}
	for k := range args {
		if k == 0 || k > printfArg {
			continue
		}
		m[argString[k]] = args[k]
	}

	if queryString, ok := m["query_string"]; ok &&
		queryString == perfMonitorQuery {
		return
	}

	logging.Log.WithFields(logging.DepthLow, m).Info("query")
}
