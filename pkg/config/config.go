package config

import (
	"github.com/spf13/cast"
	viperLib "github.com/spf13/viper"
	"gohub/pkg/helpers"
	"os"
)

// viper library example
var viper *viperLib.Viper

// ConfigFunc Dynamically load configuration information
type ConfigFunc func() map[string]any

// ConfigFuncs Load into this array first, loadConfig dynamically generates configuration information
var ConfigFuncs map[string]ConfigFunc

func init() {
	// Initialize the Viper library
	viper = viperLib.New()
	// Configuration type
	viper.SetConfigType("env")
	// Path to look for the environment variable file
	viper.AddConfigPath(".")
	// Set environment variable prefix
	viper.SetEnvPrefix("appEnv")
	// Read environment variables
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

// InitConfig Initialize configuration information,
// complete the loading of environment variables and config information
func InitConfig(env string) {
	// Load environment variables
	loadEnv(env)
	// Load config information
	loadConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

func loadEnv(envSuffix string) {
	// Load .env file by default, if there is a parameter --env=name, load the .env.name file
	envPath := ".env"
	if len(envSuffix) > 0 {
		filePath := ".env." + envSuffix
		if _, err := os.Stat(filePath); err == nil {
			// if it is .env.testing or .env.stage
			envPath = filePath
		}
	}

	// Load env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Monitor .env files and reload when changed
	viper.WatchConfig()
}

// Env Read environment variables, support default values
func Env(envName string, defaultValue ...any) any {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

// Add To add configuration items
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// Get To get configuration items
// param 'path' Dot passing is allowed, e.g.: app.name
func Get(path string, defaultValue ...any) string {
	return GetString(path, defaultValue...)
}

func internalGet(path string, defaultValue ...any) any {
	// config or environment variable does not exist
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString To get the configuration information of String type
func GetString(path string, defaultValue ...any) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt To get the configuration information of Int64 type
func GetInt(path string, defaultValue ...any) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetInt64 To get the configuration information of Int64 type
func GetInt64(path string, defaultValue ...any) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetFloat64 To get the configuration information of Float64 type
func GetFloat64(path string, defaultValue ...any) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetUint To get the configuration information of Uint type
func GetUint(path string, defaultValue ...any) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool To get the configuration information of Bool type
func GetBool(path string, defaultValue ...any) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString To get struct data
func GetStringMapString(path string) map[string]string {
	return cast.ToStringMapString(path)
}
