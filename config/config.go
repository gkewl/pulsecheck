package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gkewl/pulsecheck/logger"
)

const (
	PULSE_APP_TIMEOUT            = "PULSE_APP_TIMEOUT"
	PULSE_CLIENT_TIMEOUT         = "PULSE_CLIENT_TIMEOUT"
	PULSE_TOKEN_TIMEOUT_IN_HOURS = "PULSE_TOKEN_TIMEOUT_IN_HOURS"
	PULSE_PRIVATE_KEY_SECRET     = "PULSE_PRIVATE_KEY_SECRET"
	PULSE_PUBLIC_KEY_SECRET      = "PULSE_PUBLIC_KEY_SECRET"

	PULSE_DB_TYPE            = "PULSE_DB_TYPE"
	PULSE_DB_USER            = "PULSE_DB_USER"
	PULSE_DB_NET_PROTOCOL    = "PULSE_DB_NET_PROTOCOL"
	PULSE_DB_ADDRESS         = "PULSE_DB_ADDRESS"
	PULSE_DB_TIMEOUT         = "PULSE_DB_TIMEOUT"
	PULSE_DB_DBNAME          = "PULSE_DB_DBNAME"
	PULSE_DB_MAX_IDLE_CONNS  = "PULSE_DB_MAX_IDLE_CONNS"
	PULSE_DB_MAX_OPEN_CONNS  = "PULSE_DB_MAX_OPEN_CONNS"
	PULSE_DB_PASSWORD_SECRET = "PULSE_DB_PASSWORD_SECRET"
	PULSE_DB_CONN_LIFETIME   = "PULSE_DB_CONN_LIFETIME"
)

const hostname = "__hostname__"

var configMap map[string]string
var missingKeys []string
var initialized = false

func LoadConfigurations() {
	if initialized {
		logger.LogError("Initializing the app more than once; probably a bug!", "N/A")
	} else {
		configMap = make(map[string]string)

		var err error
		configMap[hostname], err = os.Hostname()
		if err != nil {
			configMap[hostname] = "localhost"
		}

		setConfigValue(PULSE_APP_TIMEOUT)
		setConfigValue(PULSE_PRIVATE_KEY_SECRET)
		setConfigValue(PULSE_PUBLIC_KEY_SECRET)
		setConfigValue(PULSE_CLIENT_TIMEOUT)
		setConfigValue(PULSE_TOKEN_TIMEOUT_IN_HOURS)

		setConfigValue(PULSE_DB_PASSWORD_SECRET)
		setConfigValue(PULSE_DB_TYPE)
		setConfigValue(PULSE_DB_TIMEOUT)
		setConfigValue(PULSE_DB_USER)
		setConfigValue(PULSE_DB_NET_PROTOCOL)
		setConfigValue(PULSE_DB_ADDRESS)
		setConfigValue(PULSE_DB_DBNAME)
		setConfigValue(PULSE_DB_MAX_IDLE_CONNS)
		setConfigValue(PULSE_DB_MAX_OPEN_CONNS)
		setConfigValue(PULSE_DB_CONN_LIFETIME)

		if len(missingKeys) > 0 {
			// Do NOT proceed to start the app if any required config vars are unset.
			fmt.Println("Missing Environment Variables: " + strings.Join(missingKeys[:], ", "))
			panic("Missing Environment Variables: " + strings.Join(missingKeys[:], ", "))

		} else {
			initialized = true
		}
	}
}

func setConfigValue(configKey string) {
	if configValue := os.Getenv(configKey); configValue != "" {
		configMap[configKey] = configValue
	} else {
		missingKeys = append(missingKeys, configKey)
	}
}

func setOptionalConfigValue(configKey string) {
	if configValue := os.Getenv(configKey); configValue != "" {
		configMap[configKey] = configValue
	}
}

// GetEnv gets a config key from the environment loaded at the start of the program.
// Instead of this, consider making a config module function to retrieve your value
func GetEnv(configKey string) string {
	configValue, isPresent := configMap[configKey]
	if !isPresent {
		logger.LogError(fmt.Sprintf("Accessing an unset configuration key: %s; probably a bug!", configKey), "N/A")
	}
	return configValue
}

// GetAltEnv gets a numbered 'alternate' config variable, for alternate
// configuration variables identified with a suffix
func LookupAltEnv(configKey string, alt ...string) (string, bool) {
	if len(alt) == 0 {
		configValue, isPresent := configMap[configKey]
		return configValue, isPresent
	}
	envKey := strings.Join(append([]string{configKey}, alt...), "_")
	return os.LookupEnv(envKey)
}

// GetEnvOrAlt gets configuration, defaulting to the "primary" version
// but preferring the alternate.
func GetEnvOrAlt(configKey string, alt ...string) string {
	altVal, ok := LookupAltEnv(configKey, alt...)
	if ok {
		return altVal
	}
	return GetEnv(configKey)
}

// GetOptionalEnvOrAlt gets configuration, defaulting to the "primary" version
// but preferring the alternate.
func GetOptionalEnvOrAlt(configKey string, alt ...string) string {
	altVal, ok := LookupAltEnv(configKey, alt...)
	if ok {
		return altVal
	}
	return configMap[configKey]
}

// HermesEnabled returns true if Hermes is configured for the alt env.  Does not default;
// the minimum to get an extra (identical) connection to Hermes is to set PULSE_HERMES_CONN_ENABLED_2=yes

func toBool(envVal string, defaultVal bool) bool {
	var retVal = defaultVal
	switch strings.ToLower(envVal) {
	case "true", "t", "y", "yes", "on", "1":
		retVal = true
	case "false", "f", "n", "no", "off", "0":
		retVal = false
	}
	return retVal
}
