package project

import (
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
	CreateProject() (string, error)
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

type ApplicationManager struct {
	Config ApplicationConfig
}

func NewProjectManager() *ApplicationManager {
	configLocation, err := configOnStartup()
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

func (a *ApplicationManager) CreateProject() (string, error) {
	projectDirectory, err := promptForProjectDirectory()
	if err != nil {
		return "", fmt.Errorf("could not select project directory: %v", err)
	}

	err = createProjectFile(projectDirectory)
	if err != nil {
		return "", fmt.Errorf("could create the project config file: %v", err)
	}

	return projectDirectory, nil
}

func configOnStartup() (string, error) {
	configFilePath, err := xdg.SearchConfigFile(ApplicationName + "/config.json")

	//If te config doesn't exist, create it
	if err != nil {
		configFilePath, err = xdg.ConfigFile(ApplicationName + "/config.json")
		if err != nil {
			return "", fmt.Errorf("could not create the necessary config file path: %v", err)
		}
		err = createConfigFile(configFilePath)
		if err != nil {
			return "", fmt.Errorf("could not create the necessary config file: %v", err)
		}
	}

	return configFilePath, nil
}

func loadConfig(configLocation string) (ApplicationConfig, error) {
	configBytes, err := os.ReadFile(configLocation)
	if err != nil {
		return ApplicationConfig{}, fmt.Errorf("could not read config file: %v", err)
	}

	var config ApplicationConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return ApplicationConfig{}, fmt.Errorf("could not unmarshal config file: %v", err)
	}

	return config, nil
}

func promptForProjectDirectory() (string, error) {
	projectDirectory, err := application.OpenFileDialog().CanChooseDirectories(true).PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("failed to prompt for project directory: %v", err)
	}
	return projectDirectory, nil
}

func createProjectFile(projectDirectory string) error {
	filePath := filepath.Join(projectDirectory, "project.json")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	projectName := filepath.Base(projectDirectory)
	config := Config{Name: projectName}

	return writeStructToFile(config, file)
}

func createConfigFile(configFilePath string) error {
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	config := ApplicationConfig{
		Projects: make(map[string]Project),
	}

	return writeStructToFile(config, file)
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
