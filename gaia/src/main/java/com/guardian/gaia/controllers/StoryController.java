package com.guardian.gaia.controllers;


import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.service.StoryService;
import com.guardian.gaia.util.StoryDTOMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.text.ParseException;
import java.util.List;
import java.util.stream.Collectors;

@Controller
public class StoryController {

    private final StoryService storyService;
    private final StoryDTOMapper<Story, StoryDTO> storyMapper;

    public StoryController(StoryService storyService, StoryDTOMapper<Story, StoryDTO> storyMapper) {
        this.storyService = storyService;
        this.storyMapper = storyMapper;
    }


    //Get
    @GetMapping("/story/{id}")
    @ResponseBody
    public Story getStory(int id){
        return new Story();
    }

    @GetMapping("/story/user/{userid}")
    @ResponseBody
    public List<Story> getAllStories(@PathVariable String userid){
        return storyService.getStoryList(userid).stream()
                .map(storyMapper::convertToDto)
                .collect(Collectors.toList());
    }

    //Post
    @PostMapping("/story")
    @ResponseStatus(HttpStatus.OK)
    public void PostStory(@RequestBody StoryDTO storyDto){
        Story story = storyMapper.convertToEntity(storyDto);
        storyService.saveStory(story);
    }

    //Put
    @PutMapping("/story")
    @ResponseStatus(HttpStatus.OK)
    public void PutStory(@RequestBody StoryDTO storyDto){
        Story story = storyMapper.convertToEntity(storyDto);
        storyService.editStory(story);
    }




}
