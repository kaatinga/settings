package env_loader

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"time"
)

// Engine — data model to process settings.
type Engine struct {
	Value          reflect.Value
	Type           reflect.Type
	NumberOfFields int
	Validate       *validator.Validate
	Loop           Loop
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
	fieldValue      reflect.Value
	hasEnvTag       bool
	mustBeValidated bool
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
		Name:           engine.Loop.field.Name,
		Type:           engine.Loop.fieldValue.Type().String(),
		ValidationRule: engine.Loop.validationRule,
	}
}

// validate validates the current value using `validate` tag.
func (engine *Engine) validate() error {

	// получаем значение тега 'validate' для поля
	engine.Loop.validationRule, engine.Loop.mustBeValidated = engine.Loop.field.Tag.Lookup("validate")
	if engine.Loop.mustBeValidated {
		err := engine.Validate.Var(engine.Loop.fieldValue.Interface(), engine.Loop.validationRule)
		if err != nil {
			//fmt.Println(err)
			return engine.validationFailed()
		}
	}

	return nil
}

// startIteration launches field processing.
func (engine *Engine) startIteration(i int) {
	engine.Loop.field = engine.Type.Field(i)
	engine.Loop.fieldValue = engine.Value.FieldByName(engine.Loop.field.Name)

	// receiving env tag
	engine.Loop.envTag, engine.Loop.hasEnvTag = engine.Loop.field.Tag.Lookup(env)
	//fmt.Println(engine.Loop.envTag, engine.Loop.hasEnvTag)
}
