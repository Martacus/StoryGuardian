package project

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
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
	relationManager := RelationManager{
		EntityManager: entityManager,
		Relations:     make(map[string]*Relation),
	}

	return &relationManager
}

func (r *Relation) getOther(entityId string) string {
	if r.EntityOne == entityId {
		return r.EntityTwo
	}
	return r.EntityOne
}

func (r *RelationManager) LoadRelations() error {
	var relationList []Relation
	relationFolderPath := filepath.Join(r.EntityManager.StoryManager.Story.Location, "relations")

	files, err := os.ReadDir(relationFolderPath)
	if err != nil {
		return fmt.Errorf("could not load relations: %v", err)
	}

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
	_ = paginationStart
	_ = amount

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

		var relatedEntity *Entity
		relatedEntity, err = r.EntityManager.GetEntity(relation.getOther(entityId))
		if err != nil {
			return nil, fmt.Errorf("failed to load relation info: %v", err)
		}

		//Might need to remove depointer on relation
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

	//Write the relation to the system
	filePath := r.getRelationPath(relation.Id)
	if err := fileio.WriteStructToFilePath(relation, filePath); err != nil {
		return "", fmt.Errorf("failed to create relation: %v", err)
	}

	//Write the relation id to the entity
	entity.Relations = append(entity.Relations, relation.Id)
	if err = fileio.WriteStructToFilePath(entity, r.EntityManager.getEntityFilePath(entityId)); err != nil {
		return "", fmt.Errorf("failed to create relation: %v", err)
	}

	r.Relations[relation.Id] = &relation

	return relation.Id, nil
}

func (r *RelationManager) getRelationPath(relationId string) string {
	return filepath.Join(r.EntityManager.StoryManager.Story.Location, constants.RelationsFolderName, relationId+".json")
}

func (r *RelationManager) GetRelation(relationId string) (*Relation, error) {
	relation, ok := r.Relations[relationId]
	if !ok {
		return nil, fmt.Errorf("could retrieve relation: %v", relationId)
	}

	return relation, nil
}

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
