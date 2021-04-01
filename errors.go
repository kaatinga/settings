package env_loader

import cer "github.com/kaatinga/const-errs"

const (
	ErrUnsupportedField    cer.Error = "unsupported field type. only strings and bytes are supported"
	ErrTheModelHasNoFields cer.Error = "the input structure has not fields"
)
