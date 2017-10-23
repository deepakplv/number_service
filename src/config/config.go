/*
Package config contains all the configurations for the service
*/
package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

var config *viper.Viper

// Init is an exported method that takes the environment, starts the viper (external lib),
// and returns the configuration struct.
func Init(env string) {
	var err error
	v := viper.New()
	v.SetConfigType("toml")
	if env == "test" {
		v.AddConfigPath("../config/")
	} else {
		v.AddConfigPath("src/config/")
	}
	v.SetConfigName("base")
	err = v.MergeInConfig()
	if err != nil {
		log.Fatal("Error on parsing configuration file. Error " + err.Error())
	}

	v.SetConfigName(strings.ToLower(env))
	err = v.MergeInConfig()
	if err != nil {
		log.Fatal("Error on parsing configuration file. Error " + err.Error())
	}

	config = v
}

// GetConfig function to expose the config object
func GetConfig() *viper.Viper {
	return config
}