package project

import (
	"fmt"
	"os"
	"path/filepath"
	"storyguardian/internal/constants"
)

type Story struct {
	ProjectDetails
	Description string   `json:"description"`
	Entities    []Entity `json:"entities"`
}

type StoryManager struct {
	ApplicationManager *ApplicationManager
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

func (p *StoryManager) GetStory(projectId string) (*Story, error) {
	for key, projectDetails := range p.ApplicationManager.Config.Projects {
		if key == projectId {
			story := Story{}
			err := writeFilePathToStruct(filepath.Join(projectDetails.Location, constants.ProjectConfigName), &story)
			if err != nil {
				return nil, fmt.Errorf("could not find the story with id: %v | %v", projectId, err)
			}

			return &story, nil
		}
	}
	return nil, fmt.Errorf("could not find the story with id %v", projectId)
}

func (p *StoryManager) SetStoryTitle(storyId string, name string) error {
	story, err := p.GetStory(storyId)
	if err != nil {
		return fmt.Errorf("could not find the story: %v", err)
	}

	story.Name = name
	if err = writeStructToFilePath(story, filepath.Join(story.Location, constants.ProjectConfigName)); err != nil {
		return fmt.Errorf("could not save the title change: %v", err)
	}

	if err = p.ApplicationManager.writeProjectToAppConfig(*story); err != nil {
		return fmt.Errorf("could not save the title change to application config: %v", err)
	}

	return nil
}

func (p *StoryManager) SetStoryDescription(storyId string, description string) (string, error) {
	story, err := p.GetStory(storyId)
	if err != nil {
		return "", fmt.Errorf("could not find the story with id: %v | %v", storyId, err)
	}

	story.Description = description
	if err = writeStructToFilePath(story, filepath.Join(story.Location, constants.ProjectConfigName)); err != nil {
		return "", fmt.Errorf("could not save the description change: %v", err)
	}

	return description, nil
}

func (p *StoryManager) GetStoryImages(storyId string) ([]ImageFile, error) {
	story, err := p.GetStory(storyId)
	if err != nil {
		return make([]ImageFile, 0), err
	}

	images := make([]ImageFile, 0)
	imagesFolderPath := filepath.Join(story.Location, "images")
	files, err := os.ReadDir(imagesFolderPath)
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
