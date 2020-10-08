package env_loader

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/kaatinga/assets"
)

// NewSettings creates a new settings struct
func NewSettings() {
	//TBI
}

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
		return errors.New("the input structure has not fields")
	}

	// временные переменные для цикла обработки параметров окружения ниже
	var field, tag string

	for i := 0; i < numberOfFields; i++ {
		field = t.Field(i).Name // имя поля

		// значение тега env
		tag = t.Field(i).Tag.Get("env")
		if tag == "" {
			return errors.New(strings.Join([]string{"Ошибка чтения параметра окружения:", field}, " "))
		}

		envPar, ok := os.LookupEnv(tag)
		if !ok {
			return errors.New(strings.Join([]string{"environment variable '", tag, "' has not been found for the field '", field, "'"}, ""))
		}

		switch v.Field(i).Kind() {
		case reflect.String:
			v.FieldByName(field).SetString(envPar)
		case reflect.Uint8:
			value, ok := assets.StByte(envPar)
			if !ok {
				return errors.New(strings.Join([]string{"environment variable '", tag, "' has been found but has incorrect value"}, ""))
			}

			v.FieldByName(field).SetUint(uint64(value))
		default:
			return errors.New(strings.Join([]string{"unsupported field type. only strings and bytes are supported"}, ""))
		}
	}

	return nil
}
