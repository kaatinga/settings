package settings

import (
	"log/syslog"
	"testing"
)

func TestPriorityDescription(t *testing.T) {

	tests := []struct {
		priority syslog.Priority
		want     string
	}{
		{syslog.LOG_CRIT, "crit"},
		{syslog.LOG_WARNING, "warning"},
		{syslog.LOG_ERR, "error"},
		{syslog.LOG_DEBUG, "debug"},
		{syslog.LOG_INFO, "info"},
		{syslog.LOG_EMERG, "emerg"},
	}
	
	//nolint
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := PriorityDescription(tt.priority); got != tt.want {
				t.Errorf("PriorityDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
