package env_loader

import (
	"github.com/go-playground/validator/v10"
	"github.com/kaatinga/assets"
	"os"
	"reflect"
)

// LoadUsingReflect loads a struct. The struct must contain tag 'env' on every struct field and must be set
// via pointer. It supports only byte and string field types.
func LoadUsingReflect(settings interface{}) error {

	var settingsEngine *SettingsStruct

	var ok bool
	if settingsEngine, ok = settings.(*SettingsStruct); !ok {
		//fmt.Println("this is root struct")
		settingsEngine = newSettingsStruct(settings)
	}

	settingsEngine.getStruct()

	// the main model must be a struct
	if settingsEngine.Type.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	// Reading the number of fields in the settings structure
	numberOfFields := settingsEngine.Type.NumField()
	if numberOfFields == 0 {
		return ErrTheModelHasNoFields
	}

	// temporary variables that are reused beneath
	var envTag, envPar, validateTag string
	var fieldType reflect.StructField
	var fieldValue reflect.Value

	validate := validator.New()

	for i := 0; i < numberOfFields; i++ {
		fieldType = settingsEngine.Type.Field(i)
		fieldValue = settingsEngine.Value.FieldByName(fieldType.Name)

		//fmt.Println(fieldValue.Type().Name())
		//fmt.Println("reflect field check passed", fieldType.Name)

		// getting the fieldValue of the tag 'env' for the field
		envTag, ok = fieldType.Tag.Lookup("env")
		if !ok {

			if fieldValue.Kind() == reflect.Ptr {
				err := LoadUsingReflect(&SettingsStruct{
					Value: fieldValue,
					Type:  fieldValue.Type(),
				})
				if err != nil {
					return err
				}
				continue
			}
			//fmt.Println("the field", fieldType.Name, "will be omitted as it has no 'env' tag")
			continue

		} else {

			//fmt.Println("envTag", envTag)
			envPar, ok = os.LookupEnv(envTag)
			if !ok {
				return EnvironmentVariableNotFound(envTag)
			}

			// getting the fieldValue of the tag 'validate' for the field
			validateTag = fieldType.Tag.Get("validate")

			if settingsEngine.Value.Field(i).IsValid() {
				// The struct must be send via pointer.
				if !settingsEngine.Value.Field(i).CanSet() {
					return ErrNotAddressable
				}
			} else {
				return ErrInternalFailure
			}

			switch fieldValue.Kind() {
			case reflect.String:

				fieldValue.SetString(envPar)

				if validateTag != "" {
					if err := validate.Var(fieldValue.String(), validateTag); err != nil {
						return err
					}
				}
			case reflect.Uint8:
				var byteValue byte
				byteValue, ok = assets.StByte(envPar)
				if !ok {
					return IncorrectFieldValue(envTag)
				}

				fieldValue.SetUint(uint64(byteValue))

				if validateTag != "" {
					if err := validate.Var(fieldValue.Uint(), validateTag); err != nil {
						return err
					}
				}
			default:
				return ErrUnsupportedField
			}
		}
	}

	return nil
}
