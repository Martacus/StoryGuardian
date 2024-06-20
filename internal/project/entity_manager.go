package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/internal/constants"
	"storyguardian/internal/fileio"
	"time"
)

type Entity struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	StoryId     string     `json:"storyId"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"createdAt"`
	Tags        []string   `json:"tags"`
	Relations   []Relation `json:"relations"`
}

type Relation struct {
	Name        string `json:"name"`
	To          string `json:"to"`
	Description string `json:"description"`
}

type RelationInfo struct {
	Relation
	ToName string `json:"toName"`
}

type EntityManager struct {
	ProjectManager *StoryManager
	Entities       map[string]*Entity
}

func NewEntityManager(projectManager *StoryManager) *EntityManager {
	entityManager := EntityManager{
		ProjectManager: projectManager,
		Entities:       make(map[string]*Entity),
	}

	return &entityManager
}

func (e *EntityManager) LoadEntities(projectId string) ([]Entity, error) {
	project, err := e.ProjectManager.GetStory(projectId)
	if err != nil {
		return make([]Entity, 0), err
	}

	entityList := make([]Entity, 0)

	entitiesPath := filepath.Join(project.Location, "entities")
	files, err := os.ReadDir(entitiesPath)
	for _, file := range files {
		newEntity := Entity{}
		if err := fileio.WriteFilePathToStruct(filepath.Join(entitiesPath, file.Name()), &newEntity); err != nil {
			return make([]Entity, 0), fmt.Errorf("could not load entity: %v", err)
		}
		entityList = append(entityList, newEntity)
	}

	entityMap := make(map[string]*Entity)
	for _, entity := range entityList {
		entityMap[entity.Id] = &entity
	}
	e.Entities = entityMap

	return entityList, nil
}

func (e *EntityManager) LoadRelationInfo(entityId string) ([]RelationInfo, error) {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity: %v", err)
	}

	relationInfo := make([]RelationInfo, 0)
	for _, relation := range entity.Relations {
		relationEntity, err := e.GetEntity(relation.To)
		if err != nil {
			return nil, fmt.Errorf("failed to get entity: %v", err)
		}

		relationInfo = append(relationInfo, RelationInfo{
			ToName:   relationEntity.Name,
			Relation: relation,
		})
	}
	return relationInfo, nil

}

func (e *EntityManager) GetEntity(entityId string) (*Entity, error) {
	entity, ok := e.Entities[entityId]
	if !ok {
		return nil, fmt.Errorf("could not find the entity with id: %v", entityId)
	}
	return entity, nil
}

func (e *EntityManager) CreateEntity(entity Entity) (Entity, error) {
	project := e.ProjectManager.ApplicationManager.Config.OpenProject
	entity.CreatedAt = time.Now().String()

	filePath := filepath.Join(project.Location, constants.EntityFolderName, entity.Id+".json")
	if err := fileio.WriteStructToFilePath(entity, filePath); err != nil {
		return entity, fmt.Errorf("could save new entity to file: %v", err)
	}

	e.Entities[entity.Id] = &entity

	return entity, nil
}

func (e *EntityManager) SetEntityDescription(entityId string, description string) (string, error) {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return "", fmt.Errorf("failed to get entity: %v", err)
	}

	entity.Description = description
	if err = fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entityId)); err != nil {
		return "", fmt.Errorf("could not save the entity description change: %v", err)
	}

	return description, nil
}

func (e *EntityManager) getEntityFilePath(entityID string) string {
	return filepath.Join(e.ProjectManager.ApplicationManager.Config.OpenProject.Location, constants.EntityFolderName, entityID+".json")
}

func (e *EntityManager) CreateRelation(entityId string, relation Relation) error {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("failed to get entity: %v", err)
	}

	entity.Relations = append(entity.Relations, relation)

	if err = fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entityId)); err != nil {
		return fmt.Errorf("could not save new relation: %v", err)
	}

	return nil
}
