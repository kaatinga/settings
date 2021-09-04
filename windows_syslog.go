// +build windows

package settings

// parseSyslog заглушает обработку сислога под Windows
func (field *Loop) parseSyslog() (int64, error) {
	return 0, nil
}
