package settings

// LoggerOptions — модель данных настроек логгера.
type LoggerOptions struct {
	Syslog         string `env:"SYSLOG" toml:"log.syslog_addr" default:"127.0.0.1:514" validate:"tcp_addr"`
	SyslogProtocol string `env:"SYSLOG_PROTOCOL" toml:"log.syslog_protocol" default:"udp" validate:"min=3,max=3"`
	Colour         bool   `env:"COLOUR" toml:"log.colour" default:"false"`
	StdOut         bool   `env:"STDOUT" toml:"log.stdout" default:"false"`
}
