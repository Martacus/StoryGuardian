package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/src/utility"

	"github.com/google/uuid"
	"storyguardian/src/constants"
	"storyguardian/src/fileio"
)

type Relation struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	EntityOne   string                 `json:"entityOne"`
	EntityTwo   string                 `json:"entityTwo"`
	Description string                 `json:"description"`
	Modules     map[string]StoryModule `json:"modules"`
}

type RelationInfo struct {
	Relation
	ToName string `json:"toName"`
}

type RelationManager struct {
	EntityManager *EntityManager
	Relations     map[string]*Relation
}

func NewRelationManager(entityManager *EntityManager) *RelationManager {
	return &RelationManager{
		EntityManager: entityManager,
		Relations:     make(map[string]*Relation),
	}
}

func (r *Relation) getOther(entityId string) string {
	if r.EntityOne == entityId {
		return r.EntityTwo
	}
	return r.EntityOne
}

func (r *RelationManager) LoadRelations() error {
	relationFolderPath := filepath.Join(r.EntityManager.StoryManager.Story.Location, "relations")
	files, err := os.ReadDir(relationFolderPath)
	if err != nil {
		return fmt.Errorf("could not load relations: %v", err)
	}

	var relationList []Relation
	for _, file := range files {
		var relation Relation
		if err := fileio.WriteFilePathToStruct(filepath.Join(relationFolderPath, file.Name()), &relation); err != nil {
			return fmt.Errorf("could not load relations: %v", err)
		}
		relationList = append(relationList, relation)
	}

	r.Relations = make(map[string]*Relation)
	for _, relation := range relationList {
		r.Relations[relation.Id] = &relation
	}

	return nil
}

func (r *RelationManager) LoadRelationInfo(entityId string, paginationStart int, amount int) ([]RelationInfo, error) {
	entity, err := r.EntityManager.GetEntity(entityId)
	if err != nil {
		return nil, fmt.Errorf("failed to load relation info: %v", err)
	}

	var relationInfo []RelationInfo
	for _, relationId := range entity.Relations {
		relation, ok := r.Relations[relationId]
		if !ok {
			return nil, fmt.Errorf("failed to load relation info: %v", err)
		}

		if relation.EntityTwo == "" {
			relationInfo = append(relationInfo, RelationInfo{
				ToName:   "-Empty-",
				Relation: *relation,
			})
			return relationInfo, nil
		}

		relatedEntity, err := r.EntityManager.GetEntity(relation.getOther(entityId))
		if err != nil {
			return nil, fmt.Errorf("failed to load relation info: %v", err)
		}

		relationInfo = append(relationInfo, RelationInfo{
			ToName:   relatedEntity.Name,
			Relation: *relation,
		})
	}
	return relationInfo, nil
}

func (r *RelationManager) CreateRelation(entityId string) (string, error) {
	entity, err := r.EntityManager.GetEntity(entityId)
	if err != nil {
		return "", fmt.Errorf("failed to create relation: %v", err)
	}

	infoModule := StoryModule{
		Name: RelationInfoModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
		},
	}

	descModule := StoryModule{
		Name: DescriptionModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
		},
	}

	relation := Relation{
		Id:          uuid.New().String(),
		EntityOne:   entityId,
		EntityTwo:   "",
		Name:        "Placeholder",
		Description: "Placeholder",
		Modules: map[string]StoryModule{
			RelationInfoModuleID: infoModule,
			DescriptionModuleID:  descModule,
		},
	}

	filePath := r.getRelationPath(relation.Id)
	if err := fileio.WriteStructToFilePath(relation, filePath); err != nil {
		return "", fmt.Errorf("failed to create relation: %v", err)
	}

	entity.Relations = append(entity.Relations, relation.Id)
	if err = fileio.WriteStructToFilePath(entity, r.EntityManager.getEntityFilePath(entityId)); err != nil {
		return "", fmt.Errorf("failed to create relation: %v", err)
	}

	r.Relations[relation.Id] = &relation

	return relation.Id, nil
}

