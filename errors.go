package env_loader

import cer "github.com/kaatinga/const-errs"

const (
	ErrUnsupportedField    cer.Error = "unsupported field type"
	ErrTheModelHasNoFields cer.Error = "the input structure has no fields"
	ErrNotAStruct          cer.Error = "the configuration must be a struct"
	ErrNotAddressable      cer.Error = "the main struct must be pointed out via pointer"

	ErrInternalFailure     cer.Error = "an internal package error"
	ErrIncorrectFieldValue cer.Error = "variable has been found but has incorrect value"
	ErrValidationFailed    cer.Error = "field validation failed"
	ErrIncorrectPriority   cer.Error = "incorrect syslog priority"
)

type unsupportedField string

func (err unsupportedField) Error() string {
	return "environment variable " + string(err) + " has been found but has incorrect value"
}

func (err unsupportedField) Is(target error) bool {
	return target == ErrUnsupportedField
}

type incorrectFieldValue string

func (err incorrectFieldValue) Error() string {
	return "environment variable " + string(err) + " has been found but has incorrect value"
}

func (err incorrectFieldValue) Is(target error) bool {
	return target == ErrIncorrectFieldValue
}

type validationFailed struct {
	Name           string
	Type           string
	ValidationRule string
}

func (err *validationFailed) Error() string {
	return "validation with rule '" + err.ValidationRule + "' failed on the field '" + err.Name + "' of '" + err.Type + "' type"
}

func (err *validationFailed) Is(target error) bool {
	return target == ErrValidationFailed
}

// incorrectPriority — error for ParseSyslogPriority()
type incorrectPriority string

func (err incorrectPriority) Error() string {
	return "syslog priority " + string(err) + " is incorrect"
}

func (err incorrectPriority) Is(target error) bool {
	return target == ErrIncorrectPriority
}
