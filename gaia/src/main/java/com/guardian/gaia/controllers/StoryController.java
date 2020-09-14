package com.guardian.gaia.controllers;


import com.guardian.custodian.Story;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller
public class StoryController {

    @GetMapping("/story/{id}")
    public Story getStory(int id){
        return new Story();
    }

}
