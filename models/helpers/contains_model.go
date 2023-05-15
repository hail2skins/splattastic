package helpers

import (
	"reflect"
)

// ContainsModel checks if a slice of models contains a model with the given name.
func ContainsModel(models []interface{}, name string) bool {
	for _, m := range models {
		modelValue := reflect.ValueOf(m)
		modelName := modelValue.FieldByName("Name")
		if modelName.IsValid() && modelName.String() == name {
			return true
		}
	}
	return false
}

// ContainsModel checks if a slice of models contains a model with the given name.
func ContainsModelDescription(models []interface{}, description string) bool {
	for _, m := range models {
		modelValue := reflect.ValueOf(m)
		modelName := modelValue.FieldByName("Description")
		if modelName.IsValid() && modelName.String() == description {
			return true
		}
	}
	return false
}
