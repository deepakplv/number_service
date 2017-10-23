package server

import (
	"config"
	"logger"
	"sync"
)

var onceRest sync.Once

// Init function to initialize the service
func Init() {
	onceRest.Do(func() {
		conf := config.GetConfig()
		logger.Log.Info("Initializing Rest server")
		r := NewRouter()
		if err := r.Start(conf.GetString("general.rest_server_port")); err != nil {
			logger.Log.Fatal("Unable to bring service up: " + err.Error())
		}

	})

}