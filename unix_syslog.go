// +build !windows

package settings

// parseSyslog возвращает максимальное значение типов uint.
func (field *Loop) parseSyslog() (int64, error) {

	priority, err := ParseSyslogPriority(field.envValue)
	if err != nil {
		return 0, err
	}
	return int64(priority), nil
}
