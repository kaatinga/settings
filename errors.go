package settings

import "errors"

var (
	ErrTheModelHasEmptyStruct = errors.New("an input struct has no fields")
	ErrNotAStruct             = errors.New("the configuration must be a struct")
	ErrNotAddressable         = errors.New("the main struct must be pointed out via pointer")
	ErrNotAddressableField    = errors.New("the value is not addressable or main struct is not indicated via pointer")
	ErrInternalFailure        = errors.New("an internal package error")
)

type unsupportedFieldError string

func (err unsupportedFieldError) Error() string {
	return "environment variable '" + string(err) + "' has been found but the field type is unsupported"
}

func (err unsupportedFieldError) Is(target error) bool {
	_, ok := target.(unsupportedFieldError)
	return ok
}

type incorrectFieldValueError string

func (err incorrectFieldValueError) Error() string {
	return "environment variable '" + string(err) + "' has been found but has incorrect value"
}

func (err incorrectFieldValueError) Is(target error) bool {
	_, ok := target.(incorrectFieldValueError)
	return ok
}

type validationFailedError struct {
	Name           string
	Type           string
	ValidationRule string
}

func (err *validationFailedError) Error() string {
	return "validation with rule '" + err.ValidationRule + "' failed on the field '" + err.Name + "' of '" + err.Type + "' type"
}

func (err *validationFailedError) Is(target error) bool {
	_, ok := target.(*validationFailedError)
	return ok
}

func NewUnsupportedFieldError(fieldName string) error {
	return unsupportedFieldError(fieldName)
}

func NewIncorrectFieldValueError(fieldName string) error {
	return incorrectFieldValueError(fieldName)
}

func NewValidationFailedError(name, fieldType, validationRule string) error {
	return &validationFailedError{
		Name:           name,
		Type:           fieldType,
		ValidationRule: validationRule,
	}
}
