package project

import (
	"fmt"
	"log"
	"os"
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

func (a *ApplicationManager) writeConfigToFile() error {
	configLocation, err := getConfigLocation()
	if err != nil {
		return fmt.Errorf("cannot find application config file: %v", err)
	}

	file, err := os.OpenFile(configLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot open config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	err = writeStructToFile(a.Config, file)
	if err != nil {
		return err
	}
	return nil
}

func (a *ApplicationManager) writeProjectDetailsToAppConfig(projectDetails ProjectDetails) error {
	a.Config.Projects[projectDetails.Id] = projectDetails

	if err := a.writeConfigToFile(); err != nil {
		return fmt.Errorf("cannot write project details to app config: %v", err)
	}
	return nil
}

func (a *ApplicationManager) writeProjectToAppConfig(project Story) error {
	a.Config.Projects[project.Id] = ProjectDetails{
		Name:     project.Name,
		Id:       project.Id,
		Location: project.Location,
	}

	if err := a.writeConfigToFile(); err != nil {
		return fmt.Errorf("cannot write project details to app config: %v", err)
	}
	return nil
}
