package logger

import (
	"config"
	"github.com/plivo/go-plivolog/plivolog"
)

// Log - log object
var Log *plivolog.PlivoLogger

// Init function for initializing the logger
func Init() {
	conf := config.GetConfig()
	var err error

	Log, err = plivolog.New()
	if conf.GetString("general.config") == "test" {
		Log.Info("For testing environment err will be ignored")
		return
	}
	if err != nil {
		panic("logger could not be initialized. Error " + err.Error())
	}
}
