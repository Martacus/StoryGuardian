package project

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"storyguardian/src/constants"
	"storyguardian/src/fileio"
)

type ProjectDetails struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Config struct {
	ProjectDetails
}

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

	err = fileio.WriteStructToFile(a.Config, file)
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

func (a *ApplicationManager) CreateProject() (string, error) {
	projectDirectory, err := promptForProjectDirectory()
	if err != nil {
		return "", fmt.Errorf("could not select project directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(projectDirectory, constants.EntityFolderName), 0755); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(projectDirectory, constants.RelationsFolderName), 0755); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(projectDirectory, constants.ImagesFolderName), 0755); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	return projectDirectory, nil
}

func (a *ApplicationManager) OpenProject(projectId string) error {
	for key, val := range a.Config.Projects {
		if key == projectId {
			a.Config.OpenProject = val
		}
	}

	if err := a.writeConfigToFile(); err != nil {
		return fmt.Errorf("cannot write openProject to app config: %v", err)
	}

	return nil
}

func promptForProjectDirectory() (string, error) {
	projectDirectory, err := application.OpenFileDialog().CanChooseDirectories(true).PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("failed to prompt for project directory: %v", err)
	}
	return projectDirectory, nil
}

func (a *ApplicationManager) OpenProjectFolder(folder string) error {
	var cmd *exec.Cmd
	path := filepath.Join(a.Config.OpenProject.Location, folder)

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
