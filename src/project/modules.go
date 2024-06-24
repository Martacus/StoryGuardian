package project

const (
	ImageModuleName   = "images"
	TagListModuleName = "tagList"
)

var availableStoryModules = []string{"description", "entityList", TagListModuleName, ImageModuleName}

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
	case TagListModuleName:
		if err := addStoryTagsModule(s); err != nil {
			return err
		}
		break
	case ImageModuleName:
		if err := addStoryImagesModule(s); err != nil {
			return err
		}
		break
	}

	return nil
}

func addStoryImagesModule(manager *StoryManager) error {
	newImageModule := StoryModule{
		Name: ImageModuleName,
		Configuration: map[string]string{
			"columnSize": "4",
		},
	}
	manager.Story.Modules[ImageModuleName] = newImageModule
	if err := manager.SaveStory(); err != nil {
		return err
	}

	return nil
}

func addStoryTagsModule(manager *StoryManager) error {
	newTagListModule := StoryModule{
		Name: TagListModuleName,
		Configuration: map[string]string{
			"columnSize": "4",
			"itemView":   "list",
		},
	}
	manager.Story.Modules[TagListModuleName] = newTagListModule
	if err := manager.SaveStory(); err != nil {
		return err
	}

	return nil
}
