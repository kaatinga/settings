package settings

import (
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Engine — data model to process settings.
type Engine struct {
	Type           reflect.Type
	Validate       *validator.Validate
	Value          reflect.Value
	Field          Loop
	NumberOfFields int
}

// newEngine creates new model to process settings.
func newEngine(settings any) *Engine {
	return &Engine{
		Value:    reflect.ValueOf(settings),
		Type:     reflect.TypeOf(settings),
		Validate: validator.New(),
	}
}

// Loop — variables that used during field processing.
type Loop struct {
	value             reflect.Value
	envTag            string
	envValue          string
	validationRule    string
	defaultSetting    string
	field             reflect.StructField
	durationValue     time.Duration
	int64Value        int64
	uint64Value       uint64
	float64Value      float64
	hasEnvTag         bool
	mustBeOmitted     bool
	hasEnvValue       bool
	mustBeValidated   bool
	required          bool
	hasDefaultSetting bool
}

// exceedsMaximumUint returns true if the value exceeds the maximum uint value.
func (field *Loop) exceedsMaximumUint() bool {
	var kind reflect.Kind
	if kind = field.value.Kind(); kind == reflect.Uint {
		kind = reflect.Uint64
	}

	//nolint:gomnd // it is a formula
	return field.uint64Value > 1<<(2<<(uint64(kind)-6))-1
}

// notInIntRange returns true if the value is not in the int range.
func (field *Loop) notInIntRange() bool {
	kind := field.value.Kind()
	var minimum, maximum int64

	switch kind {
	case reflect.Int8:
		minimum, maximum = math.MinInt8, math.MaxInt8
	case reflect.Int16:
		minimum, maximum = math.MinInt16, math.MaxInt16
	case reflect.Int32:
		minimum, maximum = math.MinInt32, math.MaxInt32
	case reflect.Int, reflect.Int64:
		minimum, maximum = math.MinInt64, math.MaxInt64
	default:
		return false
	}

	return field.int64Value > maximum || field.int64Value < minimum
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

// validationFailedError forms validation error.
func (engine *Engine) validationFailed() error {
	return &validationFailedError{
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
