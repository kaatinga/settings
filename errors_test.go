package settings

import (
	"errors"
	"testing"
)

func TestUnsupportedFieldError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "basic field",
			field:    "PORT",
			expected: "environment variable 'PORT' has been found but the field type is unsupported",
		},
		{
			name:     "empty field",
			field:    "",
			expected: "environment variable '' has been found but the field type is unsupported",
		},
		{
			name:     "field with special chars",
			field:    "DB_URL",
			expected: "environment variable 'DB_URL' has been found but the field type is unsupported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unsupportedFieldError(tt.field)
			if err.Error() != tt.expected {
				t.Errorf("unsupportedFieldError.Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestIncorrectFieldValueError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "basic field",
			field:    "PORT",
			expected: "environment variable 'PORT' has been found but has incorrect value",
		},
		{
			name:     "empty field",
			field:    "",
			expected: "environment variable '' has been found but has incorrect value",
		},
		{
			name:     "field with special chars",
			field:    "DB_URL",
			expected: "environment variable 'DB_URL' has been found but has incorrect value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := incorrectFieldValueError(tt.field)
			if err.Error() != tt.expected {
				t.Errorf("incorrectFieldValueError.Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestValidationFailedError(t *testing.T) {
	tests := []struct {
		name           string
		fieldName      string
		fieldType      string
		validationRule string
		expected       string
	}{
		{
			name:           "basic validation",
			fieldName:      "Port",
			fieldType:      "int",
			validationRule: "required",
			expected:       "validation with rule 'required' failed on the field 'Port' of 'int' type",
		},
		{
			name:           "min validation",
			fieldName:      "CacheSize",
			fieldType:      "byte",
			validationRule: "min=10",
			expected:       "validation with rule 'min=10' failed on the field 'CacheSize' of 'byte' type",
		},
		{
			name:           "empty values",
			fieldName:      "",
			fieldType:      "",
			validationRule: "",
			expected:       "validation with rule '' failed on the field '' of '' type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &validationFailedError{
				Name:           tt.fieldName,
				Type:           tt.fieldType,
				ValidationRule: tt.validationRule,
			}
			if err.Error() != tt.expected {
				t.Errorf("validationFailedError.Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestUnsupportedFieldErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "same type",
			err:      unsupportedFieldError("PORT"),
			target:   unsupportedFieldError("DB"),
			expected: true,
		},
		{
			name:     "different type",
			err:      unsupportedFieldError("PORT"),
			target:   incorrectFieldValueError("PORT"),
			expected: false,
		},
		{
			name:     "nil target",
			err:      unsupportedFieldError("PORT"),
			target:   nil,
			expected: false,
		},
		{
			name:     "standard error",
			err:      unsupportedFieldError("PORT"),
			target:   errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, result, tt.expected)
			}
		})
	}
}

func TestIncorrectFieldValueErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "same type",
			err:      incorrectFieldValueError("PORT"),
			target:   incorrectFieldValueError("DB"),
			expected: true,
		},
		{
			name:     "different type",
			err:      incorrectFieldValueError("PORT"),
			target:   unsupportedFieldError("PORT"),
			expected: false,
		},
		{
			name:     "nil target",
			err:      incorrectFieldValueError("PORT"),
			target:   nil,
			expected: false,
		},
		{
			name:     "standard error",
			err:      incorrectFieldValueError("PORT"),
			target:   errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, result, tt.expected)
			}
		})
	}
}

func TestValidationFailedErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name: "same type",
			err: &validationFailedError{
				Name:           "Port",
				Type:           "int",
				ValidationRule: "required",
			},
			target: &validationFailedError{
				Name:           "DB",
				Type:           "string",
				ValidationRule: "min=5",
			},
			expected: true,
		},
		{
			name: "different type",
			err: &validationFailedError{
				Name:           "Port",
				Type:           "int",
				ValidationRule: "required",
			},
			target:   unsupportedFieldError("Port"),
			expected: false,
		},
		{
			name: "nil target",
			err: &validationFailedError{
				Name:           "Port",
				Type:           "int",
				ValidationRule: "required",
			},
			target:   nil,
			expected: false,
		},
		{
			name: "standard error",
			err: &validationFailedError{
				Name:           "Port",
				Type:           "int",
				ValidationRule: "required",
			},
			target:   errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, result, tt.expected)
			}
		})
	}
}

