package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"os"
	"path/filepath"
)

const (
	ApplicationName = "StoryGuardian"
)

type ConfigManager interface {
	GetConfig() ApplicationConfig
}

type ProjectManager interface {
	CreateProject()
}

type ApplicationManager struct {
	ctx    *context.Context
	Config ApplicationConfig
}

type Project struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type ApplicationConfig struct {
	Projects map[string]Project `json:"projects"`
}

type Config struct {
	Name string `json:"name"`
}

func NewProjectManager(ctx *context.Context) *ApplicationManager {
	configLocation, err := configOnStartup()
	if err != nil {
		log.Fatalf("Could not find/create the config file: %v", err)
	}

	configBytes, err := os.ReadFile(configLocation)
	if err != nil {
		log.Fatalf("Could not read config file while it has been found: %v", err)
	}

	var config ApplicationConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config file into object: %v", err)
	}

	return &ApplicationManager{
		ctx:    ctx,
		Config: config,
	}
}

func (a *ApplicationManager) GetConfig() ApplicationConfig {
	return a.Config
}

func configOnStartup() (string, error) {
	configFilePath, err := xdg.SearchConfigFile(ApplicationName + "/config.json")
	if err != nil {
		err := createConfigFile()
		if err != nil {
			return "", fmt.Errorf("could not create the necessary config file: %v", err)
		}
	}

	return configFilePath, nil
}

func (a *ApplicationManager) CreateProject() (string, error) {
	projectDirectory, err := application.OpenFileDialog().CanChooseDirectories(true).PromptForSingleSelection()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(projectDirectory, "project.json")
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	projectName := filepath.Base(filePath)
	config := Config{Name: projectName}

	err = writeStructToFile(config, file)
	if err != nil {
		return "", err
	}

	return projectDirectory, nil
}

func createConfigFile() error {
	configFilePath, err := xdg.ConfigFile(ApplicationName + "/config.json")
	if err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	config := ApplicationConfig{
		Projects: make(map[string]Project),
	}

	err = writeStructToFile(config, file)
	if err != nil {
		return err
	}

	return nil
}

func writeStructToFile(toWrite interface{}, file *os.File) error {
	jsonData, err := json.Marshal(toWrite)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}
