package config

import (
	"strconv"
)

//EnvVariableLookuper provides an interface to fetch env variable values
type EnvVariableLookuper interface {
	GetEnv(configKey string) string
}

//TestingEnvVariableLookup used by our tests to mock env variable fetching
var TestingEnvVariableLookup EnvVariableLookuper

//OsEnvVariableLookuper implements interface EnvVariableLookuper and performs lookup
//from configuration
type OsEnvVariableLookuper struct{}

//GetEnv wraps the config.GetEnv variable to be able to mock it for testing
//if TestingEnvVariableLookup != null returns mock otherwise returns real
func (lkp OsEnvVariableLookuper) GetEnv(envVarName string) string {
	if TestingEnvVariableLookup != nil {
		return TestingEnvVariableLookup.GetEnv(envVarName)
	}

	return GetEnv(envVarName)
}

//GetEnvValueFromOs calls OsEnvVariableLookuper
func GetEnvValueFromOs(envVarName string) string {
	return OsEnvVariableLookuper{}.GetEnv(envVarName)
}

//GetEnvVariableIntValue casts an env variable value to int64
func GetEnvVariableIntValue(envVarName string, defValue int64) int64 {
	envVal := GetEnvValueFromOs(envVarName)

	retVal, err := strconv.ParseInt(envVal, 10, 64)

	if err != nil {
		return defValue
	}

	return retVal
}
