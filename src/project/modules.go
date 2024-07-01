package project

import "fmt"

const (
	ImageStoryModuleID      = "images"
	TagListModuleID         = "tagList"
	DescriptionModuleID     = "description"
	EntityListStoryModuleID = "entityList"
	RelationsEntityModuleID = "relations"
	RelationInfoModuleID    = "relationInfo"
)

var availableStoryModules = []string{DescriptionModuleID, EntityListStoryModuleID, TagListModuleID, ImageStoryModuleID}
var availableEntityModules = []string{DescriptionModuleID, TagListModuleID, RelationsEntityModuleID}
var availableRelationModules = []string{DescriptionModuleID, RelationInfoModuleID}

func (s *StoryManager) GetStoryModules(unusedModulesOnly bool) []string {
	if !unusedModulesOnly {
		return availableStoryModules
	}

	var unusedStories []string
	for _, module := range availableStoryModules {
		if _, ok := s.Story.Modules[module]; !ok {
			unusedStories = append(unusedStories, module)
		}
	}
	return unusedStories
}

func (s *StoryManager) AddStoryModule(module string) error {
	switch module {
	case TagListModuleID:
		if err := addStoryTagsModule(s); err != nil {
			return err
		}
		break
	case ImageStoryModuleID:
		if err := addStoryImagesModule(s); err != nil {
			return err
		}
		break
	}

	return nil
}

func (s *StoryManager) EditStoryModuleConfig(module string, config string, value string) error {
	s.Story.Modules[module].Configuration[config] = value
	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save module configuration edit: %v", err)
	}
	return nil
}

func (e *EntityManager) GetEntityModules(entityId string, unusedModulesOnly bool) []string {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return make([]string, 0)
	}

	if !unusedModulesOnly {
		return availableEntityModules
	}

	var unusedEntityModules []string
	for _, module := range availableEntityModules {
		if _, ok := entity.Modules[module]; !ok {
			unusedEntityModules = append(unusedEntityModules, module)
		}
	}
	return unusedEntityModules
}

func (e *EntityManager) AddEntityModule(entityId string, module string) error {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("could not add module to entity: %v", err)
	}

	switch module {
	case TagListModuleID:
		if err := addEntityTagsModule(entity, e); err != nil {
			return err
		}
		break
	}

	return nil
}

func (e *EntityManager) EditEntityModuleConfig(entityId, module string, config string, value string) error {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("could not edit entity module config: %v", err)
	}

	entity.Modules[module].Configuration[config] = value

	if err := e.SaveEntity(*entity); err != nil {
		return fmt.Errorf("could not edit entity module config: %v", err)
	}
	return nil
}

func (r *RelationManager) GetRelationModules(relationid string, unusedModulesOnly bool) []string {
	relation, err := r.GetRelation(relationid)
	if err != nil {
		return make([]string, 0)
	}

	if !unusedModulesOnly {
		return availableRelationModules
	}

	var unusedModules []string
	for _, module := range availableRelationModules {
		if _, ok := relation.Modules[module]; !ok {
			unusedModules = append(unusedModules, module)
		}
	}
	return unusedModules
}

func (r *RelationManager) AddRelationModule(relationId string, module string) error {
	relation, err := r.GetRelation(relationId)
	if err != nil {
		return fmt.Errorf("could not add module to relation: %v", err)
	}
	_ = relation
	switch module {

	}

	return nil
}

func (r *RelationManager) EditRelationModuleConfig(relationId, module string, config string, value string) error {
	relation, err := r.GetRelation(relationId)
	if err != nil {
		return fmt.Errorf("could not edit relation module config: %v", err)
	}

	relation.Modules[module].Configuration[config] = value

	if err := r.SaveRelation(*relation); err != nil {
		return fmt.Errorf("could not edit relation module config: %v", err)
	}
	return nil
}

func addStoryImagesModule(manager *StoryManager) error {
	newImageModule := StoryModule{
		Name: ImageStoryModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
		},
	}
	manager.Story.Modules[ImageStoryModuleID] = newImageModule
	if err := manager.SaveStory(); err != nil {
		return err
	}

	return nil
}

func addStoryTagsModule(manager *StoryManager) error {
	newTagListModule := StoryModule{
		Name: TagListModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
			"itemView":   "list",
		},
	}
	manager.Story.Modules[TagListModuleID] = newTagListModule
	if err := manager.SaveStory(); err != nil {
		return err
	}

	return nil
}

func addEntityTagsModule(entity *Entity, manager *EntityManager) error {
	newTagListModule := StoryModule{
		Name: TagListModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
			"itemView":   "list",
		},
	}

	entity.Modules[TagListModuleID] = newTagListModule
	if err := manager.SaveEntity(*entity); err != nil {
		return fmt.Errorf("unable to add tag module to entity: %v", err)
	}

	return nil
}
