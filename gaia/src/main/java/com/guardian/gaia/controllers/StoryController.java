package com.guardian.gaia.controllers;


import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.service.StoryService;
import com.guardian.gaia.util.DTOMapper;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@CrossOrigin(origins = "http://localhost:8000")
@RestController
@RequestMapping("/story")
public class StoryController {

    private final StoryService storyService;
    private final DTOMapper<Story, StoryDTO> storyMapper;

    public StoryController(StoryService storyService, DTOMapper<Story, StoryDTO> storyMapper) {
        this.storyService = storyService;
        this.storyMapper = storyMapper;
    }


    //GET
    @GetMapping("/story/{uuid}")
    public StoryDTO getStory(@PathVariable UUID uuid){
        Optional<Story> story = storyService.getById(uuid);
        return story.map(storyMapper::convertToDto).orElse(null);
    }
    @GetMapping("/story/user/{uuid}")
    public List<StoryDTO> getAllStories(@PathVariable String uuid){
        return storyService.getStoryList(uuid).stream()
                .map(storyMapper::convertToDto)
                .collect(Collectors.toList());
    }

    //POST
    @PostMapping("/story")
    public StoryDTO PostStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        story = storyService.saveStory(story);
        return storyMapper.convertToDto(story);
    }

    //PUT
    @PutMapping("/story")
    @ResponseStatus(HttpStatus.OK)
    public void PutStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        storyService.editStory(story);
    }
    
    //DELETE
    @DeleteMapping("/story/{uuid}")
    @ResponseStatus(HttpStatus.OK)
    public void DeleteStory(@PathVariable UUID uuid){
        storyService.deleteStory(uuid);
    }





}
