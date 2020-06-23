package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// ConsolidateEnvVars looks for discrepancies between environment variables and the
// internal config struct, preferring the value set in the environment variable.
// In this function, variables starting with "field" track values associated
// with the struct, and a "ev" prefix indicates the value is associated with the
// environment variable.
func ConsolidateEnvVars(conf *Config) {
	numSubStructs := reflect.ValueOf(conf).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		var iter reflect.Value
		var evPrefix string
		if i == 0 {
			iter = reflect.ValueOf(&conf.Sources).Elem()
			evPrefix = "ARKEN_SOURCES_"
		} else if i == 1 {
			iter = reflect.ValueOf(&conf.Database).Elem()
			evPrefix = "ARKEN_DB_"
		} else {
			iter = reflect.ValueOf(&conf.General).Elem()
			evPrefix = "ARKEN_GENERAL_"
		}
		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			if fieldVal != "Version" {
				fieldName := structType.Field(j).Name
				evName := evPrefix + strings.ToUpper(fieldName)
				evVal, evExists := os.LookupEnv(evName)
				if evExists && evVal != fieldVal {
					iter.FieldByName(fieldName).SetString(evVal)
					fmt.Printf("Env. var. \"%v\" does not match internal"+
						" memory. Updating memory value to \"%v\".\n", evName, evVal)
				}
			}
		}
	}
}
