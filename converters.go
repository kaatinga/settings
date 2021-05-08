package env_loader

import (
	"log/syslog"
	"strings"
)

// ParseSyslogPriority converts string to syslog.Priority.
func ParseSyslogPriority(lvl string) (syslog.Priority, error) {

	switch strings.ToLower(lvl) {
	case "panic":
		return syslog.LOG_EMERG, nil
	case "fatal":
		return syslog.LOG_CRIT, nil
	case "error":
		return syslog.LOG_ERR, nil
	case "warn", "warning":
		return syslog.LOG_WARNING, nil
	case "info":
		return syslog.LOG_INFO, nil
	case "debug":
		return syslog.LOG_DEBUG, nil
	case "trace":
		return syslog.LOG_NOTICE, nil
	}

	return 0, incorrectPriority(lvl)
}

// PriorityDescription returns description for syslog.Priority.
func PriorityDescription(priority syslog.Priority) string { //nolint:unused
	switch priority {
	case syslog.LOG_DEBUG:
		return "debug"
	case syslog.LOG_ERR:
		return "error"
	case syslog.LOG_CRIT:
		return "crit"
	case syslog.LOG_INFO:
		return "info"
	case syslog.LOG_EMERG:
		return "emerg"
	case syslog.LOG_WARNING:
		return "warning"
	default:
		return "unknown syslog level"
	}
}
