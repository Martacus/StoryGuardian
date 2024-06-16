package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/constants"
	"time"
)

type Entity struct {
	Id          string
	Name        string
	StoryId     string
	Description string
	CreatedAt   string
	Tags        []string
}
type EntityManager struct {
	ProjectManager *StoryManager
}

func NewEntityManager(projectManager *StoryManager) *EntityManager {
	return &EntityManager{
		ProjectManager: projectManager,
	}
}

func (e *EntityManager) GetEntities(projectId string) ([]Entity, error) {
	project, err := e.ProjectManager.GetStory(projectId)
	if err != nil {
		return make([]Entity, 0), err
	}

	entityList := make([]Entity, 0)

	entitiesPath := filepath.Join(project.Location, "entities")
	files, err := os.ReadDir(entitiesPath)
	for _, file := range files {
		newEntity := Entity{}
		if err := writeFilePathToStruct(filepath.Join(entitiesPath, file.Name()), &newEntity); err != nil {
			return make([]Entity, 0), fmt.Errorf("could not load entity: %v", err)
		}
		entityList = append(entityList, newEntity)
	}

	return entityList, nil
}

func (e *EntityManager) CreateEntity(entity Entity) (Entity, error) {
	project := e.ProjectManager.ApplicationManager.Config.OpenProject
	entity.CreatedAt = time.Now().String()

	filePath := filepath.Join(project.Location, constants.EntityFolderName, entity.Id+".json")
	if err := writeStructToFilePath(entity, filePath); err != nil {
		return entity, fmt.Errorf("could save new entity to file: %v", err)
	}

	return entity, nil
}
