package env_loader

import "reflect"

// SettingsStruct — data model to process settings.
type SettingsStruct struct {
	Value reflect.Value
	Type  reflect.Type
}

// newSettingsStruct creates new model to process settings.
func newSettingsStruct(settings interface{}) *SettingsStruct {
	return &SettingsStruct{
		reflect.ValueOf(settings),
		reflect.TypeOf(settings),
	}
}

// getStruct returns the final struct to load settings.
func (settingsStruct *SettingsStruct) getStruct() {
	if settingsStruct.Value.Kind() == reflect.Ptr {
		settingsStruct.Value = settingsStruct.Value.Elem()
		settingsStruct.Type = settingsStruct.Type.Elem()
		settingsStruct.getStruct()
	}
}
