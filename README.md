# env_loader
The package looks up necessary environment variables and use them to set settings for application.

The settings must be formed as struct with byte and string fields.

Example:

```go
...
type EnvironmentSettings struct {
	Port       string `env:"PORT"`
	ShadowPath string `env:"USER_LIST"`
	KeyPath    string `env:"SESSION_KEYS"`
	Database   string `env:"DATABASE"`
	CacheSize  uint64 `env:"CACHE_SIZE"`
	LaunchMode string `env:"LAUNCH_MODE"`
}
...
```
