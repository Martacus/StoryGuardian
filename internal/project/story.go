package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/internal/constants"
	"storyguardian/internal/fileio"
)

type Story struct {
	ProjectDetails
	Description string   `json:"description"`
	Entities    []Entity `json:"entities"`
	Tags        []string `json:"tags"`
}

type StoryManager struct {
	ApplicationManager *ApplicationManager
	Story              *Story
}

type ImageFile struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

func NewStoryManager(appManager *ApplicationManager) *StoryManager {
	return &StoryManager{
		ApplicationManager: appManager,
	}
}

func (s *StoryManager) GetStory(projectId string) (*Story, error) {
	if s.Story != nil && s.Story.Id == projectId {
		return s.Story, nil
	}

	for key, projectDetails := range s.ApplicationManager.Config.Projects {
		if key == projectId {
			story := Story{}
			err := fileio.WriteFilePathToStruct(filepath.Join(projectDetails.Location, constants.ProjectConfigName), &story)
			if err != nil {
				return nil, fmt.Errorf("could not find the story with id: %v | %v", projectId, err)
			}
			s.Story = &story
			return &story, nil
		}
	}
	return nil, fmt.Errorf("could not find the story with id %v", projectId)
}

func (s *StoryManager) SaveStory() error {
	if err := fileio.WriteStructToFilePath(s.Story, filepath.Join(s.Story.Location, constants.ProjectConfigName)); err != nil {
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
