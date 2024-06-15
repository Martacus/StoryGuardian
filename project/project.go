package project

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"os"
	"path/filepath"
)

type ProjectDetails struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Config struct {
	ProjectDetails
}

func (a *ApplicationManager) CreateProject() (string, error) {
	projectDirectory, err := promptForProjectDirectory()
	if err != nil {
		return "", fmt.Errorf("could not select project directory: %v", err)
	}

	projectId := uuid.New().String()
	err = createProjectFile(projectDirectory, projectId)
	if err != nil {
		return "", fmt.Errorf("could not create the project config file: %v", err)
	}

	err = writeProjectToAppConfig(projectId, projectDirectory, a)
	if err != nil {
		return "", fmt.Errorf("could not write project to application config file: %v", err)
	}

	return projectId, nil
}

func (a *ApplicationManager) OpenProject(projectDirectory string) (string, error) {
	//write nw project to config

	return projectDirectory, nil
}

func promptForProjectDirectory() (string, error) {
	projectDirectory, err := application.OpenFileDialog().CanChooseDirectories(true).PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("failed to prompt for project directory: %v", err)
	}
	return projectDirectory, nil
}

func createProjectFile(projectDirectory string, projectId string) error {
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
	config := Config{
		ProjectDetails{
			Id:       projectId,
			Name:     projectName,
			Location: projectDirectory,
		},
	}

	return writeStructToFile(config, file)
}

func writeProjectToAppConfig(projectId string, projectDirectory string, a *ApplicationManager) error {
	projectName := filepath.Base(projectDirectory)
	a.Config.Projects[projectId] = ProjectDetails{
		Id:       projectId,
		Name:     projectName,
		Location: projectDirectory,
	}

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
