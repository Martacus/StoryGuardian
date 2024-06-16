package project

import (
	"fmt"
	"path/filepath"
	"storyguardian/constants"
)

type Project struct {
	ProjectDetails
	Description string `json:"description"`
}

type ProjectManager struct {
	ApplicationManager *ApplicationManager
}

func NewProjectManager(appManager *ApplicationManager) *ProjectManager {
	return &ProjectManager{
		ApplicationManager: appManager,
	}
}

func (p *ProjectManager) GetProject(projectId string) (*Project, error) {
	for key, projectDetails := range p.ApplicationManager.Config.Projects {
		if key == projectId {
			project := Project{}
			err := writeFileToStruct(filepath.Join(projectDetails.Location, constants.ProjectConfigName), &project)
			if err != nil {
				return nil, fmt.Errorf("could not write project json to struct: %v", err)
			}

			return &project, nil
		}
	}
	return nil, fmt.Errorf("could not find the project with id %v", projectId)
}

func (p *ProjectManager) SetProjectTitle(projectId string, name string) error {
	project, err := p.GetProject(projectId)
	if err != nil {
		return fmt.Errorf("could not find the project: %v", err)
	}

	project.Name = name
	if err = writeStructToFilePath(project, filepath.Join(project.Location, constants.ProjectConfigName)); err != nil {
		return fmt.Errorf("could not save the title change: %v", err)
	}

	if err = p.ApplicationManager.writeProjectToAppConfig(*project); err != nil {
		return fmt.Errorf("could not save the title change to application config: %v", err)
	}

	return nil
}

func (p *ProjectManager) SetProjectDescription(projectId string, description string) (string, error) {
	project, err := p.GetProject(projectId)
	if err != nil {
		return "", fmt.Errorf("could not find the project with id: %v | %v", projectId, err)
	}

	project.Description = description
	if err = writeStructToFilePath(project, filepath.Join(project.Location, constants.ProjectConfigName)); err != nil {
		return "", fmt.Errorf("could not save the description change: %v", err)
	}

	return description, nil
}
