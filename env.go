package env_loader

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kaatinga/assets"
)

// Deprecated: LoadUsingReflect loads a struct. The struct must contain tag 'env' on every struct field and must be set
// via pointer. It supports only byte and string field types.
func LoadUsingReflect(settings interface{}) error {

	// Getting the type
	t := reflect.TypeOf(settings).Elem()

	// Getting the value
	v := reflect.ValueOf(settings).Elem()

	// Reading the number of fields in the settings structure
	numberOfFields := t.NumField()
	if numberOfFields == 0 {
		return ErrTheModelHasNoFields
	}

	// temporary variables that are reused beneath
	var field, envTag, validateTag string

	validate := validator.New()

	for i := 0; i < numberOfFields; i++ {
		field = t.Field(i).Name // field name

		// getting the value of the tag 'env' for the field
		envTag = t.Field(i).Tag.Get("env")
		if envTag == "" {
			return errors.New(strings.Join([]string{"reading environment variables failed:", field}, " "))
		}

		envPar, ok := os.LookupEnv(envTag)
		if !ok {
			return errors.New(strings.Join([]string{"environment variable '", envTag, "' has not been found for the field '", field, "'"}, ""))
		}

		// getting the value of the tag 'validate' for the field
		validateTag = t.Field(i).Tag.Get("validate")

		switch v.Field(i).Kind() {
		case reflect.String:
			v.FieldByName(field).SetString(envPar)

			if validateTag != "" {
				if err := validate.Var(v.FieldByName(field).String(), validateTag); err != nil {
					return err
				}
			}
		case reflect.Uint8:
			value, ok := assets.StByte(envPar)
			if !ok {
				return errors.New(strings.Join([]string{"environment variable '", envTag, "' has been found but has incorrect value"}, ""))
			}

			v.FieldByName(field).SetUint(uint64(value))

			if validateTag != "" {
				if err := validate.Var(v.FieldByName(field).Uint(), validateTag); err != nil {
					return err
				}
			}
		default:
			return ErrUnsupportedField
		}
	}

	return nil
}
