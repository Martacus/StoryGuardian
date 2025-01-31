package project

import "fmt"

type TagManager struct {
	EntityManager *EntityManager
}

// NewTagManager creates a new TagManager
func NewTagManager(entityManager *EntityManager) *TagManager {
	tagManager := TagManager{
		EntityManager: entityManager,
	}

	return &tagManager
}

func (t *TagManager) RemoveTag(tagName string) error {
	for _, entity := range t.EntityManager.Entities {
		for i, entityTag := range entity.Tags {
			if entityTag == tagName {
				entity.Tags = append(entity.Tags[:i], entity.Tags[i+1:]...)
				if err := t.EntityManager.SaveEntity(*entity); err != nil {
					fmt.Printf("could not save entity %v with removed tag: %v", entity.Name, err)
				}
				break
			}
		}
	}

	s := *t.EntityManager.StoryManager
	for i, tag := range s.Story.Tags {
		if tag == tagName {
			s.Story.Tags = append(s.Story.Tags[:i], s.Story.Tags[i+1:]...)
			if err := s.SaveStory(); err != nil {
				return fmt.Errorf("could not save when removing tag from story: %v", err)
			}
			break
		}
	}

	return nil
}

func (t *TagManager) EditTagName(oldName, newName string) error {
	for _, entity := range t.EntityManager.Entities {
		for i, entityTag := range entity.Tags {
			if entityTag == oldName {
				entity.Tags[i] = newName
				if err := t.EntityManager.SaveEntity(*entity); err != nil {
					fmt.Printf("could not save entity %v with edited tag: %v", entity.Name, err)
				}
				break
			}
		}
	}

	s := *t.EntityManager.StoryManager
	for i, tag := range s.Story.Tags {
		if tag == oldName {
			s.Story.Tags[i] = newName
			if err := s.SaveStory(); err != nil {
				return fmt.Errorf("could not save when editing tag in story: %v", err)
			}
			break
		}
	}

	return nil
}

// == Entity ==//

// RemoveTagFromEntity removes a tag from an entity
func (e *EntityManager) RemoveTagFromEntity(entityId string, tag string) error {
	entity, err := e.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("could not retrieve entity to remove tag: %v", err)
	}

	for i, existingTag := range entity.Tags {
		if existingTag == tag {
			entity.Tags = append(entity.Tags[:i], entity.Tags[i+1:]...)
			break
		}
	}

	if err := e.SaveEntity(*entity); err != nil {
		return fmt.Errorf("could not save entity with removed tag: %v", err)
	}

	return nil
}

// GetEntitiesByTag returns a list of entities with the specified tag
func (e *EntityManager) GetEntitiesByTag(tag string) []Entity {
	var entities []Entity
	for _, entity := range e.Entities {
		for _, entityTag := range entity.Tags {
			if entityTag == tag {
				entities = append(entities, *entity)
				break
			}
		}
	}

	return entities
}

// == Story == //

// CreateTag creates a new tag
func (s *StoryManager) CreateTag(tagName string) error {
	s.Story.Tags = append(s.Story.Tags, tagName)
	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save the tag insertion: %v", err)
	}
	return nil
}

// GetStoryTags returns the tags of the story
func (s *StoryManager) GetStoryTags() []string {
	return s.Story.Tags
}
