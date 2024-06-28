package project

import "fmt"

const (
	ImageStoryModuleID        = "images"
	TagListStoryModuleID      = "tagList"
	DescriptionStoryModuleID  = "description"
	EntityListStoryModuleID   = "entityList"
	TagListEntityModuleID     = "tagList"
	DescriptionEntityModuleID = "description"
	RelationsEntityModuleID   = "relations"
	RelationInfoModuleID      = "relationInfo"
)

var availableStoryModules = []string{DescriptionStoryModuleID, EntityListStoryModuleID, TagListStoryModuleID, ImageStoryModuleID}
var availableEntityModules = []string{DescriptionEntityModuleID, TagListEntityModuleID, RelationsEntityModuleID}
var availableRelationModules = []string{DescriptionEntityModuleID, RelationInfoModuleID}

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
	case TagListStoryModuleID:
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
	case TagListEntityModuleID:
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

//func (e *EntityManager) GetEntityModules(entityId string, unusedModulesOnly bool) []string {
//	entity, err := e.GetEntity(entityId)
//	if err != nil {
//		return make([]string, 0)
//	}
//
//	if !unusedModulesOnly {
//		return availableEntityModules
//	}
//
//	var unusedEntityModules []string
//	for _, module := range availableEntityModules {
//		if _, ok := entity.Modules[module]; !ok {
//			unusedEntityModules = append(unusedEntityModules, module)
//		}
//	}
//	return unusedEntityModules
//}
//
//func (e *EntityManager) AddEntityModule(entityId string, module string) error {
//	entity, err := e.GetEntity(entityId)
//	if err != nil {
//		return fmt.Errorf("could not add module to entity: %v", err)
//	}
//
//	switch module {
//	case TagListEntityModuleID:
//		if err := addEntityTagsModule(entity, e); err != nil {
//			return err
//		}
//		break
//	}
//
//	return nil
//}
//
//func (e *EntityManager) EditEntityModuleConfig(entityId, module string, config string, value string) error {
//	entity, err := e.GetEntity(entityId)
//	if err != nil {
//		return fmt.Errorf("could not edit entity module config: %v", err)
//	}
//
//	entity.Modules[module].Configuration[config] = value
//
//	if err := e.SaveEntity(*entity); err != nil {
//		return fmt.Errorf("could not edit entity module config: %v", err)
//	}
//	return nil
//}

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
		Name: TagListStoryModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
			"itemView":   "list",
		},
	}
	manager.Story.Modules[TagListStoryModuleID] = newTagListModule
	if err := manager.SaveStory(); err != nil {
		return err
	}

	return nil
}

func addEntityTagsModule(entity *Entity, manager *EntityManager) error {
	newTagListModule := StoryModule{
		Name: TagListEntityModuleID,
		Configuration: map[string]string{
			"columnSize": "4",
			"itemView":   "list",
		},
	}

	entity.Modules[TagListEntityModuleID] = newTagListModule
	if err := manager.SaveEntity(*entity); err != nil {
		return fmt.Errorf("unable to add tag module to entity: %v", err)
	}

	return nil
}
