package project

import (
	"fmt"
)

const (
	ImageStoryModuleID      = "images"
	TagListModuleID         = "tagList"
	DescriptionModuleID     = "description"
	EntityListStoryModuleID = "entityList"
	RelationsEntityModuleID = "relations"
	RelationInfoModuleID    = "relationInfo"
)

var (
	availableStoryModules    = []string{DescriptionModuleID, EntityListStoryModuleID, TagListModuleID, ImageStoryModuleID}
	availableEntityModules   = []string{DescriptionModuleID, TagListModuleID, RelationsEntityModuleID}
	availableRelationModules = []string{DescriptionModuleID, RelationInfoModuleID}
)

func (s *StoryManager) GetStoryModules(unusedModulesOnly bool) []string {
	return getUnusedModules(unusedModulesOnly, availableStoryModules, s.Story.Modules)
}

func (s *StoryManager) AddStoryModule(module string) error {
	switch module {
	case TagListModuleID:
		return addStoryTagsModule(s)
	case ImageStoryModuleID:
		return addStoryImagesModule(s)
	default:
		return fmt.Errorf("unknown story module: %s", module)
	}
}

func (s *StoryManager) EditStoryModuleConfig(module, config, value string) error {
	if _, exists := s.Story.Modules[module]; !exists {
		return fmt.Errorf("story module %s does not exist", module)
	}
	s.Story.Modules[module].Configuration[config] = value
	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save module configuration edit: %v", err)
	}
	return nil
}

func (e *EntityManager) GetEntityModules(entityID string, unusedModulesOnly bool) []string {
	entity, err := e.GetEntity(entityID)
	if err != nil {
		return []string{}
	}
	return getUnusedModules(unusedModulesOnly, availableEntityModules, entity.Modules)
}

func (e *EntityManager) AddEntityModule(entityID, module string) error {
	entity, err := e.GetEntity(entityID)
	if err != nil {
		return fmt.Errorf("could not add module to entity: %v", err)
	}

	switch module {
	case TagListModuleID:
		return addEntityTagsModule(entity, e)
	default:
		return fmt.Errorf("unknown entity module: %s", module)
	}
}

func (e *EntityManager) EditEntityModuleConfig(entityID, module, config, value string) error {
	entity, err := e.GetEntity(entityID)
	if err != nil {
		return fmt.Errorf("could not edit entity module config: %v", err)
	}
	if _, exists := entity.Modules[module]; !exists {
		return fmt.Errorf("entity module %s does not exist", module)
	}
	entity.Modules[module].Configuration[config] = value
	if err := e.SaveEntity(*entity); err != nil {
		return fmt.Errorf("could not save entity module configuration: %v", err)
	}
	return nil
}

func (r *RelationManager) GetRelationModules(relationID string, unusedModulesOnly bool) []string {
	relation, err := r.GetRelation(relationID)
	if err != nil {
		return []string{}
	}
	return getUnusedModules(unusedModulesOnly, availableRelationModules, relation.Modules)
}

func (r *RelationManager) AddRelationModule(relationID, module string) error {
	relation, err := r.GetRelation(relationID)
	if err != nil {
		return fmt.Errorf("could not add module to relation: %v", err)
	}
	_ = relation
	switch module {

	}
	return fmt.Errorf("unknown relation module: %s", module)
}

func (r *RelationManager) EditRelationModuleConfig(relationID, module, config, value string) error {
	relation, err := r.GetRelation(relationID)
	if err != nil {
		return fmt.Errorf("could not edit relation module config: %v", err)
	}
	if _, exists := relation.Modules[module]; !exists {
		return fmt.Errorf("relation module %s does not exist", module)
	}
	relation.Modules[module].Configuration[config] = value
	if err := r.SaveRelation(*relation); err != nil {
		return fmt.Errorf("could not save relation module configuration: %v", err)
	}
	return nil
}

func getUnusedModules(unusedModulesOnly bool, availableModules []string, currentModules map[string]StoryModule) []string {
	if !unusedModulesOnly {
		return availableModules
	}

	var unusedModules []string
	for _, module := range availableModules {
		if _, ok := currentModules[module]; !ok {
			unusedModules = append(unusedModules, module)
		}
	}
	return unusedModules
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
		return fmt.Errorf("unable to add image module to story: %v", err)
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
		return fmt.Errorf("unable to add tag module to story: %v", err)
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
