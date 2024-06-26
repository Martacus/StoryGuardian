package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/src/constants"
	"storyguardian/src/fileio"
	"time"
)

type Entity struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	StoryId     string                 `json:"storyId"`
	Description string                 `json:"description"`
	CreatedAt   string                 `json:"createdAt"`
	Tags        []string               `json:"tags"`
	Relations   []string               `json:"relations"`
	Modules     map[string]StoryModule `json:"modules"`
}

type EntityManager struct {
	StoryManager *StoryManager
	Entities     map[string]*Entity
}

type Category struct {
}

func NewEntityManager(projectManager *StoryManager) *EntityManager {
	entityManager := EntityManager{
		StoryManager: projectManager,
		Entities:     make(map[string]*Entity),
	}

	return &entityManager
}

func (e *EntityManager) LoadEntities() ([]Entity, error) {
	entitiesPath := filepath.Join(e.StoryManager.Story.Location, "entities")
	files, err := os.ReadDir(entitiesPath)
	if err != nil {
		return nil, fmt.Errorf("could not read entities directory: %v", err)
	}

	var entityList []Entity
	for _, file := range files {
		var newEntity Entity
		if err := fileio.WriteFilePathToStruct(filepath.Join(entitiesPath, file.Name()), &newEntity); err != nil {
			return make([]Entity, 0), fmt.Errorf("could not load entity: %v", err)
		}
		entityList = append(entityList, newEntity)
	}

	e.Entities = make(map[string]*Entity)
	for _, entity := range entityList {
		e.Entities[entity.Id] = &entity
	}

	return entityList, nil
}

func (e *EntityManager) GetEntity(entityId string) (*Entity, error) {
	entity, ok := e.Entities[entityId]
	if !ok {
		return nil, fmt.Errorf("could not find the entity with id: %v", entityId)
	}
	return entity, nil
}

func (e *EntityManager) RefreshEntity(entityId string) (*Entity, error) {
	entitiesPath := filepath.Join(e.StoryManager.Story.Location, "entities")
	var newEntity Entity
	if err := fileio.WriteFilePathToStruct(filepath.Join(entitiesPath, entityId+".json"), &newEntity); err != nil {
		return nil, fmt.Errorf("could not load entity: %v", err)
	}

	e.Entities[newEntity.Id] = &newEntity

	return &newEntity, nil
}

func (e *EntityManager) CreateEntity(entity Entity) (Entity, error) {
	layout := "2006-01-02 15:04:05"
	entity.CreatedAt = time.Now().Format(layout)

	filePath := e.getEntityFilePath(entity.Id)
	if err := fileio.WriteStructToFilePath(entity, filePath); err != nil {
		return entity, fmt.Errorf("could save new entity to file: %v", err)
	}

	e.Entities[entity.Id] = &entity

	return entity, nil
}

func (e *EntityManager) SaveEntityById(entityId string) error {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("failed to get entity: %v", err)
	}

	if err = fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entityId)); err != nil {
		return fmt.Errorf("could not save the entity: %v", err)
	}

	return nil
}

func (e *EntityManager) SaveEntity(entity Entity) error {
	if err := fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entity.Id)); err != nil {
		return fmt.Errorf("could not save the entity: %v", err)
	}

	return nil
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

func (e *EntityManager) SetEntityName(entityId string, name string) (string, error) {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return "", fmt.Errorf("failed to get entity: %v", err)
	}

	entity.Name = name
	if err = fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entityId)); err != nil {
		return "", fmt.Errorf("could not save the entity description change: %v", err)
	}

	return name, nil
}

func (e *EntityManager) getEntityFilePath(entityID string) string {
	return filepath.Join(e.StoryManager.Story.Location, constants.EntityFolderName, entityID+".json")
}
