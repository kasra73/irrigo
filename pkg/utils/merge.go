package utils

import (
	"irrigo/pkg/config"
	"reflect"
)

func MergeConfigs(globalConfig, localConfig config.GeneralSettings) config.GeneralSettings {
	// Deep clone the global config
	mergedConfig := globalConfig

	// Get the reflect.Value of the local and merged configs
	localVal := reflect.ValueOf(localConfig)
	mergedVal := reflect.ValueOf(&mergedConfig).Elem()

	// Iterate over the fields of the local config
	for i := 0; i < localVal.NumField(); i++ {
		// Get the field from the local and merged configs
		localField := localVal.Field(i)
		mergedField := mergedVal.Field(i)

		// If the local field is not the zero value for its type, override the merged field
		if localField.Interface() != reflect.Zero(localField.Type()).Interface() {
			mergedField.Set(localField)
		}
	}

	return mergedConfig
}
