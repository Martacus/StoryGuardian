package com.guardian.gaia.controllers;


import com.guardian.custodian.Entity;
import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.service.StoryService;
import com.guardian.gaia.util.DTOMapper;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@CrossOrigin
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
    @GetMapping(value="/story/{uuid}", produces = {"application/json"})
    public ResponseEntity<StoryDTO> getStory(@PathVariable UUID uuid){
        Optional<Story> story = storyService.getById(uuid);
        if(story.isPresent()){
            return ResponseEntity.ok(story.map(storyMapper::convertToDto).get());
        } else {
            return ResponseEntity.notFound().build();
        }
    }
    @GetMapping(value="/story/user/{uuid}", produces = {"application/json"})
    public ResponseEntity<List<StoryDTO>> getAllStories(@PathVariable String uuid){
        return ResponseEntity.ok(
            storyService.getStoryList(uuid)
                .stream()
                .map(storyMapper::convertToDto)
                .collect(Collectors.toList())
        );
    }

    //POST
    @PostMapping(value="/story", produces = {"application/json"})
    public ResponseEntity<StoryDTO> PostStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        story = storyService.saveStory(story);
        return ResponseEntity.ok(storyMapper.convertToDto(story));
    }

    //PUT
    @PutMapping(value="/story")
    @ResponseStatus(HttpStatus.OK)
    public void PutStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        storyService.editStory(story);
    }
    
    //DELETE
    @DeleteMapping(value="/story/{uuid}")
    @ResponseStatus(HttpStatus.OK)
    public void DeleteStory(@PathVariable UUID uuid){
        storyService.deleteStory(uuid);
    }





}
