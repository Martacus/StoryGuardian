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

func getUnusedModules(unusedModulesOnly bool, availableModules []string, currentModules []StoryModule) []string {
	if !unusedModulesOnly {
		return availableModules
	}

	var unusedModules []string
	for _, module := range availableModules {
		found := false
		for _, currentModule := range currentModules {
			if currentModule.Configuration["name"] == module {
				found = true
				break
			}
		}
		if !found {
			unusedModules = append(unusedModules, module)
		}
	}
	return unusedModules
}

// === Story Modules === //

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
	recModule, exists := storyModuleExists(module, s.Story.Modules)

	if !exists {
		return fmt.Errorf("story module %s does not exist", module)
	}
	recModule.Configuration[config] = value

	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save module configuration edit: %v", err)
	}
	return nil
}

func addStoryImagesModule(manager *StoryManager) error {
	newImageModule := StoryModule{
		Name: ImageStoryModuleID,
		Configuration: map[string]any{
			"columnSize": "4",
		},
	}
	manager.Story.Modules = append(manager.Story.Modules, newImageModule)
	if err := manager.SaveStory(); err != nil {
		return fmt.Errorf("unable to add image module to story: %v", err)
	}
	return nil
}

func addStoryTagsModule(manager *StoryManager) error {
	newTagListModule := StoryModule{
		Name: TagListModuleID,
		Configuration: map[string]any{
			"columnSize": "4",
			"itemView":   "list",
		},
	}
	manager.Story.Modules = append(manager.Story.Modules, newTagListModule)
	if err := manager.SaveStory(); err != nil {
		return fmt.Errorf("unable to add tag module to story: %v", err)
	}
	return nil
}

func storyModuleExists(module string, modules []StoryModule) (*StoryModule, bool) {
	for _, m := range modules {
		if m.Configuration["name"] == module {
			return &m, true
		}
	}

	return nil, false
}

// === Entity Modules === //

func (e *EntityManager) GetEntityModules(entityID string, unusedModulesOnly bool) []string {
	_, err := e.GetEntity(entityID)
	if err != nil {
		return []string{}
	}
	return nil //getUnusedModules(unusedModulesOnly, availableEntityModules, entity.Modules)
}

func (e *EntityManager) AddEntityModule(entityID, module string) error {
	entity, err := e.GetEntity(entityID)
	if err != nil {
		return fmt.Errorf("could not add module to entity: %v", err)
	}

	defaultConfig := map[string]any{
		"columnSize": "4",
		"itemView":   "list",
		"open":       "true",
	}

	switch module {
	case TagListModuleID:
		return addModuleToEntity(entity, e, TagListModuleID, defaultConfig)
	case RelationsEntityModuleID:
		return addModuleToEntity(entity, e, RelationsEntityModuleID, defaultConfig)
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

func addModuleToEntity(entity *Entity, manager *EntityManager, moduleID string, config map[string]any) error {
	newModule := StoryModule{
		Name:          moduleID,
		Configuration: config,
	}

	if entity.Modules == nil {
		entity.Modules = make(map[string]StoryModule)
	}
	entity.Modules[moduleID] = newModule

	if err := manager.SaveEntity(*entity); err != nil {
		return fmt.Errorf("unable to add %s module to entity: %v", moduleID, err)
	}
	return nil
}

// === Relation Modules === //

func (r *RelationManager) GetRelationModules(relationID string, unusedModulesOnly bool) []string {
	_, err := r.GetRelation(relationID)
	if err != nil {
		return []string{}
	}
	return nil // getUnusedModules(unusedModulesOnly, availableRelationModules, relation.Modules)
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
	fmt.Println("0")
	relation, err := r.GetRelation(relationID)
	if err != nil {
		return fmt.Errorf("could not edit relation module config: %v", err)
	}
	fmt.Println("1")
	if _, exists := relation.Modules[module]; !exists {
		return fmt.Errorf("relation module %s does not exist", module)
	}
	fmt.Println("2")
	relation.Modules[module].Configuration[config] = value
	if err := r.SaveRelation(*relation); err != nil {
		return fmt.Errorf("could not save relation module configuration: %v", err)
	}
	fmt.Println("3")
	return nil
}
