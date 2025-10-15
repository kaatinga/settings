package settings

import (
	"errors"
	"testing"
	"time"
)

type emptySettings struct{}

type Settings struct {
	SessionName   string `env:"SESSION"  validate:"required"`
	PublicKeyPath string `env:"KEY_PATH" validate:"required"`
}

type settingWithRequiredIf struct {
	RequiredIf string `env:"NOTFOUND" validate:"required_if=Trigger true"`
	Trigger    bool   `env:"STDOUT"`
}

type Int16AndFloat64 struct {
	PORT  int16   `env:"PORT"`
	FLOAT float64 `env:"FLOAT"`
}

type settings4 struct {
	Port    uint16        `env:"PORT"`
	Timeout time.Duration `env:"TIMEOUT"`
	Stdout  bool          `env:"STDOUT"`
}

type Int8 struct {
	Port int8 `env:"PORT"`
}

type NotAStruct string

type settingsWithStruct struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       *InternalStruct
}

type settingsWithStruct2 struct {
	Port           int64  `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       InternalStruct
}

type settingsWithStruct3 struct {
	Port           int64           `env:"PORT"`
	PathToDatabase string          `env:"DB"`
	Internal       *InternalStruct `env:"-"`
}

type settingsWithRequiredTag struct {
	PathToDatabase string `env:"DB2" validate:"required"`
}

type InternalStruct struct {
	CacheSize string `env:"CACHE"`
}

type goodEnvironmentSettings1 struct {
	Port           string `env:"PORT" validate:"required"`
	PathToDatabase string `env:"DB"`
}

type goodEnvironmentSettings3withEmptyString struct {
	Port           string
	PathToDatabase string `env:"DB"`
}

type goodEnvironmentSettings2 struct {
	Port           uint32 `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"CACHE"`
}

type badEnvironmentSettings2 struct {
	Port           uint32 `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE1"`
}

type badEnvironmentSettings3 struct {
	Port           int64  `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE2"`
}

type badEnvironmentSettings4 struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE3"`
}

type goodEnvironmentSettings1PlusValidation struct {
	Port           string `env:"PORT" validate:"numeric"`
	PathToDatabase string `env:"DB"   validate:"required"`
}

type badEnvironmentSettings1PlusValidation struct {
	Port           string `env:"BADPORT" validate:"numeric"`
	PathToDatabase string `env:"DB"      validate:"required"`
}

type badEnvironmentSettings2PlusValidation struct {
	Port           string `env:"PORT"  validate:"numeric"`
	PathToDatabase string `env:"DB"    validate:"required"`
	CacheSize      byte   `env:"CACHE" validate:"min=10"`
}

type simpleConfig struct {
	DBURL   string        `default:"127.0.0.1" env:"DB_URL"      validate:"required"`
	Timeout time.Duration `default:"5s"        env:"DB_TIMEOUT"`
}

type pigPort struct {
	Port uint16 `env:"BIG_PORT" validate:"required"`
}

type stringSliceConfig struct {
	StringSlice []string `env:"STRING_SLICE"`
}

// setupTestEnvironment sets up all required environment variables for testing
func setupTestEnvironment(t *testing.T) {
	t.Helper()

	envVars := map[string]string{
		"PORT":         "80",
		"FLOAT":        "80.1",
		"DB":           "db/file",
		"CACHE":        "5",
		"BADCACHE1":    "i",
		"BADCACHE2":    "300",
		"BADCACHE3":    "-1",
		"LOG_LEVEL":    "debug",
		"SYSLOG_LEVEL": "info",
		"TIMEOUT":      "20s",
		"BADPORT":      "a",
		"STDOUT":       "true",
		"SESSION":      "session",
		"KEY_PATH":     "/etc",
		"PROD":         "true",
		"HAS_DB":       "true",
		"DOMAIN":       "example.com",
		"EMAIL":        "email@example.com",
		"BIG_PORT":     "25060",
		"STRING_SLICE": "a,b,c",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}
}

