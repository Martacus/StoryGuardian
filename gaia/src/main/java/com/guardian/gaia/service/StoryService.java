package com.guardian.gaia.service;

import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.repository.StoryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
public class StoryService {

    private final StoryRepository storyRepository;

    public StoryService(StoryRepository storyRepository) {
        this.storyRepository = storyRepository;
    }

    public List<Story> getStoryList(String userid){
        return storyRepository.findAll();
    }

    public Story saveStory(Story story) {
        return storyRepository.save(story);
    }

    public void editStory(Story story) {
        Optional<Story> story2 = storyRepository.findById(story.getId());
        if(story2.isPresent()){
            storyRepository.save(story);
        }
    }

    public Optional<Story> getById(UUID id) {
        return storyRepository.findById(id);
    }

    public void deleteStory(UUID id) {
        storyRepository.findById(id).ifPresent(storyRepository::delete);
    }
}
