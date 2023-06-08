package settings

import cer "github.com/kaatinga/const-errs"

const (
	ErrUnsupportedField       cer.Error = "unsupported field type"
	ErrTheModelHasEmptyStruct cer.Error = "an input struct has no fields"
	ErrNotAStruct             cer.Error = "the configuration must be a struct"
	ErrNotAddressable         cer.Error = "the main struct must be pointed out via pointer"
	ErrNotAddressableField    cer.Error = "the value is not addressable or main struct is not indicated via pointer"

	ErrInternalFailure     cer.Error = "an internal package error"
	ErrIncorrectFieldValue cer.Error = "variable has been found but has incorrect value"
	ErrValidationFailed    cer.Error = "field validation failed"
)

type unsupportedFieldError string

func (err unsupportedFieldError) Error() string {
	return "environment variable " + string(err) + " has been found but the field type is unsupported"
}

func (err unsupportedFieldError) Is(target error) bool {
	return target == ErrUnsupportedField //nolint:goerr113
}

type incorrectFieldValueError string

func (err incorrectFieldValueError) Error() string {
	return "environment variable " + string(err) + " has been found but has incorrect value"
}

func (err incorrectFieldValueError) Is(target error) bool {
	return target == ErrIncorrectFieldValue //nolint:goerr113
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
	return target == ErrValidationFailed //nolint:goerr113
}
