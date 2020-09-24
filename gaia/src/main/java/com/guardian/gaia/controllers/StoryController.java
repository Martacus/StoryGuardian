package com.guardian.gaia.controllers;


import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.service.StoryService;
import com.guardian.gaia.util.DTOMapper;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@CrossOrigin(origins = "http://localhost:8000")
@Controller
public class StoryController {

    private final StoryService storyService;
    private final DTOMapper<Story, StoryDTO> storyMapper;

    public StoryController(StoryService storyService, DTOMapper<Story, StoryDTO> storyMapper) {
        this.storyService = storyService;
        this.storyMapper = storyMapper;
    }


    //GET
    @RequestMapping(value="/story/{uuid}", method=RequestMethod.GET)
    @ResponseBody
    public StoryDTO getStory(@PathVariable UUID uuid){
        Optional<Story> story = storyService.getById(uuid);
        return story.map(storyMapper::convertToDto).orElse(null);
    }
    @RequestMapping(value="/story/user/{uuid}", method=RequestMethod.GET)
    @ResponseBody
    public List<StoryDTO> getAllStories(@PathVariable String uuid){
        return storyService.getStoryList(uuid).stream()
                .map(storyMapper::convertToDto)
                .collect(Collectors.toList());
    }

    //POST
    @RequestMapping(value="/story", headers="Accept=application/json", method=RequestMethod.POST)
    @ResponseBody
    public StoryDTO PostStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        story = storyService.saveStory(story);
        return storyMapper.convertToDto(story);
    }

    //PUT
    @RequestMapping(value="/story", headers="Accept=application/json", method=RequestMethod.PUT)
    @ResponseStatus(HttpStatus.OK)
    public void PutStory(@RequestBody StoryDTO storyDTO){
        Story story = storyMapper.convertToEntity(storyDTO);
        storyService.editStory(story);
    }
    
    //DELETE
    @RequestMapping(value = "/story/{uuid}", headers = "Accept=application/json", method = RequestMethod.DELETE)
    @ResponseStatus(HttpStatus.OK)
    public void DeleteStory(@PathVariable UUID uuid){
        storyService.deleteStory(uuid);
    }





}