// getRelationPath returns the path to the relation file
func (r *RelationManager) getRelationPath(relationId string) string {
	return filepath.Join(r.EntityManager.StoryManager.Story.Location, constants.RelationsFolderName, relationId+".json")
}

// GetRelation returns the relation with the given id
func (r *RelationManager) GetRelation(relationId string) (*Relation, error) {
	relation, ok := r.Relations[relationId]
	if !ok {
		return nil, fmt.Errorf("could not retrieve relation: %v", relationId)
	}

	return relation, nil
}

// SetRelationDescription sets the description of the relation with the given id
func (r *RelationManager) SetRelationDescription(relationId string, description string) (string, error) {
	relation, ok := r.Relations[relationId]
	if !ok {
		return "", fmt.Errorf("could not set description: relation not found")
	}

	relation.Description = description
	if err := r.SaveRelation(*relation); err != nil {
		return "", fmt.Errorf("could not set description: %v", err)
	}

	return description, nil
}

// SetRelationName sets the name of the relation with the given id
func (r *RelationManager) SetRelationName(relationId string, name string) (string, error) {
	relation, ok := r.Relations[relationId]
	if !ok {
		return "", fmt.Errorf("could not set relation name: relation not found")
	}

	relation.Name = name
	if err := r.SaveRelation(*relation); err != nil {
		return "", fmt.Errorf("could not set relation name: %v", err)
	}

	return name, nil
}

func (r *RelationManager) SaveRelation(relation Relation) error {
	if err := fileio.WriteStructToFilePath(relation, r.getRelationPath(relation.Id)); err != nil {
		return fmt.Errorf("could not save the relation to file system: %v", err)
	}

	return nil
}

func (r *RelationManager) SetRelationEntities(entityOneId string, entityTwoId string, relationId string) error {
	relation, ok := r.Relations[relationId]
	if !ok {
		return fmt.Errorf("could not set relation entities: relation not found")
	}

	if relation.EntityOne != entityOneId {
		if err := r.removeRelationFromEntity(relation.EntityOne, relationId); err != nil {
			return err
		}
		relation.EntityOne = entityOneId
		if err := r.addRelationToEntity(entityOneId, relationId); err != nil {
			return err
		}
	}

	if relation.EntityTwo != entityTwoId {
		if err := r.removeRelationFromEntity(relation.EntityTwo, relationId); err != nil {
			return err
		}
		relation.EntityTwo = entityTwoId
		if err := r.addRelationToEntity(entityTwoId, relationId); err != nil {
			return err
		}
	}

	if err := r.SaveRelation(*relation); err != nil {
		return fmt.Errorf("could not set relation entities: %v", err)
	}

	return nil
}

func (r *RelationManager) DeleteRelation(relationId string) error {
	relation, ok := r.Relations[relationId]
	if !ok {
		return fmt.Errorf("could not delete relation: relation not found")
	}

	if err := r.removeRelationFromEntity(relation.EntityOne, relationId); err != nil {
		return err
	}

	if err := r.removeRelationFromEntity(relation.EntityTwo, relationId); err != nil {
		return err
	}

	filePath := r.getRelationPath(relationId)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete relation: %v", err)
	}

	delete(r.Relations, relationId)
	return nil
}

func (r *RelationManager) removeRelationFromEntity(entityId, relationId string) error {
	if entityId == "" {
		return nil
	}

	entity, err := r.EntityManager.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("could not find entity in relation: %v", err)
	}

	entity.Relations = utility.RemoveStringFromSlice(entity.Relations, relationId)
	if err := r.EntityManager.SaveEntity(*entity); err != nil {
		return err
	}

	return nil
}

func (r *RelationManager) addRelationToEntity(entityId, relationId string) error {
	entity, err := r.EntityManager.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("could not find entity in relation: %v", err)
	}

	if entity.Relations == nil {
		entity.Relations = []string{}
	}
	entity.Relations = append(entity.Relations, relationId)
	if err := r.EntityManager.SaveEntity(*entity); err != nil {
		return err
	}

	return nil
}
