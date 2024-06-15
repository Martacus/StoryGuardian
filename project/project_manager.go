package project

import (
	"fmt"
	"log"
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
		log.Printf("%v : %v", key, projectId)
		if key == projectId {
			project := Project{}
			err := writeFileToStruct(projectDetails.Location+"/project.json", &project)
			if err != nil {
				return nil, fmt.Errorf("could not write project json to struct: %v", err)
			}

			return &project, nil
		}
	}
	return nil, fmt.Errorf("could not find the project with id %v", projectId)
}
