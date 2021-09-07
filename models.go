package settings

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
	envTag            string
	envValue          string
	validationRule    string
	durationValue     time.Duration
	int64Value        int64
	uint64Value       uint64
	field             reflect.StructField
	value             reflect.Value
	hasEnvTag         bool
	mustBeOmitted     bool
	hasEnvValue       bool
	mustBeValidated   bool
	required          bool
	defaultSetting    string
	hasDefaultSetting bool
	float64Value      float64
}

// exceedsMaximumUint возвращает максимальное значение типов uint.
func (field *Loop) exceedsMaximumUint() bool {
	kind := field.value.Kind()
	if kind == reflect.Uint { // TODO: uint не всегда = uint64
		kind = reflect.Uint64
	}

	return field.uint64Value > 1<<(2<<(uint64(kind)-6))-1
}

// notInIntRange возвращает максимальное и минимальное значение типов int.
func (field *Loop) notInIntRange() bool {
	kind := field.value.Kind()
	if kind == reflect.Int { // TODO: int не всегда = int64
		kind = reflect.Int64
	}

	return field.int64Value > 1<<((2<<(int64(kind)-1))-1)-1 ||
		field.int64Value < -1<<((2<<(int64(kind)-1))-1)
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
		return ErrTheModelHasEmptyStruct
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

func (engine *Engine) validateRequired() {

	engine.Field.required = false

	// receiving the 'validate' tag value
	engine.Field.validationRule, engine.Field.mustBeValidated = engine.Field.field.Tag.Lookup("validate")

	// process validation rule to ascertain the required status
	rules := strings.Split(engine.Field.validationRule, ",")
	for _, value := range rules {
		if value == required {
			engine.Field.required = true
			break
		}
	}
}

// validate validates the current value using `validate` tag.
//func (engine *Engine) validate() error {
//
//	if engine.Field.mustBeValidated &&
//		engine.Validate.Var(engine.Field.value.Interface(), engine.Field.validationRule) != nil {
//		return engine.validationFailed()
//	}
//
//	return nil
//}

// startIteration launches field processing.
func (engine *Engine) startIteration(i int) {
	engine.Field.field = engine.Type.Field(i)
	engine.Field.value = engine.Value.FieldByName(engine.Field.field.Name)

	// receiving env tag
	engine.Field.envTag, engine.Field.hasEnvTag = engine.Field.field.Tag.Lookup(env)
	if engine.Field.hasEnvTag && engine.Field.envTag == omit {
		engine.Field.mustBeOmitted = true
		return
	}
	engine.Field.mustBeOmitted = false

	// receiving default setting
	engine.Field.defaultSetting, engine.Field.hasDefaultSetting = engine.Field.field.Tag.Lookup(defaultSetting)
}
