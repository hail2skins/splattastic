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
