package settings

import (
	"log/syslog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

// LoadSettings loads settings to a struct.
func LoadSettings(settings interface{}) error {

	engine, nestedStruct := settings.(*Engine)
	if !nestedStruct {
		engine = newEngine(settings)
	}

	err := engine.getStruct()
	if err != nil {
		return err
	}

	for i := 0; i < engine.NumberOfFields; i++ {
		engine.startIteration(i)

		if engine.Field.value.Kind() == reflect.Ptr ||
			engine.Field.value.Kind() == reflect.Struct {
			// we check whether the field is pointer or struct

			err = LoadSettings(&Engine{
				Value: engine.Field.value,
				Type:  engine.Field.value.Type(),
			})
			if err != nil {
				return err
			}
			continue

		} else {

			// if a field has no env tag, we pass such a field
			if !engine.Field.hasEnvTag {
				continue
			}

			// we check if it is required
			engine.validateRequired()
			//fmt.Println(engine.Field.field.Name, "required:", engine.Field.required)

			// if a field has env tag, but the env was not found, and if it is required
			// we return error
			engine.Field.envValue, engine.Field.hasEnvValue = os.LookupEnv(engine.Field.envTag)
			if !engine.Field.hasEnvValue {
				if engine.Field.required {
					return engine.validationFailed()
				}

				if engine.Field.hasDefaultSetting {
					// substitute the envValue with default setting
					engine.Field.envValue = engine.Field.defaultSetting
				} else {
					// finish processing the current field
					continue
				}
			}

			if !engine.Value.Field(i).IsValid() {
				return ErrInternalFailure
			}

			// The struct must be send via pointer.
			if !engine.Value.Field(i).CanSet() {
				return ErrNotAddressable
			}

			switch engine.Field.value.Kind() {
			case reflect.String:
				engine.Field.value.SetString(engine.Field.envValue)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:

				if engine.Field.value.Kind() == reflect.Uint32 &&
					engine.Field.value.Type().String() == logrusLevel {
					// check if it is logrus.level

					var level logrus.Level
					level, err = logrus.ParseLevel(engine.Field.envValue)
					if err != nil {
						return incorrectFieldValue(engine.Field.envTag)
					}
					engine.Field.uint64Value = uint64(level)

				} else {
					// uint

					engine.Field.uint64Value, err = strconv.ParseUint(engine.Field.envValue, 10, 64)
					if err != nil {
						return incorrectFieldValue(engine.Field.envTag)
					}

					// check if whether the value exceeds the type maximum or not
					if engine.Field.uint64Value > maximum(engine.Field.value.Kind()) {
						return incorrectFieldValue(engine.Field.envTag)
					}
				}

				engine.Field.value.SetUint(engine.Field.uint64Value)

			case reflect.Int8:

				if engine.Field.value.Type().String() == zerologLevel {
					// check if it is zerolog.Level

					var level zerolog.Level
					level, err = zerolog.ParseLevel(engine.Field.envValue)
					if err != nil {
						return incorrectFieldValue(engine.Field.envTag)
					}

					engine.Field.value.SetInt(int64(level))
					break
				}

				return unsupportedField(engine.Field.value.Type().Name())

			case reflect.Int64, reflect.Int:

				if engine.Field.value.Kind() == reflect.Int &&
					engine.Field.value.Type().String() == syslogPriority {
					// check if it is syslog.Priority

					var priority syslog.Priority
					priority, err = ParseSyslogPriority(engine.Field.envValue)
					if err != nil {
						return err
					}

					engine.Field.int64Value = int64(priority)

				} else if engine.Field.value.Kind() == reflect.Int64 &&
					engine.Field.value.Type().String() == duration {
					// check if it is time.Duration

					engine.Field.durationValue, err = time.ParseDuration(engine.Field.envValue)
					if err != nil {
						return err
					}
					engine.Field.int64Value = engine.Field.durationValue.Nanoseconds()

				} else {
					// int

					engine.Field.int64Value, err = strconv.ParseInt(engine.Field.envValue, 10, 64)
					if err != nil {
						return err
					}
				}

				engine.Field.value.SetInt(engine.Field.int64Value)

			case reflect.Bool:
				engine.Field.value.SetBool(strings.ToLower(engine.Field.envValue) == "true")
			default:
				return unsupportedField(engine.Field.value.Type().Name())
			}
		}
	}

	if !nestedStruct {
		// we execute entire struct validation
		return engine.Validate.Struct(engine.Value.Interface())
	}

	return nil
}

// maximum returns the type's maximum value.
func maximum(kind reflect.Kind) uint64 {
	switch kind {
	case reflect.Uint8:
		return 255
	case reflect.Uint16:
		return 65535
	case reflect.Uint32:
		return 4294967295
	case reflect.Uint64, reflect.Uint:
		return 18446744073709551615
	default:
		return 0
	}
}
