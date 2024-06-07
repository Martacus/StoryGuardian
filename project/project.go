package project

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"log"
	"os"
)

const (
	ProjectName = "StoryGuardian"
)

type ConfigManager interface {
	GetConfig() ProjectConfig
}

type ProjectManager struct {
	Config ProjectConfig
}

type Project struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type ProjectConfig struct {
	Projects map[string]Project `json:"projects"`
}

func NewProjectManager() *ProjectManager {
	configLocation, err := configOnStartup()
	if err != nil {
		log.Fatalf("Could not find/create the config file: %v", err)
	}

	configBytes, err := os.ReadFile(configLocation)
	if err != nil {
		log.Fatalf("Could not read config file while it has been found: %v", err)
	}

	var config ProjectConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config file into object: %v", err)
	}

	return &ProjectManager{
		Config: config,
	}
}

func (p *ProjectManager) GetConfig() ProjectConfig {
	return p.Config
}

func configOnStartup() (string, error) {
	configFilePath, err := xdg.SearchConfigFile(ProjectName + "/config.json")
	if err != nil {
		err := createConfigFile()
		if err != nil {
			return "", fmt.Errorf("could not create the necessary config file: %v", err)
		}
	}

	return configFilePath, nil
}

func createConfigFile() error {
	configFilePath, err := xdg.ConfigFile(ProjectName + "/config.json")
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

	config := ProjectConfig{
		Projects: make(map[string]Project),
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
