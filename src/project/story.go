package project

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"storyguardian/src/constants"
	"storyguardian/src/fileio"
)

type Story struct {
	ProjectDetails
	Description string                 `json:"description"`
	Entities    []Entity               `json:"entities"`
	Tags        []string               `json:"tags"`
	Modules     map[string]StoryModule `json:"modules"`
}

type StoryManager struct {
	ApplicationManager *ApplicationManager
	Story              *Story
}

type ImageFile struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type StoryModule struct {
	Name          string            `json:"name"`
	Configuration map[string]string `json:"configuration"`
}

func NewStoryManager(appManager *ApplicationManager) *StoryManager {
	return &StoryManager{
		ApplicationManager: appManager,
	}
}

func (s *StoryManager) NewStory(projectDirectory string) (*Story, error) {
	projectId := uuid.New().String()
	projectName := filepath.Base(projectDirectory)

	projectDetails := ProjectDetails{
		Id:       projectId,
		Name:     projectName,
		Location: projectDirectory,
	}

	moduleMap := addStoryModules(make(map[string]StoryModule))

	story := Story{
		ProjectDetails: projectDetails,
		Description:    "Placeholder",
		Modules:        moduleMap,
	}

	filePath := filepath.Join(projectDirectory, "story.json")

	if err := fileio.WriteStructToFilePath(story, filePath); err != nil {
		return nil, fmt.Errorf("could not write story to story.json: %w", err)
	}

	if err := s.ApplicationManager.writeProjectDetailsToAppConfig(projectDetails); err != nil {
		return nil, fmt.Errorf("could not write project to application config file: %w", err)
	}

	return &story, nil
}

func (s *StoryManager) GetStory(projectId string, refresh bool) (*Story, error) {
	if !refresh && s.Story != nil && s.Story.Id == projectId {
		return s.Story, nil
	}

	projectDetails, exists := s.ApplicationManager.Config.Projects[projectId]
	if !exists {
		return nil, fmt.Errorf("could not find the story with id %v", projectId)
	}

	story := &Story{}
	err := fileio.WriteFilePathToStruct(filepath.Join(projectDetails.Location, constants.StoryConfigName), story)
	if err != nil {
		return nil, fmt.Errorf("could not find the story with id: %v | %v", projectId, err)
	}

	s.Story = story
	return story, nil
}

func (s *StoryManager) SaveStory() error {
	if err := fileio.WriteStructToFilePath(s.Story, filepath.Join(s.Story.Location, constants.StoryConfigName)); err != nil {
		return fmt.Errorf("could not save story to file system: %v", err)
	}

	return nil
}

func (s *StoryManager) SetStoryTitle(name string) error {
	s.Story.Name = name
	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save the title change: %v", err)
	}

	if err := s.ApplicationManager.writeProjectToAppConfig(*s.Story); err != nil {
		return fmt.Errorf("could not save the title change to application config: %v", err)
	}

	return nil
}

func (s *StoryManager) SetStoryDescription(description string) (string, error) {
	s.Story.Description = description
	if err := s.SaveStory(); err != nil {
		return "", fmt.Errorf("could not save the description change: %v", err)
	}

	return description, nil
}

func (s *StoryManager) GetStoryImages() ([]ImageFile, error) {
	images := make([]ImageFile, 0)
	imagesFolderPath := filepath.Join(s.Story.Location, "images")

	files, err := os.ReadDir(imagesFolderPath)
	if err != nil {
		return nil, fmt.Errorf("could not read image folder: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
				images = append(images, ImageFile{
					Name:     file.Name(),
					Location: filepath.Join(imagesFolderPath, file.Name()),
				})
			}
		}
	}

	return images, nil
}

func (s *StoryManager) CreateTag(tagName string) error {
	s.Story.Tags = append(s.Story.Tags, tagName)
	if err := s.SaveStory(); err != nil {
		return fmt.Errorf("could not save the tag insertion: %v", err)
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

func addStoryModules(moduleMap map[string]StoryModule) map[string]StoryModule {
	moduleMap["description"] = StoryModule{
		Name: "Description",
		Configuration: map[string]string{
			"columnSize": "4",
			"name":       "description",
		},
	}

	moduleMap["entityList"] = StoryModule{
		Name: "Entities",
		Configuration: map[string]string{
			"columnSize": "4",
			"listView":   "list",
			"name":       "entityList",
		},
	}

	moduleMap["images"] = StoryModule{
		Name: "Images",
		Configuration: map[string]string{
			"columnSize": "4",
			"name":       "images",
		},
	}

	return moduleMap
}
