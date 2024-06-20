package project

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
	"storyguardian/internal/constants"
	"storyguardian/internal/fileio"
)

const (
	ApplicationName = "StoryGuardian"
)

type ApplicationConfig struct {
	Projects    map[string]ProjectDetails `json:"projects"`
	OpenProject ProjectDetails            `json:"openProject"`
}

// getConfigLocation retrieves the path to the configuration file for the application.
// It first attempts to locate an existing configuration file using xdg.SearchConfigFile.
// If the configuration file does not exist, it creates a new one using xdg.ConfigFile
// and then calls createConfigFile to initialize it.
//
// If any error occurs during the file search or creation process, an error is returned.
//
// Returns:
//   - string: The path to the configuration file.
//   - error: An error if the configuration file path cannot be retrieved or created.
func getConfigLocation() (string, error) {
	configFilePath, err := xdg.SearchConfigFile(filepath.Join(ApplicationName, constants.ConfigFileName))

	//If te config doesn't exist, create it
	if err != nil {
		configFilePath, err = xdg.ConfigFile(filepath.Join(ApplicationName, constants.ConfigFileName))
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
		Projects: make(map[string]ProjectDetails),
	}

	return fileio.WriteStructToFile(config, file)
}
