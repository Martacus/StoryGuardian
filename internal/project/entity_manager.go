package project

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"storyguardian/internal/constants"
	"storyguardian/internal/fileio"
	"time"
)

type Entity struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	StoryId     string   `json:"storyId"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"createdAt"`
	Tags        []string `json:"tags"`
	Relations   []string `json:"relations"`
}

type Relation struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	EntityOne   string `json:"entityOne"`
	EntityTwo   string `json:"entityTwo"`
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
	return filepath.Join(e.ProjectManager.ApplicationManager.Config.OpenProject.Location, constants.EntityFolderName, entityID+".json")
}

// Relations
func (r *Relation) getOther(entityId string) string {
	if r.EntityOne == entityId {
		return r.EntityTwo
	} else {
		return r.EntityOne
	}
}

func (e *EntityManager) LoadRelationInfo(entityId string, paginationStart int, amount int) ([]RelationInfo, error) {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity: %v", err)
	}

	relationInfo := make([]RelationInfo, 0)
	for _, relationId := range entity.Relations {
		relation := Relation{}
		if err := fileio.WriteFilePathToStruct(e.getRelationPath(relationId), &relation); err != nil {
			return nil, fmt.Errorf("failed to write relation to struct: %v", err)
		}

		relationEntity, err := e.GetEntity(relation.getOther(entityId))
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

func (e *EntityManager) CreateRelation(entityId string) (string, error) {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return "", fmt.Errorf("failed to get entity: %v", err)
	}

	relation := Relation{
		Id:          uuid.New().String(),
		EntityOne:   entityId,
		EntityTwo:   "e01e7ce9-36e1-4943-af20-6da01f77038a",
		Name:        "Placeholder",
		Description: "Placeholder",
	}

	//Write the relation to the system
	filePath := filepath.Join(e.ProjectManager.ApplicationManager.Config.OpenProject.Location, constants.RelationsFolderName, relation.Id+".json")
	if err := fileio.WriteStructToFilePath(relation, filePath); err != nil {
		return "", fmt.Errorf("could save new relation to file: %v", err)
	}

	//Write the relation id to the entity
	entity.Relations = append(entity.Relations, relation.Id)
	if err = fileio.WriteStructToFilePath(entity, e.getEntityFilePath(entityId)); err != nil {
		return "", fmt.Errorf("could not save entity with new relation: %v", err)
	}

	return relation.Id, nil
}

func (e *EntityManager) getRelationPath(relationId string) string {
	return filepath.Join(e.ProjectManager.ApplicationManager.Config.OpenProject.Location, constants.RelationsFolderName, relationId+".json")
}
