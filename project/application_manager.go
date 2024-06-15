package project

import (
	"log"
)

type ApplicationManager struct {
	Config ApplicationConfig
}

func NewApplicationManager() *ApplicationManager {
	configLocation, err := getConfigLocation()
	if err != nil {
		log.Fatalf("Could not find/create the config file: %v", err)
	}

	config, err := loadConfig(configLocation)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	return &ApplicationManager{
		Config: config,
	}
}

func (a *ApplicationManager) GetConfig() ApplicationConfig {
	return a.Config
}
