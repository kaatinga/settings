package env_loader

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"log/syslog"
	"os"
	"testing"
	"time"
)

// settings with logrus.Level and time.Duration 4
type settings4 struct {
	Port     uint16        `env:"PORT"`
	LogLevel logrus.Level  `env:"LOG_LEVEL"`
	Timeout  time.Duration `env:"TIMEOUT"`
	Stdout   bool          `env:"STDOUT"`
}

// settings with unsupported int8
type badInt8 struct {
	Port int8 `env:"PORT"`
}

type NotAStruct string

type settingsWithZerologSyslog struct {
	LogLevel       zerolog.Level   `env:"LOG_LEVEL"`
	SyslogPriority syslog.Priority `env:"SYSLOG_LEVEL"`
}

// complex example
type settingsWithStruct struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       *InternalStruct
}

// complex example
type settingsWithStruct2 struct {
	Port           int64  `env:"PORT"`
	PathToDatabase string `env:"DB"`
	Internal       InternalStruct
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
	PathToDatabase string `env:"DB" validate:"required"`
}

// test structure #8
type badEnvironmentSettings1PlusValidation struct {
	Port           string `env:"BADPORT" validate:"numeric"`
	PathToDatabase string `env:"DB" validate:"required"`
}

// test structure #9
type badEnvironmentSettings2PlusValidation struct {
	Port           string `env:"PORT" validate:"numeric"`
	PathToDatabase string `env:"DB" validate:"required"`
	CacheSize      byte   `env:"CACHE" validate:"min=10"`
}

func TestLoadUsingReflect(t *testing.T) {

	// ENV settings PORT=80;DB=db/file;CACHE=5;BADCACHE1=i;BADCACHE2=300
	var err error
	err = os.Setenv("PORT", "80")           // nolint
	err = os.Setenv("DB", "db/file")        // nolint
	err = os.Setenv("CACHE", "5")           // nolint
	err = os.Setenv("BADCACHE1", "i")       // nolint
	err = os.Setenv("BADCACHE2", "300")     // nolint
	err = os.Setenv("BADCACHE3", "-1")      // nolint
	err = os.Setenv("LOG_LEVEL", "debug")   // nolint
	err = os.Setenv("SYSLOG_LEVEL", "info") // nolint
	err = os.Setenv("TIMEOUT", "20s")       // nolint
	err = os.Setenv("BADPORT", "a")
	err = os.Setenv("STDOUT", "true")
	if err != nil {
		t.Errorf("ENV variables has not been set")
	}

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
	var zerologSyslog = settingsWithZerologSyslog{}
	var logrusDurationUint16 = settings4{}
	var unsupportedInt8 = badInt8{}

	tests := []struct {
		name     string
		settings interface{}
		wantErr  error
	}{
		{"ok1", &goodSettings1, nil},
		{"ok2", &goodSettings3withEmptyString, nil},
		{"!ok1", good2, ErrNotAddressable},
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
		{"zerolog and syslog fields", &zerologSyslog, nil},
		{"logrus and duration", &logrusDurationUint16, nil},
		{"int8", &unsupportedInt8, ErrUnsupportedField},
	}

	//nolint
	for _, tt := range tests {

		v := validator.New()

		t.Run(tt.name, func(t *testing.T) {
			err = LoadUsingReflect(tt.settings)

			spew.Dump(tt.settings)

			if errors.Is(err, tt.wantErr) {
				if err != nil {
					t.Log(err)
				}

				err = v.Struct(tt.settings)
				if err != nil {
					t.Logf("additional validation result: %s", err)
				}

				return
			}

			t.Errorf("LoadUsingReflect() error is incorrect\nhave %v\nwant %v", err, tt.wantErr)
		})
	}
}
