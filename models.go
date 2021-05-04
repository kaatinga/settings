package env_loader

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"time"
)

// Engine — data model to process settings.
type Engine struct {
	Value          reflect.Value
	Type           reflect.Type
	NumberOfFields int
	Validate       *validator.Validate
	Field          Loop
}

// newEngine creates new model to process settings.
func newEngine(settings interface{}) *Engine {
	return &Engine{
		Value:    reflect.ValueOf(settings),
		Type:     reflect.TypeOf(settings),
		Validate: validator.New(),
	}
}

// Loop — variables that used during field processing.
type Loop struct {
	envTag          string
	envValue        string
	validationRule  string
	tmpStringValue  string
	durationValue   time.Duration
	int64Value      int64
	uint64Value     uint64
	field           reflect.StructField
	value           reflect.Value
	hasEnvTag       bool
	mustBeValidated bool
	required        bool
}

// getStruct checks and returns a struct to process.
func (engine *Engine) getStruct() error {
	if engine.Value.Kind() == reflect.Ptr {

		// initing struct if it is added via pointer
		if engine.Value.IsNil() {

			// the main struct must be inited
			// the third law of reflection
			// https://blog.golang.org/laws-of-reflection
			// we have no parameter to change value
			if !engine.Value.CanSet() {
				return ErrNotAddressable
			}

			newStruct := reflect.New(engine.Type.Elem())
			newValue := reflect.ValueOf(newStruct.Interface())
			engine.Value.Set(newValue)
		}

		engine.Value = engine.Value.Elem()
		engine.Type = engine.Type.Elem()

		err := engine.getStruct()
		if err != nil {
			return err
		}
	}

	// checking that kind us a struct
	if engine.Type.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	// checking the number of the fields in the struct.
	engine.NumberOfFields = engine.Type.NumField()
	if engine.NumberOfFields == 0 {
		return ErrTheModelHasNoFields
	}

	return nil
}

// validationFailed forms validation error.
func (engine *Engine) validationFailed() error {
	return &validationFailed{
		Name:           engine.Field.field.Name,
		Type:           engine.Field.value.Type().String(),
		ValidationRule: engine.Field.validationRule,
	}
}

func (engine *Engine) prevalidate() {

	// receiving the 'validate' tag value
	engine.Field.validationRule, engine.Field.mustBeValidated = engine.Field.field.Tag.Lookup("validate")

	// process validation rule to ascertain the required status
	if strings.Contains(engine.Field.validationRule, required) {
		engine.Field.required = true
	}
}

// validate validates the current value using `validate` tag.
func (engine *Engine) validate() error {

	if engine.Field.mustBeValidated {
		err := engine.Validate.Var(engine.Field.value.Interface(), engine.Field.validationRule)
		if err != nil {
			//fmt.Println(err)
			return engine.validationFailed()
		}
	}

	return nil
}

// startIteration launches field processing.
func (engine *Engine) startIteration(i int) {
	engine.Field.field = engine.Type.Field(i)
	engine.Field.value = engine.Value.FieldByName(engine.Field.field.Name)

	// receiving env tag
	engine.Field.envTag, engine.Field.hasEnvTag = engine.Field.field.Tag.Lookup(env)
	//fmt.Println(engine.Field.envTag, engine.Field.hasEnvTag)
}
