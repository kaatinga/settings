package env_loader

import (
	"log/syslog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LoadUsingReflect loads a struct. The struct must contain tag 'env' on every struct field and must be set
// via pointer. It supports only byte and string field types.
func LoadUsingReflect(settings interface{}) error {

	engine, ok := settings.(*Engine)
	if !ok {
		engine = newEngine(settings)
	}

	err := engine.getStruct()
	if err != nil {
		return err
	}

	for i := 0; i < engine.NumberOfFields; i++ {
		engine.startIteration(i)

		//fmt.Println(engine.Loop.field.Name, "has toml tag:", engine.Loop.hasTomlTag)

		if engine.Loop.fieldValue.Kind() == reflect.Ptr ||
			engine.Loop.fieldValue.Kind() == reflect.Struct {
			// Мы проверяем а не вложенная ли это структура

			err = LoadUsingReflect(&Engine{
				Value: engine.Loop.fieldValue,
				Type:  engine.Loop.fieldValue.Type(),
			})
			if err != nil {
				return err
			}
			continue

		} else {

			// if a field has no env tag, we pass such a field
			if !engine.Loop.hasEnvTag {
				continue
			}

			if !engine.Value.Field(i).IsValid() {
				return ErrInternalFailure
			}

			// The struct must be send via pointer.
			if !engine.Value.Field(i).CanSet() {
				return ErrNotAddressable
			}

			//fmt.Println(engine.Value.Field(i).Type().String(), "can be changed")

			switch engine.Loop.fieldValue.Kind() {
			case reflect.String:
				engine.Loop.fieldValue.SetString(engine.Loop.envValue)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:

				if engine.Loop.fieldValue.Kind() == reflect.Uint32 &&
					engine.Loop.fieldValue.Type().String() == logLevel {
					// check if it is logrus.level

					var level logrus.Level
					level, err = logrus.ParseLevel(engine.Loop.envValue)
					if err != nil {
						return incorrectFieldValue(engine.Loop.envTag)
					}
					engine.Loop.uint64Value = uint64(level)

				} else {
					// uint

					engine.Loop.uint64Value, err = strconv.ParseUint(engine.Loop.envValue, 10, 64)
					if err != nil {
						return incorrectFieldValue(engine.Loop.envTag)
					}

					// check if whether the value exceeds the type maximum or not
					if engine.Loop.uint64Value > maximum(engine.Loop.fieldValue.Kind()) {
						return incorrectFieldValue(engine.Loop.envTag)
					}
				}

				engine.Loop.fieldValue.SetUint(engine.Loop.uint64Value)

			case reflect.Int64, reflect.Int:

				if engine.Loop.fieldValue.Kind() == reflect.Int &&
					engine.Loop.fieldValue.Type().String() == syslogPriority {
					// check if it is syslog.Priority

					var priority syslog.Priority
					priority, err = ParseSyslogPriority(engine.Loop.envValue)
					if err != nil {
						return incorrectFieldValue(engine.Loop.envTag)
					}

					engine.Loop.int64Value = int64(priority)

				} else if engine.Loop.fieldValue.Kind() == reflect.Int64 &&
					engine.Loop.fieldValue.Type().String() == duration {
					// check if it is time.Duration

					engine.Loop.durationValue, err = time.ParseDuration(engine.Loop.envValue)
					if err != nil {
						return err
					}
					engine.Loop.int64Value = engine.Loop.durationValue.Nanoseconds()

				} else {
					// int

					engine.Loop.int64Value, err = strconv.ParseInt(engine.Loop.envValue, 10, 64)
					if err != nil {
						return err
					}
				}

				engine.Loop.fieldValue.SetInt(engine.Loop.int64Value)

			case reflect.Bool:
				engine.Loop.fieldValue.SetBool(strings.ToLower(engine.Loop.envValue) == "true")
			default:
				return unsupportedField(engine.Loop.fieldValue.Type().Name())
			}

			//fmt.Println("установленное значение :", fieldValue.Interface())

			err = engine.validate()
			if err != nil {
				return err
			}
		}
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