// TestLoad tests the Load function with various configurations
func TestLoad(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name        string
		settings    any
		wantErr     error
		expectError bool
		description string
	}{
		{
			name:        "valid_settings_with_required_validation",
			settings:    &goodEnvironmentSettings1{},
			wantErr:     nil,
			expectError: false,
			description: "Should successfully load settings with required field validation",
		},
		{
			name:        "valid_settings_with_empty_string_field",
			settings:    &goodEnvironmentSettings3withEmptyString{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle fields without env tags",
		},
		{
			name:        "non_pointer_struct_should_fail",
			settings:    goodEnvironmentSettings2{},
			wantErr:     ErrNotAddressableField,
			expectError: true,
			description: "Should fail when struct is not passed as pointer",
		},
		{
			name:        "valid_settings_with_uint32_and_byte",
			settings:    &goodEnvironmentSettings2{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle uint32 and byte types correctly",
		},
		{
			name:        "invalid_byte_value_should_fail",
			settings:    &badEnvironmentSettings2{},
			wantErr:     NewIncorrectFieldValueError("BADCACHE1"),
			expectError: true,
			description: "Should fail when byte field has invalid value",
		},
		{
			name:        "byte_value_out_of_range_should_fail",
			settings:    &badEnvironmentSettings3{},
			wantErr:     NewIncorrectFieldValueError("BADCACHE2"),
			expectError: true,
			description: "Should fail when byte value exceeds range",
		},
		{
			name:        "negative_byte_value_should_fail",
			settings:    &badEnvironmentSettings4{},
			wantErr:     NewIncorrectFieldValueError("BADCACHE3"),
			expectError: true,
			description: "Should fail when byte value is negative",
		},
		{
			name:        "valid_settings_with_numeric_validation",
			settings:    &goodEnvironmentSettings1PlusValidation{},
			wantErr:     nil,
			expectError: false,
			description: "Should pass numeric validation for valid port",
		},
		{
			name:        "invalid_numeric_validation_should_fail",
			settings:    &badEnvironmentSettings1PlusValidation{},
			wantErr:     nil, // Will be a validation error from Load function
			expectError: true,
			description: "Should fail validation during load",
		},
		{
			name:        "min_validation_should_fail",
			settings:    &badEnvironmentSettings2PlusValidation{},
			wantErr:     nil, // Will be a validation error from Load function
			expectError: true,
			description: "Should fail min validation during load",
		},
		{
			name:        "non_struct_type_should_fail",
			settings:    NotAStruct("test"),
			wantErr:     ErrNotAStruct,
			expectError: true,
			description: "Should fail when input is not a struct",
		},
		{
			name:        "non_struct_pointer_should_fail",
			settings:    func() *NotAStruct { s := NotAStruct("test"); return &s }(),
			wantErr:     ErrNotAStruct,
			expectError: true,
			description: "Should fail when pointer points to non-struct",
		},
		{
			name:        "complex_settings_with_double_pointer",
			settings:    func() any { s := &settingsWithStruct{}; return &s }(),
			wantErr:     nil,
			expectError: false,
			description: "Should handle double pointer to struct",
		},
		{
			name:        "complex_settings_with_pointer",
			settings:    &settingsWithStruct{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle pointer to struct with nested struct",
		},
		{
			name:        "complex_settings_without_pointer",
			settings:    &settingsWithStruct2{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle struct with nested struct (not pointer)",
		},
		{
			name:        "required_field_missing_should_fail",
			settings:    &settingsWithRequiredTag{},
			wantErr:     NewValidationFailedError("PathToDatabase", "string", "required"),
			expectError: true,
			description: "Should fail when required field is missing",
		},
		{
			name:        "settings_with_duration",
			settings:    &settings4{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle time.Duration type correctly",
		},
		{
			name:        "settings_with_int8",
			settings:    &Int8{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle int8 type correctly",
		},
		{
			name:        "settings_with_int16_and_float64",
			settings:    &Int16AndFloat64{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle int16 and float64 types correctly",
		},
		{
			name:        "empty_struct_should_fail",
			settings:    &emptySettings{},
			wantErr:     ErrTheModelHasEmptyStruct,
			expectError: true,
			description: "Should fail when struct has no fields",
		},
		{
			name:        "required_if_validation_should_fail",
			settings:    &settingWithRequiredIf{},
			wantErr:     nil, // Will be a validation error from Load function
			expectError: true,
			description: "Should fail required_if validation during load",
		},
		{
			name:        "settings_with_default_values",
			settings:    &simpleConfig{},
			wantErr:     nil,
			expectError: false,
			description: "Should use default values when env vars are not set",
		},
		{
			name:        "settings_with_omitted_field",
			settings:    &settingsWithStruct3{},
			wantErr:     nil,
			expectError: false,
			description: "Should skip fields marked with env:\"-\"",
		},
		{
			name:        "example_settings",
			settings:    &Settings{},
			wantErr:     nil,
			expectError: false,
			description: "Should load example settings successfully",
		},
		{
			name:        "settings_with_big_port",
			settings:    &pigPort{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle large port numbers correctly",
		},
		{
			name:        "settings_with_string_slice",
			settings:    &stringSliceConfig{},
			wantErr:     nil,
			expectError: false,
			description: "Should handle string slice from comma-separated values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.settings != nil {
				tt.settings = createFreshInstance(tt.settings)
			}

			err := Load(tt.settings)

			if tt.expectError {
				if err == nil {
					t.Errorf("Load() expected error but got nil")
					return
				}

				if tt.wantErr == nil {
					t.Logf("Expected validation error: %v", err)
					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Load() unexpected error = %v", err)
					return
				}
			}
		})
	}
}

func createFreshInstance(original any) any {
	switch original.(type) {
	case *goodEnvironmentSettings1:
		return &goodEnvironmentSettings1{}
	case *goodEnvironmentSettings3withEmptyString:
		return &goodEnvironmentSettings3withEmptyString{}
	case goodEnvironmentSettings2:
		return goodEnvironmentSettings2{}
	case *goodEnvironmentSettings2:
		return &goodEnvironmentSettings2{}
	case *badEnvironmentSettings2:
		return &badEnvironmentSettings2{}
	case *badEnvironmentSettings3:
		return &badEnvironmentSettings3{}
	case *badEnvironmentSettings4:
		return &badEnvironmentSettings4{}
	case *goodEnvironmentSettings1PlusValidation:
		return &goodEnvironmentSettings1PlusValidation{}
	case *badEnvironmentSettings1PlusValidation:
		return &badEnvironmentSettings1PlusValidation{}
	case *badEnvironmentSettings2PlusValidation:
		return &badEnvironmentSettings2PlusValidation{}
	case NotAStruct:
		return NotAStruct("test")
	case *NotAStruct:
		s := NotAStruct("test")
		return &s
	case *settingsWithStruct:
		return &settingsWithStruct{}
	case *settingsWithStruct2:
		return &settingsWithStruct2{}
	case *settingsWithStruct3:
		return &settingsWithStruct3{}
	case *settingsWithRequiredTag:
		return &settingsWithRequiredTag{}
	case *settings4:
		return &settings4{}
	case *Int8:
		return &Int8{}
	case *Int16AndFloat64:
		return &Int16AndFloat64{}
	case *emptySettings:
		return &emptySettings{}
	case *settingWithRequiredIf:
		return &settingWithRequiredIf{}
	case *simpleConfig:
		return &simpleConfig{}
	case *Settings:
		return &Settings{}
	case *pigPort:
		return &pigPort{}
	case *stringSliceConfig:
		return &stringSliceConfig{}
	default:
		return original
	}
}
