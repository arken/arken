package config

import (
    "fmt"
    "os"
    "reflect"
    "strings"
)

//Looks for discrepancies between environment variables and the internal config
//struct, preferring the value set in the environment variable.
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
        var evName string
        var fieldName string
        typeOfT := iter.Type()
        for j := 0; j < iter.NumField(); j++ {
            f := iter.Field(j)
            fieldName = typeOfT.Field(j).Name
            evName = strings.ToUpper(fieldName)
            evName = evPrefix + evName
            evVal, evExists := os.LookupEnv(evName)
            if evExists && evVal != f.String() {
                iter.FieldByName(fieldName).SetString(evVal)
                fmt.Printf("Env. var. \"%v\" does not match internal" +
                    " memory. Updating memory value to \"%v\".\n", evName, evVal)
            } else if !evExists {
                fmt.Printf("Env. var \"%v\" not found, using \"%v\".\n", evName, iter.FieldByName(fieldName).String())
            }
        }
    }
}
