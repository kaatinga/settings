package settings

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Load loads settings to a struct.
func Load(settings any) error {
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

		// passing the omit-tagged fields
		if engine.Field.mustBeOmitted {
			continue
		}

		if engine.Field.value.Kind() == reflect.Ptr || //nolint:nestif
			engine.Field.value.Kind() == reflect.Struct {
			// we check whether the field is pointer or struct

			err = Load(&Engine{
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

			// if a field has env tag, but the env was not found, and if it is required
			// we return error
			engine.Field.envValue, engine.Field.hasEnvValue = os.LookupEnv(engine.Field.envTag)
			if !engine.Field.hasEnvValue {
				if engine.Field.hasDefaultSetting {
					// substitute the envValue with default setting
					engine.Field.envValue = engine.Field.defaultSetting
				} else {
					if engine.Field.required {
						return engine.validationFailed()
					}
					continue
				}
			}

			if !engine.Value.Field(i).IsValid() {
				return ErrInternalFailure
			}

			// We are checking if the field is addressable
			if !engine.Value.Field(i).CanSet() {
				return ErrNotAddressableField
			}

			switch engine.Field.value.Kind() { //nolint:exhaustive
			case reflect.Slice:
				switch engine.Field.value.Type().Elem().Kind() {
				case reflect.Uint8:
					engine.Field.value.SetBytes([]byte(engine.Field.envValue))
				case reflect.String:
					stringSlice := strings.Split(engine.Field.envValue, ",")
					engine.Field.value.Set(reflect.ValueOf(stringSlice))
				default:
					return unsupportedFieldError(engine.Field.envTag)
				}
			case reflect.String:
				engine.Field.value.SetString(engine.Field.envValue)
			case reflect.Float64:
				engine.Field.float64Value, err = strconv.ParseFloat(engine.Field.envValue, 64)
				if err != nil {
					return incorrectFieldValueError(engine.Field.envTag)
				}

				engine.Field.value.SetFloat(engine.Field.float64Value)

			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				engine.Field.uint64Value, err = strconv.ParseUint(engine.Field.envValue, 10, 64)
				if err != nil {
					return incorrectFieldValueError(engine.Field.envTag)
				}

				// check if whether the value exceeds the type maximum or not
				if engine.Field.exceedsMaximumUint() {
					return incorrectFieldValueError(engine.Field.envTag)
				}

				engine.Field.value.SetUint(engine.Field.uint64Value)

			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				if engine.Field.value.Kind() == reflect.Int64 &&
					engine.Field.value.Type().String() == duration {
					// check if it is time.Duration

					engine.Field.durationValue, err = time.ParseDuration(engine.Field.envValue)
					if err != nil {
						return err
					}
					engine.Field.int64Value = engine.Field.durationValue.Nanoseconds()
				} else {
					engine.Field.int64Value, err = strconv.ParseInt(engine.Field.envValue, 10, 64)
					if err != nil {
						return err
					}

					if engine.Field.notInIntRange() {
						return incorrectFieldValueError(engine.Field.envTag)
					}
				}

				engine.Field.value.SetInt(engine.Field.int64Value)

			case reflect.Bool:
				engine.Field.value.SetBool(strings.ToLower(engine.Field.envValue) == "true")
			default:
				return unsupportedFieldError(engine.Field.value.Type().Name())
			}
		}
	}

	if !nestedStruct {
		if err = engine.Validate.Struct(engine.Value.Interface()); err != nil {
			return err
		}

		if err = runCustomValidation(engine, err); err != nil {
			return err
		}
	}

	return nil
}

func runCustomValidation(engine *Engine, err error) error {
	if to, validatable := engine.Value.Interface().(interface {
		Validate() error
	}); validatable {
		if err = to.Validate(); err != nil {
			return err
		}
	}
	return nil
}
