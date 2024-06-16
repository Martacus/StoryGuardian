package project

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"io"
	"os"
)

const (
	ApplicationName = "StoryGuardian"
)

type ApplicationConfig struct {
	Projects map[string]ProjectDetails `json:"projects"`
}

func getConfigLocation() (string, error) {
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

func writeStructToFilePath(toWrite interface{}, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	jsonData, err := json.Marshal(toWrite)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	if _, err = file.Write(jsonData); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

func writeFileToStruct(filename string, data interface{}) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// Unmarshal the JSON data into the provided struct
	if err := json.Unmarshal(fileContents, data); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	return nil
}
