package env_loader

import cer "github.com/kaatinga/const-errs"

const (
	ErrUnsupportedField    cer.Error = "unsupported field type. only strings and bytes are supported"
	ErrTheModelHasNoFields cer.Error = "the input structure has not fields"
	ErrNotAStruct          cer.Error = "the configuration must be a struct"
	ErrNotAddressable      cer.Error = "the struct and every struct field must be added via pointer"

	ErrInternalFailure             cer.Error = "an internal package error"
	ErrIncorrectFieldValue         cer.Error = "environment variable has been found but has incorrect value"
	ErrEnvironmentVariableNotFound cer.Error = "environment variable has not been found"
)

type IncorrectFieldValue string

func (err IncorrectFieldValue) Error() string {
	return "environment variable " + string(err) + " has been found but has incorrect value"
}

func (err IncorrectFieldValue) Is(target error) bool {
	if target == ErrIncorrectFieldValue {
		return true
	}

	return false
}

type EnvironmentVariableNotFound string

func (err EnvironmentVariableNotFound) Error() string {
	return "environment variable " + string(err) + " has not been found"
}

func (err EnvironmentVariableNotFound) Is(target error) bool {
	if target == ErrEnvironmentVariableNotFound {
		return true
	}

	return false
}
