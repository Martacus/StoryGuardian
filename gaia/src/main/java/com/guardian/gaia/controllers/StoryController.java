package com.guardian.gaia.controllers;


import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.service.StoryService;
import com.guardian.gaia.util.StoryDTOMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;

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

    @GetMapping("/story")
    @ResponseBody
    public List<Story> getAllStories(){
        return storyService.getStoryList().stream()
                .map(storyMapper::convertToDto)
                .collect(Collectors.toList());
    }

    @GetMapping("/story/{id}")
    public Story getStory(int id){
        return new Story();
    }

    @PostMapping("/story")
    @ResponseStatus(HttpStatus.OK)
    public void PostStory(){

    }

}