func TestNewUnsupportedFieldError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "basic field",
			field:    "PORT",
			expected: "environment variable 'PORT' has been found but the field type is unsupported",
		},
		{
			name:     "empty field",
			field:    "",
			expected: "environment variable '' has been found but the field type is unsupported",
		},
		{
			name:     "field with special chars",
			field:    "DB_URL",
			expected: "environment variable 'DB_URL' has been found but the field type is unsupported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUnsupportedFieldError(tt.field)
			if err.Error() != tt.expected {
				t.Errorf("NewUnsupportedFieldError(%v).Error() = %v, want %v", tt.field, err.Error(), tt.expected)
			}

			// Test that it returns the correct type
			if _, ok := err.(unsupportedFieldError); !ok {
				t.Errorf("NewUnsupportedFieldError(%v) should return unsupportedFieldError type", tt.field)
			}
		})
	}
}

func TestNewIncorrectFieldValueError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{
			name:     "basic field",
			field:    "PORT",
			expected: "environment variable 'PORT' has been found but has incorrect value",
		},
		{
			name:     "empty field",
			field:    "",
			expected: "environment variable '' has been found but has incorrect value",
		},
		{
			name:     "field with special chars",
			field:    "DB_URL",
			expected: "environment variable 'DB_URL' has been found but has incorrect value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewIncorrectFieldValueError(tt.field)
			if err.Error() != tt.expected {
				t.Errorf("NewIncorrectFieldValueError(%v).Error() = %v, want %v", tt.field, err.Error(), tt.expected)
			}

			// Test that it returns the correct type
			if _, ok := err.(incorrectFieldValueError); !ok {
				t.Errorf("NewIncorrectFieldValueError(%v) should return incorrectFieldValueError type", tt.field)
			}
		})
	}
}

func TestNewValidationFailedError(t *testing.T) {
	tests := []struct {
		name           string
		fieldName      string
		fieldType      string
		validationRule string
		expected       string
	}{
		{
			name:           "basic validation",
			fieldName:      "Port",
			fieldType:      "int",
			validationRule: "required",
			expected:       "validation with rule 'required' failed on the field 'Port' of 'int' type",
		},
		{
			name:           "min validation",
			fieldName:      "CacheSize",
			fieldType:      "byte",
			validationRule: "min=10",
			expected:       "validation with rule 'min=10' failed on the field 'CacheSize' of 'byte' type",
		},
		{
			name:           "empty values",
			fieldName:      "",
			fieldType:      "",
			validationRule: "",
			expected:       "validation with rule '' failed on the field '' of '' type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationFailedError(tt.fieldName, tt.fieldType, tt.validationRule)
			if err.Error() != tt.expected {
				t.Errorf("NewValidationFailedError(%v, %v, %v).Error() = %v, want %v",
					tt.fieldName, tt.fieldType, tt.validationRule, err.Error(), tt.expected)
			}

			// Test that it returns the correct type
			if _, ok := err.(*validationFailedError); !ok {
				t.Errorf("NewValidationFailedError(%v, %v, %v) should return *validationFailedError type",
					tt.fieldName, tt.fieldType, tt.validationRule)
			}
		})
	}
}

func TestErrorTypes(t *testing.T) {
	var _ error = unsupportedFieldError("test")
	var _ error = incorrectFieldValueError("test")
	var _ error = &validationFailedError{}
}

func TestErrorWrapping(t *testing.T) {
	// Test that errors can be wrapped and unwrapped properly
	originalErr := NewUnsupportedFieldError("PORT")
	wrappedErr := errors.Join(errors.New("wrapper"), originalErr)

	// Test that the original error can be found
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should be unwrappable to original error")
	}

	// Test error unwrapping
	var targetErr unsupportedFieldError
	if !errors.As(originalErr, &targetErr) {
		t.Error("Should be able to unwrap to unsupportedFieldError")
	}
	if targetErr != "PORT" {
		t.Errorf("Unwrapped error should have field 'PORT', got %v", targetErr)
	}
}
