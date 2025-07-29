package settings

import (
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
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

// settings with logrus.Level and time.Duration 4
type settings4 struct {
	Port    uint16        `env:"PORT"`
	Timeout time.Duration `env:"TIMEOUT"`
	Stdout  bool          `env:"STDOUT"`
}

// settings with unsupported int8
type Int8 struct {
	Port int8 `env:"PORT"`
}

type NotAStruct string

// complex example
type settingsWithStruct struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       *InternalStruct
}

// complex example 2
type settingsWithStruct2 struct {
	Port           int64  `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       InternalStruct
}

// complex example 3 with omitted field
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

// test structure #1
type goodEnvironmentSettings1 struct {
	Port           string `env:"PORT" validate:"required"`
	PathToDatabase string `env:"DB"`
}

// test structure #2
type goodEnvironmentSettings3withEmptyString struct {
	Port           string
	PathToDatabase string `env:"DB"`
}

// test structure #3
type goodEnvironmentSettings2 struct {
	Port           uint32 `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"CACHE"`
}

// test structure #4
type badEnvironmentSettings2 struct {
	Port           uint32 `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE1"`
}

// test structure #5
type badEnvironmentSettings3 struct {
	Port           int64  `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE2"`
}

// test structure #6
type badEnvironmentSettings4 struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE3"`
}

// test structure #7
type goodEnvironmentSettings1PlusValidation struct {
	Port           string `env:"PORT" validate:"numeric"`
	PathToDatabase string `env:"DB"   validate:"required"`
}

// test structure #8
type badEnvironmentSettings1PlusValidation struct {
	Port           string `env:"BADPORT" validate:"numeric"`
	PathToDatabase string `env:"DB"      validate:"required"`
}

// test structure #9
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

func TestLoadUsingReflect(t *testing.T) {
	// ENV settings PORT=80;DB=db/file;CACHE=5;BADCACHE1=i;BADCACHE2=300
	t.Setenv("PORT", "80")
	t.Setenv("FLOAT", "80.1")
	t.Setenv("DB", "db/file")
	t.Setenv("CACHE", "5")
	t.Setenv("BADCACHE1", "i")
	t.Setenv("BADCACHE2", "300")
	t.Setenv("BADCACHE3", "-1")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SYSLOG_LEVEL", "info")
	t.Setenv("TIMEOUT", "20s")
	t.Setenv("BADPORT", "a")
	t.Setenv("STDOUT", "true")
	t.Setenv("SESSION", "session")
	t.Setenv("KEY_PATH", "/etc")
	t.Setenv("PROD", "true")
	t.Setenv("HAS_DB", "true")
	t.Setenv("DOMAIN", "3lines.club")
	t.Setenv("EMAIL", "email@3lines.club")
	t.Setenv("BIG_PORT", "25060")
	t.Setenv("STRING_SLICE", "a,b,c")

	var goodSettings1 goodEnvironmentSettings1
	var goodSettings3withEmptyString goodEnvironmentSettings3withEmptyString
	var good2 goodEnvironmentSettings2
	var badSettings2 badEnvironmentSettings2
	var badSettings3 badEnvironmentSettings3
	var badSettings4 badEnvironmentSettings4
	var goodSettings5 goodEnvironmentSettings1PlusValidation
	var badSettings5 badEnvironmentSettings1PlusValidation
	var badSettings6 badEnvironmentSettings2PlusValidation
	var notAStruct NotAStruct
	var complex1 = settingsWithStruct{}
	var complex2 = &complex1
	var complex3 = settingsWithStruct2{}
	var requiredField = settingsWithRequiredTag{}
	var simple simpleConfig
	var stringSlice struct {
		StringSlice []string `env:"STRING_SLICE"`
	}

	tests := []struct {
		name     string
		settings any
		wantErr  error
	}{
		{"ok1", &goodSettings1, nil},
		{"ok2", &goodSettings3withEmptyString, nil},
		{"!ok1", good2, ErrNotAddressableField},
		{"ok3", &good2, nil},
		{"!ok2", &badSettings2, ErrIncorrectFieldValue},
		{"!ok3", &badSettings3, ErrIncorrectFieldValue},
		{"!ok4", &badSettings4, ErrIncorrectFieldValue},
		{"ok4", &goodSettings5, nil},
		{"!ok5", &badSettings5, ErrValidationFailed},
		{"!ok6", &badSettings6, ErrValidationFailed},
		{"!ok7", notAStruct, ErrNotAStruct},
		{"!ok8", &notAStruct, ErrNotAStruct},
		{"complex double pointer", &complex2, nil},
		{"complex with pointer", &complex1, nil},
		{"complex with a struct without pointer", &complex3, nil},
		{"complex with required tag", &requiredField, ErrValidationFailed},
		{"duration", &settings4{}, nil},
		{"int8", &Int8{}, nil},
		{"int16", &Int16AndFloat64{}, nil},
		{"empty", &emptySettings{}, ErrTheModelHasEmptyStruct},
		{"required if failed", &settingWithRequiredIf{}, ErrValidationFailed},
		{"not_set_env", &simple, nil},
		{"omitted field", &settingsWithStruct3{}, nil},
		{"an example", &Settings{}, nil},
		{"big port", &pigPort{}, nil},
		{"string slice", &stringSlice, nil},
	}

	var err error

	//nolint
	for _, tt := range tests {
		t.Log("\n\n")

		v := validator.New()

		t.Run(tt.name, func(t *testing.T) {
			err = Load(tt.settings)
			t.Logf("settings: %v", tt.settings)

			if errors.Is(err, tt.wantErr) {
				if err != nil {
					t.Log(err)
				}

				err = v.Struct(tt.settings)
				if err != nil {
					t.Logf("additional validation result: %s", err)
				}

				return
			} else if tt.wantErr == ErrValidationFailed {
				validationError, ok := err.(validator.ValidationErrors)
				if ok {
					t.Log(validationError)
					return
				}
			}

			t.Errorf("Load() error is incorrect\nhave %v\nwant %v", err, tt.wantErr)
		})
	}
}
